package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

//成功连接的地址
var servers map[string]*Client

func init() {
	servers = make(map[string]*Client, 0)
}

//连接其他机器
func handlerBrokers(brokers []string) {

	for _, broker := range brokers {
		go clientServe(broker)
	}
}

func clientServe(server string) {
	var t1 = time.NewTimer(HEATBEAT)
	//本机不互联
	if server == thisBroker {
		return
	}

	//成功连接
	if _, ok := servers[server]; ok {
		return
	}

	//log.Println(server)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		deleteBroker(nodeName, server)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		deleteBroker(nodeName, server)
		return
	}

	client := NewSerClient(server, conn)

	//send auth
	msg := &Message{BUFAUTH, &AuthMessage{thisBroker, thisBroker}}
	client.SendMessage(msg)

	go func(client *Client) {
		for {
			select {
			case <-t1.C:
				words := fmt.Sprintf("%s", "ping")
				tmp := &IMMessage{-1, -1, -1, -1, words}
				msg := &Message{BUFVERSION, tmp}
				err := client.SendBuffMessage(msg)
				if err != nil {
					log.Println("ping err so delete server  ", server)
					delete(servers, server)
					deleteBroker(nodeName, server)
				} else {
					//log.Printf("\n send ping %s ==> %s \n ", thisBroker, server)
					t1.Reset(HEATBEAT)
				}

			}
		}
	}(client)

	//不能走以前的流程，不好避免循环发送问题
	//go client.Run()
	go func() {
		for !client.isClose {
			msg := ReaderMessage(conn)
			if msg != nil {
				if conn, ok := client.conn.(net.Conn); ok {
					conn.SetDeadline(RestTimeOut())
					log.Println("接收到信息: ", msg, msg.body)
					//HandlerBuffMessage(client, msg)
					imsg := msg.body.(*IMMessage)
					clients := FindClients(fmt.Sprintf("%d", imsg.Receiver))
					//心跳发送 -1
					if imsg.Sender <= -1 || imsg.Receiver <= -1 {
						if conn, ok := client.conn.(net.Conn); ok {
							conn.SetDeadline(RestTimeOut())
						}
					} else {
						for _, client := range clients {
							client.SendMessage(msg)
						}
					}
				}
			} else {
				client.Stop()
			}
		}
	}()
	//目标对方的连接
	servers[server] = client

}

/*
	for key, v := range servers {
			if key == thisBroker {
				continue
			} else {
				log.Println("发送给服务器: ", v, msg)
				v.SendMessage(msg)
			}
		}
*/

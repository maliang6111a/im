package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

//连接处理
func handler(conn net.Conn) {
	conn.SetDeadline(RestTimeOut())
	client := NewClient(conn)
	//保存起来
	//TODO
	//log.Println(client)
	client.Run()
}

func waitConn(netListener net.Listener, handler func(conn net.Conn)) {
	defer netListener.Close()
	for {
		c, err := netListener.Accept()
		if nil != err {
			log.Println(err)
			log.Fatal("wait connection error!")
			if IsErrClosing(err) {
				return
			}
			return
		}
		//添加超时时间
		c.SetDeadline(time.Now().Add(TIMEOUT))
		handler(c)
	}
}

//监听
//监听处理，监听地址
func listener(handler func(conn net.Conn), addr string) {
	netListener, err := InitTcpListener(addr)
	if err != nil {
		return
	}
	log.Println("TCP Serving at ", addr)
	go waitConn(netListener, handler)
	ServerGroup.Done()
}

func StartTCPServer(addr string) {
	listener(handler, addr)
}

//tcp 信息发送
func SendMessage(conn io.Writer, msg *Message) error {
	buffer := new(bytes.Buffer)
	WriterMessage(buffer, msg)
	buf := buffer.Bytes()
	n, err := conn.Write(buf)
	if err != nil {
		fmt.Printf("sock write error:", err)
		return err
	}
	if n != len(buf) {
		fmt.Printf("write less:%d %d", n, len(buf))
		return errors.New("write less")
	}
	return nil
}

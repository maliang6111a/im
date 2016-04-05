package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/googollee/go-engine.io"
)

type Client struct {
	key     string      //用户ID
	connId  string      // 连接ID  uuid
	conn    interface{} //tcp or engin.io connection
	isClose bool        //连接是否关闭
	isAuth  bool        //该连接是否认证
}

//创建连接
func NewClient(conn interface{}) *Client {
	client := &Client{"", Rand().Hex(), conn, false, false}
	return client
}

//启动该连接
func (this *Client) Run() {
	go this.Reader()
}

//停止
func (this *Client) Stop() {
	this.isAuth = false
	this.isClose = true
	defer RemoveClient(this)
	if conn, ok := this.conn.(net.Conn); ok {
		conn.Close()
	} else if conn, ok := this.conn.(engineio.Conn); ok {
		conn.Close()
	}

}

//处理认证
func (this *Client) handlerAuth(key string) {
	this.isAuth = true
	this.key = key
}

func (this *Client) IsAuthed() bool {
	return this.isAuth && this.key != ""
}

func (this *Client) SendBuffMessage(msg *Message) {
	if conn, ok := this.conn.(net.Conn); ok {
		SendMessage(conn, msg)
	}
}

func (this *Client) SendTextMessageOf(msg *Message) {
	if conn, ok := this.conn.(engineio.Conn); ok {
		immessage := msg.body.(*IMMessage)
		bs, err := json.Marshal(immessage)
		if err != nil {
			return
		}
		//内部编码
		SendEngineIOTextMessage(conn, string(bs))
	}
}

func (this *Client) SendTextMessage(msg string) {
	if conn, ok := this.conn.(engineio.Conn); ok {
		SendEngineIOTextMessage(conn, msg)
	}
}

// engin.io 连接
func (this *Client) handReadSIOConn(conn engineio.Conn) {
	for !this.isClose {
		//msg := ReadEngineIOMessage(conn)
		msg := ReadEngineIOMessageResultStr(conn)
		if msg != "" {
			HandlerTextMessage(this, msg)
		} else {
			this.Stop()
		}
	}
}

//tcp 连接
func (this *Client) handReadTcpConn(conn net.Conn) {
	for !this.isClose {
		msg := ReaderMessage(conn)
		log.Println("接收信息: ", msg)
		if msg != nil {
			conn.SetDeadline(RestTimeOut())
			HandlerBuffMessage(this, msg)
		} else {
			log.Println("服务端关闭tcp...")
			this.Stop()
		}
	}
}

//根据连接类型处理信息
func (this *Client) Reader() {
	if conn, ok := this.conn.(net.Conn); ok {
		this.handReadTcpConn(conn)
	} else if conn, ok := this.conn.(engineio.Conn); ok {
		this.handReadSIOConn(conn)
	}

}

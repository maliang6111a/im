package main

import (
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

// engin.io 连接
func (this *Client) handReadSIOConn(conn engineio.Conn) {
	for !this.isClose {
		msg := ReadEngineIOMessage(conn)
		if msg != nil {
			HandlerMessage(this, msg)
		} else {
			this.Stop()
		}
	}
}

//tcp 连接
func (this *Client) handReadTcpConn(conn net.Conn) {
	for !this.isClose {
		msg := ReaderMessage(conn)
		if msg != nil {
			conn.SetDeadline(RestTimeOut())

			HandlerMessage(this, msg)
		} else {
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

package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

var (

	//MSG_TYPE int = 1

	message_creators map[int]MessageCreator = make(map[int]MessageCreator)
)

func init() {
	//消息协议
	message_creators[BUFVERSION] = func() IMessage { return new(IMMessage) }
	//认证协议信息
	message_creators[BUFAUTH] = func() IMessage { return new(AuthMessage) }
}

type MessageCreator func() IMessage

type IMessage interface {
	ToData() []byte
	FromData(buff []byte) bool
}

//消息
type Message struct {
	version int         //1,byte,消息协议版本
	body    interface{} //len
}

func (this *Message) ToData() []byte {
	if this.body != nil {
		if m, ok := this.body.(IMessage); ok {
			return m.ToData()
		}
	}
	return make([]byte, 0)
}

func (this *Message) FromData(buff []byte) bool {
	proVersion := this.version
	if creator, ok := message_creators[proVersion]; ok {
		c := creator()
		r := c.FromData(buff)
		this.body = c
		return r
	}
	return len(buff) == 0
}

//写入消息头
//version  1
//msg_len  4
func WriterHeader(version byte, msg_len int32, buffer io.Writer) {
	buff := []byte{version}
	buffer.Write(buff)
	binary.Write(buffer, binary.BigEndian, msg_len)
}

//读取消息头信息
// version  1
// msg_len  4
func ReaderHeader(buff []byte) (int, int) {
	buffer := bytes.NewBuffer(buff)
	var msg_len int32
	version, _ := buffer.ReadByte()
	binary.Read(buffer, binary.BigEndian, &msg_len)
	return int(version), int(msg_len)
}

//写消息
func WriterMessage(w io.Writer, msg *Message) {
	body := msg.ToData()

	if len(body) <= 0 {
		return
	}
	//消息头
	WriterHeader(byte(msg.version), int32(len(body)), w)
	//消息体
	w.Write(body)
}

//消息读取
func ReaderMessage(conn io.Reader) *Message {

	buff := make([]byte, 5)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		return nil
	}

	version, msg_len := ReaderHeader(buff)
	if msg_len <= 0 || msg_len >= 32*1024 {
		return nil
	}

	buff = make([]byte, msg_len)
	if _, err = io.ReadFull(conn, buff); err != nil {
		log.Println("socket read body error: ", err)
		return nil
	}

	message := new(Message)
	message.version = version

	//正文
	if !message.FromData(buff) {
		log.Println("body from data error")
		return nil
	}

	return message

}

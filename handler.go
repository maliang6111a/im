package main

import (
	"encoding/json"
	"fmt"
	"net"

	gson "github.com/bitly/go-simplejson"
)

var (
	TmpauthId  = "123"
	TmpauthPwd = "ok"
)

func createTitle() string {
	msg := &IMMessage{0, 0, 0, 0, "消息格式不正确"}
	bs, err := json.Marshal(msg)
	if err != nil {
		return "{\"sender\":0,\"receiver\":0,\"timestamp\":0,\"msgid\":0,\"content\":\"消息格式不正确\"}"
	}
	return string(bs)
}

func createAuth() *Message {
	imsg := &IMMessage{0, 0, 0, 0, "连接未认证..."}
	msg := &Message{BUFVERSION, imsg}
	return msg
}

//文本信息处理
func HandlerTextMessage(client *Client, msg string) {

	result, err := gson.NewJson([]byte(msg))
	//没有认证处理认证信息
	if !client.IsAuthed() {

		if err != nil {
			return
		}
		authMsg := new(AuthMessage)
		authId, err := result.Get("authId").String()
		if err != nil {
			return
		}
		authMsg.authId = authId
		authPwd, err := result.Get("authPwd").String()
		if err != nil {
			return
		}
		authMsg.authPwd = authPwd

		//TODO 验证
		//var flag = false
		//if authMsg.authId == TmpauthId && authMsg.authPwd == TmpauthPwd {
		//	flag = true
		//}
		var flag = true

		if !flag {
			client.Stop()
		} else {
			client.handlerAuth(authMsg.authId)
			PushClient(client)
		}

	} else {

		//消息格式不正确的情况
		defer func() {
			if err := recover(); err != nil {
				client.SendTextMessage(createTitle())
			}
		}()

		sender, _ := result.Get("sender").Int64()
		receiver, _ := result.Get("receiver").Int64()
		content, _ := result.Get("content").String()

		if sender <= 0 || receiver <= 0 || content == "" {
			client.SendTextMessage(createTitle())
		}

		//TODO 时间消息ID处理
		msg := &IMMessage{sender, receiver, 0, 0, content}
		tmp := &Message{BUFVERSION, msg}
		router(client, tmp)
	}

}
func router(client *Client, msg *Message) {
	imsg := msg.body.(*IMMessage)
	clients := FindClients(fmt.Sprintf("%d", imsg.Receiver))
	//心跳发送 -1
	if imsg.Sender <= -1 || imsg.Receiver <= -1 {
		if conn, ok := client.conn.(net.Conn); ok {
			conn.SetDeadline(RestTimeOut())
		}
	} else {

		go func() {
			if !ISServe(client.key) {
				RedirectServer(client, msg)
			}
		}()

		for _, client := range clients {
			client.SendMessage(msg)
		}

	}
}

//buff消息处理
func HandlerBuffMessage(client *Client, msg *Message) {
	//认证信息处理
	if !client.IsAuthed() || msg.version == BUFAUTH {
		if authMessage, ok := msg.body.(*AuthMessage); ok {
			var flag = true
			if ISServe(authMessage.authId) {
				flag = true
			}

			if !flag {
				client.Stop()
			} else {
				client.handlerAuth(authMessage.authId)
				//if !ISServe(authMessage.authId) {
				PushClient(client)
				//}
			}
		}
	} else {
		router(client, msg)
	}
}

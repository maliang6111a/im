package main

import (
	"encoding/json"
	"fmt"
	"log"

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
		var flag = false
		if authMsg.authId == TmpauthId && authMsg.authPwd == TmpauthPwd {
			flag = true
		}

		if !flag {
			client.Stop()
		} else {
			client.handlerAuth(authMsg.authId)
			PushClient(client)
		}

	} else {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				client.SendTextMessage(createTitle())
			}
		}()

		sender, _ := result.Get("sender").Int64()
		receiver, _ := result.Get("receiver").Int64()
		content, _ := result.Get("content").String()

		if sender <= 0 || receiver <= 0 || content == "" {
			client.SendTextMessage(createTitle())
		}

		msg := &IMMessage{sender, receiver, 0, 0, content}
		TextMessageRouter(msg)
	}

}

func TextMessageRouter(msg *IMMessage) {
	clients := FindClients(fmt.Sprintf("%d", msg.Receiver))

	if len(clients) > 0 {
		for _, client := range clients {
			bs, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			//内部编码
			client.SendTextMessage(string(bs))
		}
	}
}

//buff消息处理
func HandlerMessage(client *Client, msg *Message) {

}

package main

import (
	"fmt"
	"log"
)

//消息处理中心
func HandlerMessage(client *Client, msg *Message) {

	fmt.Println(msg.body)

	if !client.IsAuthed() {
		if key, ok := handlerAuthMessage(msg); ok {
			key, result := client.handlerAuth(key)
			if result {
				client.handlerAuth(key)
			}
		} else {
			client.Stop()
		}
	} else {
		if imsg, ok := msg.body.(*IMMessage); ok {
			log.Println(imsg)
		}
	}
}

//认证成功，返回认证用户ID，是否成功标志
func handlerAuthMessage(msg *Message) (string, bool) {
	//des, _ := base64.StdEncoding.DecodeString(imsg.content)
	var autMsg AuthMessage
	var authId = "123"
	var authPwd = "ok"
	if body, ok := msg.body.(*AuthMessage); ok {
		autMsg = body
	}

	//TODO json 转换
	if autMsg != nil {
		if msg.msg_type == BITAUTH {

		} else if msg.msg_type == JSONAUTH {

		}
	}
	return "", true
}

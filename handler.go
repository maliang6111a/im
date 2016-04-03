package main

import (
	"encoding/base64"
	"log"
)

//消息处理中心
func HandlerMessage(msg *Message) {
	//这个地方的转换,注意
	// 根据msgtype 获取engine.io的文本消息，转换乱码根据Base64处理
	if im, ok := msg.body.(*IMMessage); ok {
		des, _ := base64.StdEncoding.DecodeString(im.content)
		log.Println(im.content, string(des))
	}
	//log.Println("信息: ", msg.body)
}

//func

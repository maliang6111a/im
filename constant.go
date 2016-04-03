package main

//msg version
const (
	VERSION     int = 1 //普通IM消息协议
	AUTHVERSION int = 2 //认证消息协议
)

//msg type
const (
	//普通tcp信息
	BITMESSAGE = iota + 1
	//普通engin.io信息
	JSONMESSAGE

	//tcp认证
	BITAUTH
	//engin.io 认证
	JSONAUTH
)

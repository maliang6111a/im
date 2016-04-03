package main

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

//注册器
package main

import (
	"log"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	SPLITCHAT string = ","
)

//启动一个连接watch，监听节点变化，连接其他服务
var brokers []string

//节点名称
var nodeName string
var thisBroker string
var zkConn *zk.Conn
var registed = false

func ISServe(authId string) bool {
	for key, _ := range servers {
		if key == authId {
			return true
		}
	}
	return false
}

func init() {
	nodeName = GetNode()
	thisBroker = GetTcpAddr()

	tmp, _, err := zk.Connect(GetZks(), time.Second*10) //连接超时时间
	if err != nil {
		log.Fatalln("连接zookeeper错误..", err)
		return
	}
	zkConn = tmp
	watch(nodeName)
}

func watch(node string) {

	bs, stat, event, err := zkConn.GetW(node)
	if err != nil {
		log.Fatalln("监听节点错误", err)
		return
	}
	brokers = strings.Split(string(bs), SPLITCHAT)

	go handlerBrokers(brokers)

	for _, broker := range brokers {
		if broker == thisBroker {
			registed = true
		}
	}

	if !registed {
		brokers = append(brokers, thisBroker)
		stat, err = zkConn.Set(nodeName, []byte(strings.Join(brokers, SPLITCHAT)), stat.Version)
		if err != nil {
			log.Fatalln("注册本机错误", err)
			return
		}
		registed = true
	}

	go change(event)
}

func change(event <-chan zk.Event) {
	select {
	case name := <-event:
		{
			if name.Type == zk.EventNodeDataChanged || name.Type == zk.EventNodeDeleted {
				watch(nodeName)
			}
		}
	}
}

func deleteBroker(server string) {
	_, stat, err := zkConn.Get(nodeName)
	if err != nil || len(brokers) <= 0 {
		log.Println("删除broker错误", err)
		return
	}
	for i, broker := range brokers {
		if broker == server {
			if len(brokers) <= 0 {
				brokers = make([]string, 0)
			} else {
				brokers = append(brokers[:i], brokers[i+1:]...)
			}
		}
	}

	zkConn.Set(nodeName, []byte(strings.Join(brokers, SPLITCHAT)), stat.Version)
}

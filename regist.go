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

func ISServe(authId string) bool {
	for key, _ := range servers {
		if key == authId {
			return true
		}
	}
	return false
}

func watch(conn *zk.Conn, registed bool) {

	thisBroker = GetTcpAddr()
	bs, stat, event, err := conn.GetW(nodeName)
	if err != nil {
		//create(conn, nodeName, []byte(thisBroker))
		log.Fatalln("主机节点不存在,不能启动服务....", nodeName)
	} else {
		brs := string(bs)
		//TODO
		//以逗号分隔主机
		//后期改进必须要要添加锁机制
		brokers = strings.Split(brs, SPLITCHAT)
		//检查是否在列表中
		for _, tmpBroker := range brokers {
			if tmpBroker == thisBroker {
				registed = true
			}
		}

		//标志是否注册过
		if !registed {
			brokers = append(brokers, thisBroker)
			conn.Set(nodeName, []byte(strings.Join(brokers, SPLITCHAT)), stat.Version)
			registed = true
		}

		//处理服务
		go handlerBrokers(brokers)

		handlerWatch(conn, event)
	}

}

//创建一个 watch 监听broker 变化情况
func init() {
	nodeName = GetNode()
	var registed = false
	if conn := getZKConn(); conn != nil {
		//启动监听会获取到其他机器信息
		go watch(conn, registed)
	}
}

func getZKConn() *zk.Conn {
	conn, _, err := zk.Connect(GetZks(), time.Second*10)
	if err != nil {
		return nil
	}
	return conn
}

func CloseConn(conn *zk.Conn) {
	if conn == nil {
		return
	}
	conn.Close()
}

//处理watch
func handlerWatch(conn *zk.Conn, event <-chan zk.Event) {
	select {
	case name := <-event:
		if name.Type == zk.EventNodeDataChanged || name.Type == zk.EventNodeDeleted {
			watch(conn, true)
		}
	}
}

func create(node string, v []byte) (string, error) {
	conn := getZKConn()
	defer conn.Close()
	//节点，值，xxx,权限
	return conn.Create(node, v, int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll))
}

func deleteBroker(node, brokerAddr string) bool {
	conn := getZKConn()
	defer conn.Close()

	if conn != nil {
		bs, stat, err := conn.Get(nodeName)
		if err != nil {
			return false
		}
		tmps := strings.Split(string(bs), SPLITCHAT)

		for i, tmp := range tmps {
			if brokerAddr == tmp {
				tmps = append(tmps[:i], tmps[i+1:]...)
			}
		}
		_, err = conn.Set(node, []byte(strings.Join(tmps, SPLITCHAT)), stat.Version)
		if err != nil {
			return false
		}
		return true
	}
	return false

}

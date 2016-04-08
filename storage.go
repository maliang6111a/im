package main

import (
	"log"
	"sync"
)

//"log"

var caches map[string][]*Client
var lock *sync.Mutex

func init() {
	log.Println("存储内存分配...")
	caches = make(map[string][]*Client)
	lock = &sync.Mutex{}
}

func connections() {
	var i = 0
	for _, v := range caches {
		i += len(v)
	}
	log.Println("当前连接: ", i)
}

func exits(key string) []*Client {
	return caches[key]
}

func delClient(connId string, clients []*Client) bool {
	lock.Lock()
	defer lock.Unlock()
	if len(clients) <= 0 {
		return false
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("storage line 40 : ", err)
		}
	}()

	var key string
	for i, tmp := range clients {
		if tmp.connId == connId {

			if i <= 0 {
				clients = make([]*Client, 0)
			} else {
				clients = append(clients[:i], clients[i+1:]...)
			}

			key = tmp.key
		}
	}

	if len(clients) <= 0 {
		delete(caches, key)
	}
	go connections()
	return true
}

func FindClients(key string) []*Client {
	lock.Lock()
	defer lock.Unlock()
	return exits(key)
}

func PushClient(client *Client) {
	lock.Lock()
	defer lock.Unlock()
	if client == nil || client.key == "" {
		return
	}
	if client.isAuth && !client.isClose {
		clients := exits(client.key)
		clients = append(clients, client)
		caches[client.key] = clients
	}
	go connections()
}

func RemoveClient(client *Client) {
	if client == nil {
		return
	}
	clients := exits(client.key)

	if len(clients) <= 0 {
		return
	}

	delClient(client.connId, clients)
}

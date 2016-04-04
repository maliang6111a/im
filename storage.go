package main

import (
	"log"
)

//"log"

var caches map[string][]*Client

func init() {
	caches = make(map[string][]*Client)
}

func exits(key string) []*Client {
	return caches[key]
}

func delClient(connId string, clients []*Client) bool {
	if len(clients) <= 0 {
		return false
	}
	var key string
	for i, tmp := range clients {
		if tmp.connId == connId {
			clients = append(clients[:i], clients[i+1:]...)
			key = tmp.key
		}
	}

	if len(clients) <= 0 {
		delete(caches, key)
	}
	log.Println("当前连接: ", len(caches))

	return true
}

func FindClients(key string) []*Client {
	return exits(key)
}

func PushClient(client *Client) {
	if client == nil || client.key == "" {
		return
	}
	if client.isAuth && !client.isClose {
		clients := exits(client.key)
		clients = append(clients, client)
		caches[client.key] = clients
		log.Println("当前连接 : ", len(caches))
	}
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

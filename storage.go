package main

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

	for i, tmp := range clients {
		if tmp.connId == connId {
			clients = append(clients[:i], clients[i+1:]...)
			return true
		}
	}

	return false
}

func PushClient(client *Client) {
	if client == nil || client.key == "" {
		return
	}
	if client.isAuth && !client.isClose {
		clients := exits(client.key)
		clients = append(clients, client)
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

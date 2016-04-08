package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net"
	"net/http"

	"log"

	"github.com/googollee/go-engine.io"
)

type SIOServer struct {
	server *engineio.Server
}

func (s *SIOServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Origin") != "" {
		w.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", `Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma,
		Last-Modified, Cache-Control, Expires, Content-Type`)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	s.server.ServeHTTP(w, req)
}

func StartSocketIO(socket_io_address string) {
	server, err := engineio.NewServer(nil)
	//server.SetPingInterval(time.Second * 20)
	//server.SetPingTimeout(TIMEOUT)
	if err != nil {
		return
	}

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Println("accept connect fail")
			}
			handlerEngineIOClient(conn)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/engine.io/", &SIOServer{server})
	log.Printf("EngineIO Serving at %s\n", socket_io_address)
	//http 服务
	httpService(socket_io_address, mux)

}

func listenAndServe(laddr string, handler http.Handler) {
	var err error
	var netListener net.Listener
	netListener, err = InitTcpListener(laddr)
	if err != nil {
		log.Fatalf("start fail: %v", err)
	}

	server := &http.Server{Handler: handler}
	err = server.Serve(netListener)
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
	ServerGroup.Done()
}

func httpService(laddr string, handler http.Handler) {
	go func() {
		listenAndServe(laddr, handler)
	}()
}

func handlerEngineIOClient(conn engineio.Conn) {
	client := NewClient(conn)
	client.Run()
}

func SendEngineIOBinaryMessage(conn engineio.Conn, msg *Message) error {
	w, err := conn.NextWriter(engineio.MessageBinary)
	if err != nil {
		log.Println("get next writer fail")
		return err
	}
	log.Println("message version:", msg.version)
	err = SendMessage(w, msg)

	if err != nil {
		log.Println("engine io write error")
		return err
	}
	defer w.Close()
	return nil
}

func SendEngineIOTextMessage(conn engineio.Conn, msg string) error {
	w, err := conn.NextWriter(engineio.MessageText)
	defer w.Close()
	if err != nil {
		log.Println("get next writer fail")
		return err
	}
	bs := base64.StdEncoding.EncodeToString([]byte(msg))
	_, err = w.Write([]byte(bs))
	if err != nil {
		return err
	}
	return nil
}

func ReadEngineIOMessageResultStr(conn engineio.Conn) string {
	t, r, err := conn.NextReader()
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}
	r.Close()
	if t == engineio.MessageText {
		//log.Println("接收信息1: ", string(b))
		bs, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			return ""
		}
		return string(bs)
	}
	return ""
}

//目前没用
func ReadEngineIOMessage(conn engineio.Conn) *Message {
	t, r, err := conn.NextReader()
	if err != nil {
		return nil
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}
	r.Close()
	if t == engineio.MessageText {
		iMsg := &IMMessage{0, 0, 0, 0, string(b)}
		msg := &Message{1, iMsg}
		return msg
	} else {
		log.Println("发信消息")
		return ReadBinaryMesage(b)
	}
}

func ReadBinaryMesage(b []byte) *Message {
	reader := bytes.NewReader(b)
	return ReaderMessage(reader)
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	//超时时间
	TIMEOUT = 40 * time.Second

	ServerGroup = &sync.WaitGroup{}
)

func RestTimeOut() time.Time {
	return time.Now().Add(TIMEOUT)
}

func InitTcpListener(laddr string) (net.Listener, error) {
	ServerGroup.Add(1)
	return net.Listen("tcp", laddr)
}

func IsErrClosing(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		err = opErr.Err
	}
	return "use of closed network connection" == err.Error()
}

//待学习
//TODO
func waitSignal() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM)
	for {
		sig := <-ch
		log.Println("singal:", sig.String())
		switch sig {

		case syscall.SIGTERM, syscall.SIGINT:
			//shutdown()
			return nil
		case syscall.SIGQUIT:
			//gracefulShutdown()
			return nil
		case syscall.SIGHUP:
			//restart(sig)
			//gracefulShutdown()
			return nil
		}
	}
}

func Wait() {
	waitSignal()
	log.Println("close main process")
}

func main() {
	log.Println("启动服务....")

	go StartTCPServer(fmt.Sprintf("0.0.0.0:%d", 9997))
	go StartSocketIO(fmt.Sprintf("0.0.0.0:%d", 9998))
	ServerGroup.Wait()
	Wait()
}

package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

//适用于mysql远程客户端, 不适用于http
func main() {
	listen, err := net.Listen("tcp", value.ADDR)
	util.ErrFatal(err)

	go func() {
		for {
			log.Printf("连接数:%s", im)
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("新连接")
		go handleConn(conn)
	}
}

var im = &v.IntMux{}

func handleConn(from net.Conn) {
	im.Inc()
	defer from.Close()
	defer im.Dec()

	to, err := net.Dial("tcp", value.TO_ADDR)
	if err != nil {
		return
	}
	defer to.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		io.Copy(to, from)
		wg.Done()
	}()
	go func() {
		io.Copy(from, to)
		wg.Done()
	}()
	wg.Wait()
}


package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

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

func handleConn(conn net.Conn) {
	im.Inc()
	defer conn.Close()
	defer im.Dec()

	dial, err := net.Dial("tcp", value.TO_ADDR)
	if err != nil {
		return
	}
	defer dial.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			_, err := io.Copy(dial, conn)
			if err != nil {
				log.Println("conn-dial error:",err)
				break
			}
		}
		log.Println("连接结束1")
		wg.Done()
	}()
	go func() {
		for {
			_, err := io.Copy(conn, dial)
			if err != nil {
				log.Println("dial-conn error:",err)
				break
			}
		}
		log.Println("连接结束2")
		wg.Done()
	}()
	wg.Wait()
	log.Println("连接结束")
}

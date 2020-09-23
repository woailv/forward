package util

import (
	"context"
	"errors"
	"ldy/forward/v"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func ErrFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var im = &v.IntMux{}

type Conf struct {
	UseIPWhite bool                `yaml:"use_ip_white"`
	IPWhite    []string            `yaml:"ip_white"`
	Mapping    []map[string]string `yaml:"mapping"`
}

func Run(conf *Conf) {
	wg := &sync.WaitGroup{}
	for _, item := range conf.Mapping {
		for a, b := range item {
			wg.Add(1)
			go listenAndForward(a, b, wg, conf)
		}
	}
	wg.Wait()
}

func listenAndForward(listenAddr, localAddr string, wg *sync.WaitGroup, conf *Conf) error {
	defer wg.Done()
	listen, err := net.Listen("tcp", ":"+listenAddr)
	ErrFatal(err)
	for {
		remoteConn, err := listen.Accept()
		if err != nil {
			continue
		}
		forbbiden := false
		if conf.UseIPWhite {
			forbbiden = true
			for k, ip := range conf.IPWhite {
				if strings.Contains(remoteConn.RemoteAddr().String(), ip) {
					forbbiden = false
					break
				}
				if k == len(conf.IPWhite)-1 {
					log.Println("不在白名单内:", remoteConn.RemoteAddr())
				}
			}
		}
		if forbbiden {
			remoteConn.Close()
			continue
		}
		go func() {
			im.Inc()
			log.Println("新建连接,当前连接数:", im)
			localConn := mustDial("localhost:" + localAddr)
			RunForward(remoteConn, localConn)
			im.Dec()
			log.Println("断开连接,当前连接数:", im)
		}()
	}
}

func RunForward(connRemote, connLocal net.Conn) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go forward(connRemote, connLocal, ctx, cancelFunc, wg)
	go forward(connLocal, connRemote, ctx, cancelFunc, wg)
	wg.Wait()
}

func forward(dst, src net.Conn, ctx context.Context, cancel func(), wg *sync.WaitGroup) {
	defer wg.Done()
	_ = readWrite(dst, src, ctx)
	cancel()
	funcPipeExe(dst.Close, src.Close)
}

func funcPipeExe(f ...interface{}) {
	for _, item := range f {
		switch x := item.(type) {
		case func():
			x()
		case func() error:
			if err := x(); err != nil {
				log.Println("funcPipeExe.warning:", err)
			}
		default:
			panic("funcPipeExe:不支持的类型")
		}
	}
}

func readWrite(dst net.Conn, src net.Conn, ctx context.Context) error {
	var buf = make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			log.Println("for readwrite ctx.Done")
			return errors.New("contextDone")
		default:
			//log.Println("for readwrite default")
			nr, er := src.Read(buf)
			if er != nil {
				return er
			}
			//log.Println("for readwrite default read n:", nr)
			if nr == 0 {
				return errors.New("没有读取到内容")
			}
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if ew != nil {
					return ew
				}
				if nr != nw {
					return errors.New("目标写入数据不完整")
				}
			}
		}
	}
}

func mustDial(host string) net.Conn {
	log.Printf("开始建立连接:%s", host)
HERE:
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("建立连接:%s失败:%s,3秒后重试", host, err)
		time.Sleep(time.Second * 3)
		goto HERE
	}
	log.Printf("建立连接%s成功", host)
	return conn
}

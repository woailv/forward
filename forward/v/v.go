package v

import (
	"fmt"
	"sync"
)

type IntMux struct {
	Val int
	sync.Mutex
}

func (this *IntMux) Inc() {
	this.Lock()
	this.Val++
	this.Unlock()
}
func (this *IntMux) Dec() {
	this.Lock()
	this.Val--
	this.Unlock()
}

func (this *IntMux) String() string {
	this.Lock()
	defer this.Unlock()
	return fmt.Sprintf("%d", this.Val)
}

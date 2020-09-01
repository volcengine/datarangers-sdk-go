package asynsdk

import (
	"sync"
	"time"
)

type mq struct {
	sync.Mutex
	queue chan *message
}

func newMq() *mq {
	mq1 := &mq{
		queue: make(chan *message, singelConfig.A.Mqlen),
	}
	return mq1
}

func (p *mq) push(item *message) {
	go func() {
		select {
		case p.queue <- item:
			break
		case <-time.After(1 * time.Second):
			a, _ := logJson(item)
			errlogger.Println(a)
			break
		}
	}()
}

func (p *mq) pop() *message {
	select {
	case item := <-p.queue:
		return item
	}
}

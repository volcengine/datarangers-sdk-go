package datarangers_sdk

import (
	"encoding/json"
	"time"
)

type mq struct {
	queue chan *dancemsg
}

func newMq() *mq {
	mq1 := &mq{
		queue: make(chan *dancemsg, confIns.Asynconf.Mqlen),
	}
	return mq1
}

func (p *mq) push(dmg *dancemsg) {

	go func() {
		select {
		case p.queue <- dmg:
			break
		case <-time.After(1 * time.Second):
			a, _ := json.Marshal(dmg)
			warn("消息未进入队列，丢弃")
			errlogger.Println(string(a))
			break
		}
	}()
}

func (p *mq) pop() *dancemsg {
	select {
	case item := <-p.queue:
		return item
	}
}

package datarangers_sdk

import (
	"encoding/json"
	"time"
)

type mq struct {
	queue chan *dancemsg
}

func newMq() *mq {
	q := &mq{
		queue: make(chan *dancemsg, confIns.Asynconf.Mqlen),
	}
	return q
}

func (q *mq) push(dmg *dancemsg) {
	go func() {
		select {
		case q.queue <- dmg:
			break
		case <-time.After(1 * time.Second):
			a, _ := json.Marshal(dmg)
			warn("消息未进入队列，丢弃")
			errlogger.Println(string(a))
			break
		}
	}()
}

func (q *mq) pop() *dancemsg {
	select {
	case item := <-q.queue:
		return item
	}
}

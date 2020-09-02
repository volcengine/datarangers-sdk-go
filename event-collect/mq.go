package event_collect

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
		case <-time.After(30 * time.Second):
			a, _ := json.Marshal(dmg)
			errlogger.Println(a)
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

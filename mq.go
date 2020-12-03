package datarangers_sdk
/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
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

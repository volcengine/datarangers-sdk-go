/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package datarangers_sdk

import (
	"errors"
	"fmt"
	"time"
)

type messageQueue struct {
	queue chan interface{}
}

func newMq() *messageQueue {
	q := &messageQueue{
		queue: make(chan interface{}, confIns.AsynConfig.QueueSize),
	}
	return q
}

func (q *messageQueue) push(dmg interface{}) {
	if confIns.AsynConfig.WaitTimeout > 0 {
		waitTimeout := time.Duration(confIns.AsynConfig.WaitTimeout)
		select {
		case q.queue <- dmg:
			// 数据发送成功
		case <-time.After(waitTimeout * time.Millisecond):
			// 发送超时，进行相应处理
			errMsg := fmt.Sprintf("datarangers send Queue reach max length: %d, and waitTiemout: %dms",
				confIns.AsynConfig.QueueSize, waitTimeout,
			)
			fatal(errMsg)
			handleErr([]interface{}{dmg}, errors.New(errMsg))
		}
	} else {
		q.queue <- dmg
	}
}

func (q *messageQueue) pop() interface{} {
	select {
	case item := <-q.queue:
		return item
	}
}

func (q *messageQueue) popBatch(size int, waitTimeMs time.Duration) []interface{} {
	items := make([]interface{}, 0)
	for len(items) < size {
		select {
		case item := <-q.queue:
			items = append(items, item)
		case <-time.After(waitTimeMs):
			return items
		}
	}
	return items

}

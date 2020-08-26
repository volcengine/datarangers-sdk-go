package mario_collector

import (
	"code.byted.org/data/datarangers-sdk-go"
	"fmt"
)

type client struct {
	producer *executePool
	consumer *consumer
	isFirst  bool
}

func NewClient(num int) (*client, error) {
	ans := &client{}
	var err error
	if ans.producer, err = NewExecutorPool(num); err != nil {
		return nil, fmt.Errorf("producer初始化错误")
	}

	if ans.consumer, err = NewConsumer(num); err != nil {
		return nil, fmt.Errorf("consumer初始化错误")
	}

	return ans, nil
}

func (this *client) isComplete() bool {
	return this.producer.isComplete() && this.consumer.register.complete
}

/**
 目前的发送线程 ， 是发完就停了。 以后再改。！！！！
	1. 会莫名其妙的停下来 。
	2. 会有大量的timeout 。
*/
func (this *client) submit(user *mario_collector.User, header *mario_collector.Header, ee interface{}) error {

	var1 := 0
	if event, ok := ee.(*mario_collector.Event); ok {
		var tt []*mario_collector.Event
		tt = append(tt, event)
		this.producer.execute(user, header, tt)
		var1 += 1
	}
	if events, ok := ee.([]*mario_collector.Event); ok {
		this.producer.execute(user, header, events)
		var1 += 1
	}
	if var1 == 0 {
		return fmt.Errorf("事件类型为[]*pb_event.Event或*pb_event.Event")
	}
	if !this.isFirst {
		this.consumer.execute()
		this.isFirst = true
	}
	return nil
}

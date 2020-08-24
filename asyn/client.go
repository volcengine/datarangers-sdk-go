package asyn

import (
	"code.byted.org/data/datarangers-sdk-go/pb_event"
	"fmt"
)

type client struct {
	producer *executePool
	consumer *consumer
}

func NewClient(num int)(*client, error){
	ans := &client{}
	var err error
	if ans.producer, err = NewExecutorPool(num); err!=nil{
		return nil, fmt.Errorf("producer初始化错误")
	}

	if ans.consumer, err = NewConsumer(num); err!=nil{
		return nil, fmt.Errorf("consumer初始化错误")
	}

	return ans,nil
}

func (this *client) isComplete()bool{
	return this.producer.isComplete()
}

func (this *client) submit(user *pb_event.User, header *pb_event.Header, ee interface{})error{

	if event,ok:=ee.(*pb_event.Event); ok{
		return this.producer.execute(user, header, event)
	}

	if events,ok:=ee.([]*pb_event.Event); ok{

		return this.producer.execute(user, header, events)
	}

	return fmt.Errorf("事件类型为[]*pb_event.Event或*pb_event.Event")
}
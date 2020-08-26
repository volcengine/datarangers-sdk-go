package mario_collector

import (
	"code.byted.org/data/datarangers-sdk-go"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"testing"
	"time"
)

var user = &mario_collector.User{
	UserUniqueId: proto.String("uuid"),
	DeviceId:     proto.Uint64(55555),
	//WebId:        proto.Uint64(55555),
	//Ssid:         proto.String("bb3217e3-dfc5-4520-9cfc-1ab79d854952"),
}

// 自定义headers，大部分情况下不需要。标准字段请使用预定义字段。
var headers = map[string]interface{}{"update_version_code": 542}
var jsonBytes, _ = json.Marshal(headers)
var header = &mario_collector.Header{
	AppId:      proto.Uint32(10000013), //tob产品使用appkey进行上报，appId设置为0即可
	Install_id: proto.Uint64(123),
	Headers:    proto.String(string(jsonBytes)),
}

// event名和params的定义即具体埋点的规范由PM制定
// 测试时请设置XStagingFlag=true
var params = map[string]interface{}{"enter_from": "click_headline", "group_id": 123}
var paramsBytes, _ = json.Marshal(params)
var event = &mario_collector.Event{
	Event: proto.String("test_go_detail"),
	//Time:         proto.Uint32(123),
	Params:    proto.String(string(paramsBytes)),
	SessionId: proto.String("11sgsgsjhbdjsad"),
}
var event1 = &mario_collector.Event{
	Event:     proto.String("test_go_detail1"),
	Time:      proto.Uint32(123),
	Params:    proto.String(string(paramsBytes)),
	SessionId: proto.String("22sgsgsjhbdjsad"),
}

//
//func TestProducesser(t *testing.T) {
//	p,err:=NewExecutorPool(GoRouNum)
//	if err!=nil{
//		fmt.Println(err)
//	}else{
//		for i:=0; i<150; i++{
//			user.UserUniqueId = proto.String("sss"+strconv.Itoa(i))
//			p.execute(user, header, event)
//		}
//		time.Sleep(1*time.Second)
//	}
//}

//测试 多个 event
func TestProducesser2(t *testing.T) {
	p, err := NewExecutorPool(GoRouNum)
	var events []*mario_collector.Event
	events = append(events, event)
	events = append(events, event1)
	if err != nil {
		fmt.Println(err)
	} else {
		for i := 0; i < 150; i++ {
			user.UserUniqueId = proto.String("sss" + strconv.Itoa(i))
			p.execute(user, header, events)
		}
		time.Sleep(time.Second)
	}
}

func TestNewConsumer(t *testing.T) {
	c, err := NewConsumer(GoRouNum)
	if err != nil {
		fmt.Println(err)
	} else {
		c.collcet(0)
		for c.ans != GoRouNum {
			fmt.Print("ok : ", c.oknum)
			fmt.Println("fail : ", c.failnum)
			time.Sleep(2 * time.Second)
		}
	}
	fmt.Println(c.ans)
}

func TestNewConsumer2(t *testing.T) {
	c, _ := NewConsumer(GoRouNum)
	c.execute()
	for c.ans != GoRouNum {
		fmt.Print("ok : ", c.oknum)
		fmt.Println("fail : ", c.failnum)
		time.Sleep(2 * time.Second)
	}
}

func TestNewClient(t *testing.T) {
	client, _ := NewClient(8)
	var events []*mario_collector.Event
	events = append(events, event)
	events = append(events, event1)
	for i := 0; i < 10000; i++ {
		user.UserUniqueId = proto.String("sss" + strconv.Itoa(i))
		client.submit(user, header, events)
	}
	for !client.isComplete() {
		fmt.Println(" - waitting -")
		time.Sleep(2 * time.Second)
	}
}

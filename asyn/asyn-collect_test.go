package asyn

import (
	"code.byted.org/data/datarangers-sdk-go/pb_event"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"testing"
	"time"
)

var user = &pb_event.User{
	UserUniqueId: proto.String("uuid"),
	DeviceId:     proto.Uint64(55555),
	//WebId:        proto.Uint64(55555),
	//Ssid:         proto.String("bb3217e3-dfc5-4520-9cfc-1ab79d854952"),
}

// 自定义headers，大部分情况下不需要。标准字段请使用预定义字段。
var headers = map[string]interface{}{"update_version_code": 542}
var jsonBytes, _ = json.Marshal(headers)
var header = &pb_event.Header{
	AppId:        proto.Uint32(10000013), //tob产品使用appkey进行上报，appId设置为0即可
	Install_id:   proto.Uint64(123),
	Headers:      proto.String(string(jsonBytes)),
}

// event名和params的定义即具体埋点的规范由PM制定
// 测试时请设置XStagingFlag=true
var params = map[string]interface{}{"enter_from": "click_headline", "group_id": 123}
var paramsBytes, _ = json.Marshal(params)
var event = &pb_event.Event{
	Event:        proto.String("test_go_detail"),
	//Time:         proto.Uint32(123),
	Params:       proto.String(string(paramsBytes)),
	SessionId:    proto.String("11sgsgsjhbdjsad"),
}
var event1 = &pb_event.Event{
	Event:        proto.String("test_go_detail1"),
	Time:         proto.Uint32(123),
	Params:       proto.String(string(paramsBytes)),
	SessionId:    proto.String("22sgsgsjhbdjsad"),
}



func TestProducesser(t *testing.T) {
	p,err:=NewExecutorPool(GoRouNum)
	if err!=nil{
		fmt.Println(err)
	}else{
		for i:=0; i<150; i++{
			user.UserUniqueId = proto.String("sss"+strconv.Itoa(i))
			p.execute(user, header, event)
		}
		time.Sleep(1*time.Second)
	}
}

//测试 多个 event
func TestProducesser2(t *testing.T) {
	p,err:=NewExecutorPool(GoRouNum)
	if err!=nil{
		fmt.Println(err)
	}else{
		var events []*pb_event.Event
		events = append(events, event)
		events = append(events, event1)
		p.execute(user, header, events)
		time.Sleep(1*time.Second)
	}
}


func TestNewConsumer(t *testing.T) {
	c,err := NewConsumer(GoRouNum)
	if err!=nil {
		fmt.Println(err)
	}else{
		c.execute(true)
		time.Sleep(2 * time.Second)
		for c.ans!=GoRouNum {
			fmt.Print("ok : ",c.oknum)
			fmt.Println("fail : ",c.failnum)
			time.Sleep(2 * time.Second)
		}
	}
	fmt.Println(c.ans)
}

func TestNewConsumer2(t *testing.T) {
	c,_ := NewConsumer(GoRouNum)
	c.execute(true)
	for c.ans!=GoRouNum  {
		fmt.Print("ok : ",c.oknum)
		fmt.Println("fail : ",c.failnum)
		time.Sleep(2 * time.Second)
	}
}
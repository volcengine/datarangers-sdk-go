package mario_collector

import (
	"code.byted.org/data/mario_collector/pb_event"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"testing"
)

var user = &pb_event.User{
	UserUniqueId: proto.String("411668129935874"),
	UserType:     proto.Uint32(14),
	UserId:       proto.Uint64(13468385056),
	UserIsAuth:   proto.Bool(false),
	UserIsLogin:  proto.Bool(false),
	DeviceId:     proto.Uint64(6902145090442234884),
	WebId:        proto.Uint64(6902145090442234884),
	Ssid:         proto.String("003217e3-dfc5-4520-9cfc-1ab79d854952"),
}

// 自定义headers，大部分情况下不需要。标准字段请使用预定义字段。
var headers = map[string]interface{}{"update_version_code": 542}
var jsonBytes, _ = json.Marshal(headers)
var header = &pb_event.Header{
	AppId:        proto.Uint32(10000010), //tob产品使用appkey进行上报，appId设置为0即可
	AppName:      proto.String("news_article"),
	AppInstallId: proto.Uint64(123),
	AppPackage:   proto.String("com.ss.android.article.news"),
	AppChannel:   proto.String("App Store"),
	AppVersion:   proto.String("5.1.3"),
	OsName:       proto.String("Android"),
	OsVersion:    proto.String("6.0.1"),
	DeviceModel:  proto.String("SM-G9250"),
	AbClient:     proto.String("a1,b1,c2,e1,f1,g2"),
	AbVersion:    proto.String("91223,83097"),
	TrafficType:  proto.String("app"),
	ClientIp:     proto.String("10.100.1.1"),
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
	//XStagingFlag: proto.Bool(true), // 设置测试标志
}
var event1 = &pb_event.Event{
	Event:        proto.String("test_go_detail1"),
	//Time:         proto.Uint32(123),
	Params:       proto.String(string(paramsBytes)),
	SessionId:    proto.String("22sgsgsjhbdjsad"),
	//XStagingFlag: proto.Bool(true), // 设置测试标志
}
var events = []*pb_event.Event{event}


var event2 = &pb_event.Event{
	Event:        proto.String("aaaaaaaa"),
	//Time:         proto.Uint32(123),
	Params:       proto.String(string(paramsBytes)),
	SessionId:    proto.String("66sgsgsjhbdjsad"),
}


var event3 = &pb_event.Event{
	//Event:        proto.String("event_no_user_id_but_header_uuid"),
	Event:        proto.String("header+++Web+installID+event_user_id_and_header_uuid333+testwebid0"),
	//Time:         proto.Uint32(123),
	Params:       proto.String(string(paramsBytes)),
	SessionId:    proto.String("dnijsh23sd23e"),
}
var appevents = []*pb_event.Event{event3}



// MCS上报测试
func TestMcsCollectEvents(t *testing.T) {
	// 不同地区使用不同的上报url，实际应使用的上报地址请咨询PM
	mcsCollector := NewMcsCollector("http://10.225.129.3:31081/v2/event/json", "e9ad30af640c2dcf20a366417e9b287d")
	resp, err := mcsCollector.WebCollectEvents(user, header, events)
	if err == nil {
		defer resp.Body.Close()                        // 保证连接复用
		fmt.Println("response code:", resp.StatusCode) // 查看resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // 查看resp.Body
	}
}



func TestWebCollector(t *testing.T) {

	mcsCollector := NewWebMpCollector("92cf8853e3647c6ed244e5cc94c7704c")
	resp, err := mcsCollector.WebCollectEvents(user, header, events)

	if err == nil {
		defer resp.Body.Close()                        // 保证连接复用
		fmt.Println("response code:", resp.StatusCode) // 查看resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // 查看resp.Body
	}else{
		fmt.Println(err)
	}
}


func TestAppCollector(t *testing.T) {

	mcsCollector := NewAppCollector("e76878c039b9e12be381eb326e247fe2")
	resp, err := mcsCollector.AppCollectEvents(user, header, appevents)

	if err == nil {
		defer resp.Body.Close()                        // 保证连接复用
		fmt.Println("response code:", resp.StatusCode) // 查看resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // 查看resp.Body
	}else{
		fmt.Println(err)
	}
}


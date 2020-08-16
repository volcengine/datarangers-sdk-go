# mario_collector
## 简介
服务端GOSDK第一版
## SDK已支持功能

 
 
## 使用方式
### Web端、小程序端的Event上报
```go
package main
 
 import (
    "code.byted.org/data/mario_collector"
    "code.byted.org/data/mario_collector/pb_event"
    "code.byted.org/gopkg/logs"
    "encoding/json"
    "github.com/golang/protobuf/proto"
 )
 
 func main() {
    user := &pb_event.User{
       UserUniqueId: proto.String("14_13468385056"),
       UserType:     proto.Uint32(14),
       UserId:       proto.Uint64(13468385056),
       UserIsAuth:   proto.Bool(false),
       UserIsLogin:  proto.Bool(false),
       DeviceId:     proto.Uint64(13468385056),
       WebId:        proto.Uint64(12345),
    }
    // 自定义headers，大部分情况下不需要。标准字段请使用预定义字段。
    headers := map[string]interface{}{}
    headers["update_version_code"] = 542
    jsonBytes, _ := json.Marshal(headers)
    header := &pb_event.Header{
       AppId:        proto.Uint32(7), //必选
       AppName:      proto.String("joke"),
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
       UtmSource:    proto.String(""),
       UtmCampaign:  proto.String(""),
       UtmMedium:    proto.String(""),
       ClientIp:     proto.String("10.100.1.1"),
       Headers:      proto.String(string(jsonBytes)),
    }
 
    // event名和params的定义即具体埋点的规范由PM制定
    params := map[string]interface{}{}
    params["enter_from"] = "click_headline"
    params["group_id"] = 123
    paramsBytes, _ := json.Marshal(params)
    event := &pb_event.Event{
       Event:  proto.String("test_go_detail"),
       Time:   proto.Uint32(123),
       Params: proto.String(string(paramsBytes)),
    }
    var events = []*pb_event.Event{event}
    //产生一个连接器。
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
 ```
### APP数据上报
```go
package main
 
 import (
    "code.byted.org/data/mario_collector"
    "code.byted.org/data/mario_collector/pb_event"
    "code.byted.org/gopkg/logs"
    "encoding/json"
    "github.com/golang/protobuf/proto"
 )
 
 func main() {
    user := &pb_event.User{
       UserUniqueId: proto.String("14_13468385056"),
       UserType:     proto.Uint32(14),
       UserId:       proto.Uint64(13468385056),
       UserIsAuth:   proto.Bool(false),
       UserIsLogin:  proto.Bool(false),
       DeviceId:     proto.Uint64(13468385056),
       WebId:        proto.Uint64(12345),
    }
    // 自定义headers，大部分情况下不需要。标准字段请使用预定义字段。
    headers := map[string]interface{}{}
    headers["update_version_code"] = 542
    jsonBytes, _ := json.Marshal(headers)
    header := &pb_event.Header{
       AppId:        proto.Uint32(7), //必选
       AppName:      proto.String("joke"),
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
       UtmSource:    proto.String(""),
       UtmCampaign:  proto.String(""),
       UtmMedium:    proto.String(""),
       ClientIp:     proto.String("10.100.1.1"),
       Headers:      proto.String(string(jsonBytes)),
    }
 
    // event名和params的定义即具体埋点的规范由PM制定
    params := map[string]interface{}{}
    params["enter_from"] = "click_headline"
    params["group_id"] = 123
    paramsBytes, _ := json.Marshal(params)
    event := &pb_event.Event{
       Event:  proto.String("test_go_detail"),
       Time:   proto.Uint32(123),
       Params: proto.String(string(paramsBytes)),
    }
    var events = []*pb_event.Event{event}
    //产生一个连接器。
    mcsCollector := NewAppCollector("e76878c039b9e12be381eb326e247fe2")
	//上报数据
    resp, err := mcsCollector.AppCollectEvents(user, header, events)

	if err == nil {
		defer resp.Body.Close()                        // 保证连接复用
		fmt.Println("response code:", resp.StatusCode) // 查看resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // 查看resp.Body
	}else{
		fmt.Println(err)
	}
 }
```
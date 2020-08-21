# mario_collector
## 简介
服务端GOSDK第一版, APP端当天没有launch事件会导致无法写入离线表。
## SDK已支持功能
- 支持数据单个、批量上报
- 自动补全DeviceId与SSID[已删除]
- 可配置的事件日志
- 异常抛出
 
## 使用方式
### 修改配置信息
```go
//constant.go中，修改IP地址
HttpAddr = "10.225.130.127"
//选择是否记录事件日志以及日志位置
ISLOG = true
LOGPATH = "kkkk.log"
```
### 定义用户属性
```go
    user := &pb_event.User{
       UserUniqueId: proto.String("14_13468385056"),//必选
       UserType:     proto.Uint32(14),
       UserId:       proto.Uint64(13468385056),
       UserIsAuth:   proto.Bool(false),
       UserIsLogin:  proto.Bool(false),
       DeviceId:     proto.Uint64(13468385056), //app必选
       WebId:        proto.Uint64(12345), //web必选
    }

 ```

### 定义header
```go
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
 ```
### 定义事件
```go
    // event名和params的定义即具体埋点的规范由PM制定
    params := map[string]interface{}{}
    params["enter_from"] = "click_headline"
    params["group_id"] = 123
    paramsBytes, _ := json.Marshal(params)
    event := &pb_event.Event{
       Event:  proto.String("EVENT_NAME1"),//必选
       Time:   proto.Uint32(123),
       Params: proto.String(string(paramsBytes)),//必选
    }
    event2 := &pb_event.Event{
       Event:  proto.String("EVENT_NAME2"),//必选
       Time:   proto.Uint32(123),
       Params: proto.String(string(paramsBytes)),//必选
    }
    var events = []*pb_event.Event{event, event2}
 ```

### 移动端的数据上报
```go
    //产生一个连接器。
    client := NewAppCollector()
	//上报数据
    resp, err := client.collect(user, header, events)
    //resp, err := client.collect(user, header, event2) //单个上报
	if err == nil {
		defer resp.Body.Close()                        // 保证连接复用
		fmt.Println("response code:", resp.StatusCode) // 查看resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // 查看resp.Body
	}else{
		fmt.Println(err)
	}
```

### Web端、小程序端的数据上报
```go
    //产生一个连接器。
	client := NewWebMpCollector()
	resp, err := client.collect(user, header, events)
    //  resp, err := client.collect(user, header, event2) //单个数据上报
	if err == nil {
		defer resp.Body.Close()                        // 保证连接复用
		fmt.Println("response code:", resp.StatusCode) // 查看resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body)) // 查看resp.Body
	}else{
		fmt.Println(err)
	}
```

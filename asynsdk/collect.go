package asynsdk

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"sync"
	"time"
)

var (
	mqlxy     *mq
	errlogger *log.Logger
	mu        sync.Mutex
	isFirst   = true
	once      sync.Once
)

func init() {

}

func AppCollect(appid uint32, uuid string, eventname string, eventParam map[string]interface{}) {
	//DCL
	if isFirst {
		mu.Lock()
		if isFirst {
			//厨师日志
			initAsynConf()
			initLogPerHour()
			exep := newExecpool(8)
			mqlxy = newMq()
			go exep.exec()
			isFirst = false
		}
		mu.Unlock()
	}

	appmess := &message{
		isapp:      true,
		appid:      appid,
		uuid:       uuid,
		eventname:  eventname,
		eventParam: eventParam,
	}
	mqlxy.push(appmess)
}

func WebCollect(appid uint32, uuid string, eventname string, eventParam map[string]interface{}) {
	webmess := &message{
		isapp:      false,
		appid:      appid,
		uuid:       uuid,
		eventname:  eventname,
		eventParam: eventParam,
	}
	mqlxy.push(webmess)
}

//每个小时执行.
func initLogPerHour() error {
	if true {
		errlogger = &log.Logger{}
		errwriter, _ := rotatelogs.New(
			singelConfig.A.Errlogpath+".%Y%m%d%H%M",
			rotatelogs.WithRotationTime(time.Duration(60)*time.Minute),
		)
		errlogger.SetOutput(errwriter)
	}
	return nil
}

func logJson(mess *message) (string, error) {
	isapp := mess.isapp
	appid := mess.appid
	uuid := mess.uuid
	eventname := mess.eventname
	eventparam := mess.eventParam
	if isapp {
		var user1 = &user{
			UserUniqueId: proto.String(uuid),
			DeviceId:     proto.Uint64(1),
		}

		var header1 = &header{
			AppId: proto.Uint32(appid), //tob产品使用appkey进行上报，appId设置为0即可
			//Custom: proto.String(""),
			Timezone: proto.Int32(8),
		}

		par, _ := json.Marshal(eventparam)
		var event1 = &event{
			Event:  proto.String(eventname),
			Params: proto.String(string(par)),
		}
		if err := motifyMatchFormatForApp(user1, header1, event1); err != nil {
			return "", err
		}
		ts := uint32(time.Now().Unix())
		events := []*event{event1}
		message := &marioEvents{
			ServerTime: &ts,
			User:       user1,
			Header:     header1,
			AppEvents:  events,
		}
		data, err := json.Marshal(message)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	//Web
	var user1 = &user{
		UserUniqueId: proto.String(uuid),
	}

	var header1 = &header{
		AppId:    proto.Uint32(appid), //tob产品使用appkey进行上报，appId设置为0即可
		Timezone: proto.Int32(8),
	}
	par, _ := json.Marshal(eventparam)
	var event1 = &event{
		Event:  proto.String(eventname),
		Params: proto.String(string(par)),
	}
	ts := uint32(time.Now().Unix())
	events := []*event{event1}
	message := &marioEvents{
		ServerTime: &ts,
		User:       user1,
		Header:     header1,
		Events:     events,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func motifyMatchFormatForApp(user *user, header *header, event *event) error {

	var par = map[string]interface{}{}
	data, ok := event.Params.(*string)
	//ok== false 意味着 空指针
	if ok {
		err := json.Unmarshal([]byte(*data), &par)
		if err == nil {
			event.Params = &par
		} else {
			fmt.Errorf("params 出错")
		}
	}

	data1, ok := event.Params.(string)
	//ok== false 意味着 空指针
	if ok {
		err := json.Unmarshal([]byte(data1), &par)
		if err == nil {
			event.Params = &par
		} else {
			fmt.Errorf("params 出错")
		}
	}
	//根据time字段修改时间
	if event.LocalTimeMs == nil {
		timeObj := time.Unix(time.Now().Unix(), 0)
		if event.Datetime == nil {
			event.Datetime = proto.String(timeObj.Format("2006-01-02 15:04:05"))
		}
		event.LocalTimeMs = proto.Uint64(uint64(time.Now().Unix() * 1000))
		event.Localtime_ms = proto.Uint64(uint64(time.Now().Unix() * 1000))
	} else {
		if event.Time == nil {
			event.Time = proto.Uint32(uint32(time.Now().Unix()))
		}
		timeObj := time.Unix(int64(*event.Time), 0)
		if event.Datetime == nil {
			event.Datetime = proto.String(timeObj.Format("2006-01-02 15:04:05"))
		}
		event.Localtime_ms = event.LocalTimeMs
	}

	if header.AppId == nil {
		err := fmt.Errorf("appid is nil")
		return err
	}
	//2. header的修改,增加 deviceID, 增加aid
	//增加uuID
	header.DeviceId = user.DeviceId
	if header.AppId == nil && header.AppAppId == nil {
		err := fmt.Errorf("appid is nil")
		return err
	}
	if header.AppAppId == nil {
		header.AppAppId = header.AppId
	}
	if header.User_unique_id == nil {
		header.User_unique_id = (user.UserUniqueId)
	}
	return nil
}

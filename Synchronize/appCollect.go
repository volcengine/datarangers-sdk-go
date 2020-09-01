package Synchronize

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
	"time"
)

type appCollector struct {
	mcsCollector
}

func AppCollect(appid uint32, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) (resp *http.Response, err error) {
	defer func() {
		if err != nil {
			data, _ := json.Marshal(&eventParam)
			logmap := map[string]interface{}{
				"appid":      appid,
				"uuid":       uuid,
				"eventname":  eventname,
				"eventparam": string(data),
			}
			data2, _ := json.Marshal(&logmap)
			fmt.Println("[ERROR]" + string(data2) + "err: " + err.Error())

		}
	}()

	if isFirst {
		firstLock.Lock()
		if isFirst {
			if err = initConfig(); err != nil {
				return nil, err
			}
			if appcollector, err = newAppCollector(); err != nil {
				return nil, err
			}
			if webcollector, err = newWebMpCollector(); err != nil {
				return nil, err
			}
			isFirst = false
		}
		firstLock.Unlock()
	}

	var user1 = &user{
		UserUniqueId: proto.String(uuid),
		DeviceId:     proto.Uint64(1),
	}

	var header1 = &header{
		AppId:    proto.Uint32(appid), //tob产品使用appkey进行上报，appId设置为0即可
		Custom:   custom,
		Timezone: proto.Int32(8),
	}

	par, _ := json.Marshal(eventParam)
	var event1 = &event{
		Event:  proto.String(eventname),
		Params: proto.String(string(par)),
	}
	return appcollector.collectEvent(user1, header1, event1)
}

//App上报的接口
func newAppCollector() (*appCollector, error) {

	//mcsurl := "http://" + httpAddr + ":31081" + appURL
	//appcollector = &appCollector{
	//	*newMcsCollector(mcsurl, ""),
	//}
	mcsurl := "http://" + syncConfIns.HttpConfig.HttpAddr + ":31081" + appURL
	appcollector = &appCollector{
		*newMcsCollector(mcsurl, ""),
	}
	return appcollector, nil
}

//事件上报
func (this *appCollector) clsollect(user *user, header *header, ee interface{}) (*http.Response, error) {
	if event, ok := ee.(*event); ok {
		return this.collectEvent(user, header, event)
	}
	if events, ok := ee.([]*event); ok {
		return this.collectEvents(user, header, events)
	}
	return nil, fmt.Errorf("事件类型为[]*pb_event.Event或*pb_event.Event")
}

//单个事件上报
func (this *appCollector) collectEvent(user *user, header *header, event1 *event) (_ *http.Response, err error) {
	events := []*event{}
	return this.collectEvents(user, header, append(events, event1))
}

func (this *appCollector) collectEvents(user *user, header *header, events []*event) (_ *http.Response, err error) {
	caller := ""
	//
	//if user.DeviceId == nil {
	//	return nil, fmt.Errorf("APP上报DeviceId不能为空")
	//}
	//1.修改pras
	//并增加datetime字段
	if err := motifyMatchFormatForApp(user, header, events); err != nil {
		return nil, err
	}

	//5  -> event_v3修改
	//增加launch
	ts := uint32(time.Now().Unix())
	message := &marioEvents{
		Caller:     &caller,
		ServerTime: &ts,
		User:       user,
		Header:     header,
		AppEvents:  events,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", this.mscUrl, strings.NewReader(string(data)))
	//fmt.Println(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X-MCS-AppKey", this.appKey)
	if syncConfIns.EventlogConfig.Islog {
		logger.Println(string(data))
		//println("app")
	}
	resp, err := this.mscHttpClient.Do(req)
	return resp, err
}

func motifyMatchFormatForApp(user *user, header *header, events []*event) error {

	for _, event := range events {
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

package Synchronize

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
	"time"
)

type webCollector struct {
	mcsCollector
}

func WebCollect(appid uint32, uuid string, eventname string, eventParam map[string]interface{}) (*http.Response, error) {
	if isFirst {
		var err error
		//if err = initL(); err != nil {
		//	return nil, err
		//}
		if appcollector, err = newAppCollector(); err != nil {
			return nil, err
		}
		if webcollector, err = newWebMpCollector(); err != nil {
			return nil, err
		}
		isFirst = false
	}

	var user1 = &user{
		UserUniqueId: proto.String(uuid),
	}

	var header1 = &header{
		AppId: proto.Uint32(appid), //tob产品使用appkey进行上报，appId设置为0即可
		//Custom: proto.String(""),
		Timezone: proto.Int32(8),
	}

	par, _ := json.Marshal(eventParam)
	var event1 = &event{
		Event:  proto.String(eventname),
		Params: proto.String(string(par)),
	}
	return webcollector.collectEvent(user1, header1, event1)
}

//Web小程序上报的接口。
func newWebMpCollector() (collector *webCollector, err error) {
	//if err = initL(); err != nil {
	//	return nil, err
	//}
	mcsurl := "http://" + httpAddr + ":31081" + webURL
	collector = &webCollector{
		*newMcsCollector(mcsurl, ""),
	}
	return
}

//事件上报
func (this *webCollector) collect(user *user, header *header, ee interface{}) (*http.Response, error) {

	if event, ok := ee.(*event); ok {
		return this.collectEvent(user, header, event)
	}
	if events, ok := ee.([]*event); ok {
		return this.collectEvents(user, header, events)
	}
	return nil, fmt.Errorf("事件类型为[]*pb_event.Event或*pb_event.Event")
}

//单个事件上报
func (this *webCollector) collectEvent(user *user, header *header, event1 *event) (*http.Response, error) {
	events := []*event{}
	return this.collectEvents(user, header, append(events, event1))
}

//多个事件上报
func (this *webCollector) collectEvents(user *user, header *header, events []*event) (*http.Response, error) {
	caller := ""
	return this.mcsCollectEvents(caller, user, header, events)
}

func (this *mcsCollector) mcsCollectEvents(caller string, user *user, header *header, events []*event) (*http.Response, error) {
	//for _, event := range events {
	//	if event.SessionId == nil {
	//		return nil, fmt.Errorf("SessionID 为空， 不能写入到离线表")
	//	}
	//}
	ts := uint32(time.Now().Unix())
	message := &marioEvents{
		Caller:     &caller,
		ServerTime: &ts,
		User:       user,
		Header:     header,
		Events:     events,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", this.mscUrl, strings.NewReader(string(data)))
	//fmt.Println(string(data))
	if err != nil {
		//logger(err)
		return nil, err
	}
	if iSLOG {
		logger.Println(string(data))
		//println("web")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := this.mscHttpClient.Do(req)
	return resp, err
}

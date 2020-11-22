package datarangers_sdk

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
)

/**
eventParam : 事件属性
custom     : 用户自定义事件公共属性
*/
func SendEvent(apptype1 Apptype, appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if apptype1 != MP && apptype1 != WEB && apptype1 != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, eventname, eventParam, custom, apptype1)
	mqlxy.push(dmg)
	return nil
}

/**
profileAction ：用户公共属性操作类型
profileParam :  用户公共属性
*/
func SendProfile(apptype1 Apptype, appid int64, uuid string, profileAction ProfileActionType, profileParam map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if apptype1 != MP && apptype1 != WEB && apptype1 != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if profileAction != SET && profileAction != SET_ONCE && profileAction != APPEND && profileAction != INCREAMENT && profileAction !=UNSET{
		fatal("请使用正确的profile操作类型")
		return nil
	}
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, string(profileAction), profileParam, map[string]interface{}{}, apptype1)
	mqlxy.push(dmg)
	return nil
}


func ProfileSet(apptype1 Apptype, appid int64, uuid string, profileParam map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if apptype1 != MP && apptype1 != WEB && apptype1 != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, string(SET), profileParam, map[string]interface{}{}, apptype1)
	mqlxy.push(dmg)
	return nil
}


func ProfileSetOnce(apptype1 Apptype, appid int64, uuid string, profileParam map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if apptype1 != MP && apptype1 != WEB && apptype1 != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, string(SET_ONCE), profileParam, map[string]interface{}{}, apptype1)
	mqlxy.push(dmg)
	return nil
}


func ProfileIncrement(apptype1 Apptype, appid int64, uuid string, profileParam map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if apptype1 != MP && apptype1 != WEB && apptype1 != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, string(INCREAMENT), profileParam, map[string]interface{}{}, apptype1)
	mqlxy.push(dmg)
	return nil
}


func ProfileUnset(apptype1 Apptype, appid int64, uuid string, profileParam map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if apptype1 != MP && apptype1 != WEB && apptype1 != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, string(UNSET), profileParam, map[string]interface{}{}, apptype1)
	mqlxy.push(dmg)
	return nil
}


func ProfileAppend(apptype1 Apptype, appid int64, uuid string, profileParam map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if apptype1 != MP && apptype1 != WEB && apptype1 != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, string(APPEND), profileParam, map[string]interface{}{}, apptype1)
	mqlxy.push(dmg)
	return nil
}



func SendEventWithDevice(apptype Apptype, appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, device devicetype, deviceKey string) error {

	if apptype != MP && apptype != WEB && apptype != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if device != ANDROID && device != IOS {
		fatal("deviceType 只能为 ANDROID IOS")
		return nil
	}
	//DCL,初始化MQ,执行池子.
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("初始化携程池完成")
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, eventname, eventParam, custom, apptype)
	//tmp, _ := json.Marshal(dmg)
	//debug("上报的 json 为 -> : " + string(tmp))
	if device == "ANDROID" {
		dmg.Header.Openudid = proto.String(deviceKey)
	} else {
		dmg.Header.Vendor_id = proto.String(deviceKey)
	}
	mqlxy.push(dmg)
	return nil
}

type execpool struct {
	max     int
	tickets chan *ticket
}

type ticket struct {
	id int
}

//单例模式
func newExecpool(x int) *execpool {
	instance = &execpool{}
	instance.max = x
	instance.tickets = make(chan *ticket, instance.max)
	for i := 0; i < instance.max; i++ {
		instance.tickets <- &ticket{id: i}
	}
	return instance
}

func (p *execpool) exec() {
	for {
		t := <-p.tickets
		go func() {
			dmg := mqlxy.pop()
			var err error
			err = appcollector.send(dmg)
			if err != nil {
				ans, _ := json.Marshal(dmg)
				errlogger.Println(string(ans))
			}
			p.tickets <- t
		}()
	}
}

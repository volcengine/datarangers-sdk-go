package event_collect

import (
	"encoding/json"
	"fmt"
)

//import synsdk "code.byted.org/data/datarangers-sdk-go/Synchronize"

/**
producer + consumer
*/
func AppCollectAsyn(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if isFirst {
		firstLock.Lock()
		if isFirst {
			if err := initConfig(); err != nil {
				return err
			}
			initAsyn()
			isFirst = false
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, eventname, eventParam, custom, "app")
	mqlxy.push(dmg)
	return nil
}

func WebCollectAsyn(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if isFirst {
		firstLock.Lock()
		if isFirst {
			if err := initConfig(); err != nil {
				return err
			}
			initAsyn()
			isFirst = false
			firstLock.Unlock()
		}
	}
	dmg := generate(appid, uuid, eventname, eventParam, custom, "web")
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
			if *dmg.App_type == "app" {
				err = appcollector.send(dmg)
			} else {
				err = webcollector.send(dmg)
			}
			if err != nil {
				ans, _ := json.Marshal(dmg)
				//发送失败
				fmt.Println("[ERROR]: " + err.Error())
				errlogger.Println(string(ans))
			}
			p.tickets <- t
		}()
	}
}

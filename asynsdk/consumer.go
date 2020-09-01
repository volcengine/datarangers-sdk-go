package asynsdk

import (
	synsdk "code.byted.org/data/datarangers-sdk-go/Synchronize"
)

var instance *execpool

type execpool struct {
	max     int
	tickets chan *ticket
}

type ticket struct {
	id int
}

//单例模式
func newExecpool(x int) *execpool {
	once.Do(func() {
		instance = &execpool{}
		instance.max = x
		instance.tickets = make(chan *ticket, instance.max)
		for i := 0; i < instance.max; i++ {
			instance.tickets <- &ticket{id: i}
		}
	})
	return instance
}

func sendMessage(mess *message) error {
	isapp := mess.isapp
	appid := mess.appid
	uuid := mess.uuid
	eventname := mess.eventname
	eventparam := mess.eventParam
	if isapp {
		_, err := synsdk.AppCollect(appid, uuid, eventname, eventparam)
		return err
	}
	_, err := synsdk.WebCollect(appid, uuid, eventname, eventparam)
	return err
}

func (p *execpool) exec() {
	for {
		t := <-p.tickets
		go func() {
			item := mqlxy.pop()
			err := sendMessage(item)
			if err != nil {
				ans, _ := logJson(item)
				errlogger.Println(ans)
			}
			p.tickets <- t
		}()
	}
}

func (p *execpool) clnErrLog() {
	for {
		t := <-p.tickets
		go func() {
			item := mqlxy.pop()
			sendMessage(item)
			p.tickets <- t
		}()
	}
}

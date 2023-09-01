package datarangers_sdk

/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// SendEventAb
// eventParam : 事件属性
// custom     : 用户自定义事件公共属性
// Deprecated instead of SendEventInfo
func SendEventAb(apptype AppType, appid int64, uuid string, abSdkVersion string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, did ...int64) error {
	var abSdkVersionList []string
	if abSdkVersion != "" {
		abSdkVersionList = []string{abSdkVersion}
	}
	dmg := getServerSdkEventWithAbMessage(appid, uuid, abSdkVersionList, []string{eventname}, []map[string]interface{}{eventParam}, custom, apptype, did...)
	mq.push(dmg)
	return nil
}

// SendEvent
// eventParam : 事件属性
// custom     : 用户自定义事件公共属性
// Deprecated: instead of SendEventInfo
func SendEvent(appType AppType, appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, did ...int64) error {
	return SendEventAb(appType, appid, uuid, "", eventname, eventParam, custom, did...)
}

// SendEvents
// eventParam : 事件属性
// custom     : 用户自定义事件公共属性
// Deprecated: instead of SendEventInfos
func SendEvents(appType AppType, appid int64, uuid string, eventnameList []string, eventParamList []map[string]interface{}, custom map[string]interface{}, did ...int64) error {
	if len(eventnameList) != len(eventParamList) {
		return fmt.Errorf("事件数目与 属性数目对不上")
	}
	dmg := getServerSdkEventMessage(appid, uuid, eventnameList, eventParamList, custom, appType, did...)
	mq.push(dmg)
	return nil
}

// SendEventInfo 上报事件信息
// event: 事件信息
// custom: 事件公共属性
func SendEventInfo(appType AppType, appId int64, uuid string, event *EventV3, custom map[string]interface{}) error {
	return SendEventInfos(appType, appId, uuid, []*EventV3{event}, custom)
}

// SendEventInfoWithItem 上报携带item的事件
// event: 事件信息
// custom: 事件公共属性
// itemList: item 信息
func SendEventInfoWithItem(appType AppType, appId int64, uuid string, event *EventV3, custom map[string]interface{}, itemList []*Item) error {
	generateItem(event.Params, itemList)
	return SendEventInfos(appType, appId, uuid, []*EventV3{event}, custom)
}

// SendEventInfos 上报多个事件
// events: 事件信息
// custom: 事件公共属性
func SendEventInfos(appType AppType, appId int64, uuid string, events []*EventV3, custom map[string]interface{}) error {
	hd := &Header{
		Aid:          &appId,
		Custom:       custom,
		UserUniqueId: &uuid,
	}
	return SendEventsWithHeader(appType, appId, hd, events)
}

// SendEventWithHeader 使用header上报事件
// header: 事件的header
// event: 事件信息
func SendEventWithHeader(appType AppType, appId int64, hd *Header, event *EventV3) error {
	return SendEventsWithHeader(appType, appId, hd, []*EventV3{event})
}

func SendEventsWithHeader(appType AppType, appId int64, hd *Header, events []*EventV3) error {
	dmg := getEventsWithHeader(appId, appType, hd, events)
	mq.push(dmg)
	return nil
}

// SendProfile
// profileAction ：用户公共属性操作类型
// profileParam :  用户公共属性
func SendProfile(apptype AppType, appid int64, uuid string, profileAction ProfileActionType, profileParam map[string]interface{}, did ...int64) error {
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(profileAction)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	mq.push(dmg)
	return nil
}

// SendItem
// Deprecated: instead of SendEventInfoWithItem
func SendItem(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, itemList []*Item) error {
	generateItem(eventParam, itemList)
	dmg := getServerSdkEventMessage(appid, uuid, []string{eventname}, []map[string]interface{}{eventParam}, custom, APP)
	mq.push(dmg)
	return nil
}

func generateItem(eventParam map[string]interface{}, itemList []*Item) {
	var __items = []interface{}{}
	itemmap := map[string][]interface{}{}
	for _, item := range itemList {
		itemIdmap := map[string]interface{}{}
		itemIdmap["id"] = item.ItemId
		itemmap[*item.ItemName] = append(itemmap[*item.ItemName], itemIdmap)
	}
	__items = append(__items, itemmap)
	eventParam["__items"] = __items
}

// ItemSet 设置item 属性
// itemName: item 名称
// itemParamList: item 属性信息
func ItemSet(appid int64, itemName string, itemParamList []map[string]interface{}) error {
	if ok := checkItemParamList(itemName, itemParamList); !ok {
		return fmt.Errorf("itemParam Must contains Id &&& id must be string")
	}
	//TODO 批量set 失效。
	batch := []string{}
	for i := 0; i < len(itemParamList); i++ {
		batch = append(batch, string(ITEM_SET))
	}
	dmg := getServerSdkEventMessage(appid, "__rangers", batch, itemParamList, map[string]interface{}{}, APP)
	dmg.MessageType = MESSAGE_ITEM
	mq.push(dmg)
	return nil
}

func checkItemParamList(itemName string, itemParamList []map[string]interface{}) bool {
	for _, itemMap := range itemParamList {
		if id, ok := itemMap["id"]; ok {
			if intId, ok := id.(int); ok {
				id = strconv.Itoa(intId)
			}
			if _, ok := id.(string); !ok {
				return false
			}
			itemMap["item_id"] = id
			itemMap["item_name"] = itemName
			delete(itemMap, "id")
		} else {
			return false
		}
	}
	return true
}

func ItemUnset(appid int64, itemName string, id string, removeKeyList []string) error {
	itemParam := map[string]interface{}{}
	itemParam["item_name"] = itemName
	itemParam["item_id"] = id
	for _, key := range removeKeyList {
		itemParam[key] = 1
	}
	dmg := getServerSdkEventMessage(appid, "__rangers", []string{string(ITEM_UNSET)}, []map[string]interface{}{itemParam}, map[string]interface{}{}, APP)
	dmg.MessageType = MESSAGE_ITEM
	mq.push(dmg)
	return nil
}

func ProfileSet(apptype AppType, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(SET)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	dmg.MessageType = MESSAGE_USER
	mq.push(dmg)
	return nil
}

func ProfileSetOnce(apptype AppType, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(SET_ONCE)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	dmg.MessageType = MESSAGE_USER
	mq.push(dmg)
	return nil
}

func ProfileIncrement(apptype AppType, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(INCREAMENT)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	dmg.MessageType = MESSAGE_USER
	mq.push(dmg)
	return nil
}

func ProfileUnset(apptype AppType, appid int64, uuid string, profileNameList []string, did ...int64) error {
	profileParam := map[string]interface{}{}
	for _, name := range profileNameList {
		profileParam[name] = 1
	}
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(UNSET)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	dmg.MessageType = MESSAGE_USER
	mq.push(dmg)
	return nil
}

func ProfileAppend(apptype AppType, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(APPEND)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	dmg.MessageType = MESSAGE_USER
	mq.push(dmg)
	return nil
}

func SendEventWithDevice(apptype AppType, appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, device deviceType, deviceKey string) error {
	dmg := getServerSdkEventMessage(appid, uuid, []string{eventname}, []map[string]interface{}{eventParam}, custom, apptype)
	if device == "ANDROID" {
		dmg.Header.Openudid = &deviceKey
	} else {
		dmg.Header.VendorId = &deviceKey
	}
	mq.push(dmg)
	return nil
}

type execPool struct {
	max     int
	tickets chan *ticket
}

type ticket struct {
	id int
}

//单例模式
func newExecpool(x int) *execPool {
	instance = &execPool{}
	instance.max = x
	instance.tickets = make(chan *ticket, instance.max)
	for i := 0; i < instance.max; i++ {
		instance.tickets <- &ticket{id: i}
	}
	return instance
}

func (p *execPool) exec() {
	for {
		t := <-p.tickets
		go func() {
			p.Send()
			p.tickets <- t
		}()
	}
}

func (p *execPool) Send() {
	if confIns.BatchConfig.Enable {
		dmgs, err := p.doSendBatch()
		if err != nil {
			for _, dmg := range dmgs {
				ans, _ := json.Marshal(dmg)
				errFileWriter.Println(string(ans))
			}
		}
		return
	}
	dmg := mq.pop()
	err := appCollector.send(dmg)
	if err != nil {
		ans, _ := json.Marshal(dmg)
		errFileWriter.Println(string(ans))
	}
}

func (p *execPool) doSend() (interface{}, error) {
	if confIns.BatchConfig.Enable {
		return p.doSendBatch()
	}
	dmg := mq.pop()
	var err error
	err = appCollector.send(dmg)
	return dmg, err
}

func (p *execPool) doSendBatch() ([]interface{}, error) {
	waitTimeMs := time.Duration(confIns.BatchConfig.WaitTimeMs) * time.Millisecond
	dmgs := mq.popBatch(confIns.BatchConfig.Size, waitTimeMs)
	if len(dmgs) < 1 {
		return dmgs, nil
	}
	debug(fmt.Sprintf("dmgs size: %d", len(dmgs)))
	err := appCollector.sendBatch(dmgs)
	return dmgs, err
}

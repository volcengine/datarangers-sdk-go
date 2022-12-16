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
	"github.com/golang/protobuf/proto"
	"strconv"
)

/**
eventParam : 事件属性
custom     : 用户自定义事件公共属性
*/
func SendEventAb(apptype Apptype, appid int64, uuid string, abSdkVersion string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, did ...int64) error {
	//DCL,初始化MQ,执行池子.
	if apptype != MP && apptype != WEB && apptype != APP {
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
	var abSdkVersionList []string
	if abSdkVersion != "" {
		abSdkVersionList = []string{abSdkVersion}
	}
	dmg := getServerSdkEventWithAbMessage(appid, uuid, abSdkVersionList, []string{eventname}, []map[string]interface{}{eventParam}, custom, apptype, did...)
	mqlxy.push(dmg)
	return nil
}

func SendEvent(apptype Apptype, appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, did ...int64) error {
	return SendEventAb(apptype, appid, uuid, "", eventname, eventParam, custom, did...)
}

/**
eventParam : 事件属性
custom     : 用户自定义事件公共属性
*/
func SendEvents(apptype Apptype, appid int64, uuid string, eventnameList []string, eventParamList []map[string]interface{}, custom map[string]interface{}, did ...int64) error {
	//DCL,初始化MQ,执行池子.
	if apptype != MP && apptype != WEB && apptype != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if len(eventnameList) != len(eventParamList) {
		return fmt.Errorf("事件数目与 属性数目对不上")
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
	dmg := getServerSdkEventMessage(appid, uuid, eventnameList, eventParamList, custom, apptype, did...)
	mqlxy.push(dmg)
	return nil
}

/**
profileAction ：用户公共属性操作类型
profileParam :  用户公共属性
*/
func SendProfile(apptype Apptype, appid int64, uuid string, profileAction ProfileActionType, profileParam map[string]interface{}, did ...int64) error {
	//DCL,初始化MQ,执行池子.
	if apptype != MP && apptype != WEB && apptype != APP {
		fatal("apptype 只能为 MP WEB APP")
		return nil
	}
	if profileAction != SET && profileAction != SET_ONCE && profileAction != APPEND && profileAction != INCREAMENT && profileAction != UNSET {
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
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(profileAction)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	mqlxy.push(dmg)
	return nil
}

//
func SendItem(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, itemList []*Item) error {

	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	generateItem(eventParam, itemList)
	dmg := getServerSdkEventMessage(appid, uuid, []string{eventname}, []map[string]interface{}{eventParam}, custom, APP)
	mqlxy.push(dmg)
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

func ItemSet(appid int64, itemName string, itemParamList []map[string]interface{}) error {
	//DCL,初始化MQ,执行池子.
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	if ok := checkItemParamList(itemName, itemParamList); !ok {
		return fmt.Errorf("itemParam Must contains Id &&& id must be string")
	}
	//TODO 批量set 失效。
	batch := []string{}
	for i := 0; i < len(itemParamList); i++ {
		batch = append(batch, "__item_set")
	}
	dmg := getServerSdkEventMessage(appid, "__rangers", batch, itemParamList, map[string]interface{}{}, APP)
	mqlxy.push(dmg)
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
	//DCL,初始化MQ,执行池子.
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	itemParam := map[string]interface{}{}
	itemParam["item_name"] = itemName
	itemParam["item_id"] = id
	for _, key := range removeKeyList {
		itemParam[key] = 1
	}
	dmg := getServerSdkEventMessage(appid, "__rangers", []string{"__item_unset"}, []map[string]interface{}{itemParam}, map[string]interface{}{}, APP)
	mqlxy.push(dmg)
	return nil
}

func ProfileSet(apptype Apptype, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	//DCL,初始化MQ,执行池子.
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(SET)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	mqlxy.push(dmg)
	return nil
}

func ProfileSetOnce(apptype Apptype, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(SET_ONCE)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	mqlxy.push(dmg)
	return nil
}

func ProfileIncrement(apptype Apptype, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(INCREAMENT)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	mqlxy.push(dmg)
	return nil
}

func ProfileUnset(apptype Apptype, appid int64, uuid string, profileNameList []string, did ...int64) error {
	//DCL,初始化MQ,执行池子.

	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	profileParam := map[string]interface{}{}
	for _, name := range profileNameList {
		profileParam[name] = 1
	}
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(UNSET)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
	mqlxy.push(dmg)
	return nil
}

func ProfileAppend(apptype Apptype, appid int64, uuid string, profileParam map[string]interface{}, did ...int64) error {
	if isFirst {
		firstLock.Lock()
		if isFirst {
			initAsyn()
			isFirst = false
			debug("init goroutine pool success")
			firstLock.Unlock()
		}
	}
	dmg := getServerSdkEventMessage(appid, uuid, []string{string(APPEND)}, []map[string]interface{}{profileParam}, map[string]interface{}{}, apptype, did...)
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
	dmg := getServerSdkEventMessage(appid, uuid, []string{eventname}, []map[string]interface{}{eventParam}, custom, apptype)
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

/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package _example

import (
	"fmt"
	sdk "github.com/volcengine/datarangers-sdk-go"
	"time"
)

type SDKExample struct {
	AppId int64
	Did   *int64
	Uuid  string
}

// SendEvent
// Deprecated: instead of SendEventInfo
func (example *SDKExample) SendEvent() {
	appId := example.AppId
	uuid := example.Uuid
	eventParams := map[string]interface{}{
		"param_string1": "value1",
		"param_int1":    12,
	}
	commonParams := map[string]interface{}{
		"common_string1": "common_value1",
		"common_int1":    13,
	}
	sdk.SendEvent(sdk.APP, appId, uuid, "app_event1", eventParams, commonParams)
	sdk.SendEvent(sdk.WEB, appId, uuid, "web_event1", eventParams, commonParams)
	sdk.SendEvent(sdk.MP, appId, uuid, "mp_event1", eventParams, commonParams)

	// 传入数组, eventNameList, eventParamsList 是一一对应的
	eventNameList := []string{
		"app_event2",
		"app_event3",
	}
	eventParamsList := []map[string]interface{}{
		{"event1param": 1},
		{"event2param": 2},
	}
	sdk.SendEvents(sdk.APP, appId, uuid, eventNameList, eventParamsList, commonParams)
}

func (example *SDKExample) SendEventInfo() {
	appId := example.AppId
	uuid := example.Uuid
	// 事件公共属性，如果不需要的化，可以传nil
	commonParams := map[string]interface{}{
		"common_string1": "common_value1",
		"common_int1":    13,
	}
	eventName := "app_send_event_info"
	// 事件发生时间
	localTimeMs := time.Now().UnixMilli()
	eventV3 := &sdk.EventV3{
		Event:       eventName,
		LocalTimeMs: &localTimeMs,
		Params: map[string]interface{}{
			"app_send_event_param1": "value1",
			"app_send_event_param2": "value2",
		},
	}
	sdk.SendEventInfo(sdk.APP, appId, uuid, eventV3, commonParams)
}

func (example *SDKExample) SendEventInfos() {
	appId := example.AppId
	uuid := example.Uuid
	// 事件公共属性，如果不需要的化，可以传nil
	commonParams := map[string]interface{}{
		"common_string1": "common_value1",
		"common_int1":    13,
	}
	event1 := &sdk.EventV3{
		Event: "app_send_event_infos1",
		Params: map[string]interface{}{
			"app_send_event_infos1_param1": "value1",
			"app_send_event_infos1_param2": "value2",
		},
	}
	event2 := &sdk.EventV3{
		Event: "app_send_event_infos2",
		Params: map[string]interface{}{
			"app_send_event_infos2_param1": "value1",
			"app_send_event_infos2_param2": "value2",
		},
	}
	events := []*sdk.EventV3{event1, event2}
	sdk.SendEventInfos(sdk.APP, appId, uuid, events, commonParams)
}

func (example *SDKExample) SetProfile() {
	appId := example.AppId
	uuid := example.Uuid
	properties := map[string]interface{}{
		"list":      []string{"a"},
		"profile_a": "param_11",
	}

	sdk.ProfileSet(sdk.APP, appId, uuid, properties)
}

func (example *SDKExample) SetItem() {
	appId := example.AppId
	itemParamList := []map[string]interface{}{
		{"id": 121, "category": "economics"},
		{"id": 122, "category": "literature"},
		{"id": 123, "category": "fiction"},
		{"id": 124, "category": "fiction"},
		{"id": 125, "category": "fiction"},
	}

	err := sdk.ItemSet(appId, "book", itemParamList)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func (example *SDKExample) SendEventWithItem() {
	appId := example.AppId
	uuid := example.Uuid
	// 事件公共属性，如果不需要的化，可以传nil
	commonParams := map[string]interface{}{
		"common_string1": "common_value1",
		"common_int1":    13,
	}
	eventName := "send_event_with_item"
	// 事件发生时间
	localTimeMs := time.Now().UnixMilli()
	eventV3 := &sdk.EventV3{
		Event:       eventName,
		LocalTimeMs: &localTimeMs,
		Params: map[string]interface{}{
			"app_send_event_param1": "value1",
			"app_send_event_param2": "value2",
		},
	}
	itemList := []*sdk.Item{
		{
			ItemName: sdk.PtrString("book"),
			ItemId:   sdk.PtrString("121"),
		},
		{
			ItemName: sdk.PtrString("book"),
			ItemId:   sdk.PtrString("122"),
		},
	}
	sdk.SendEventInfoWithItem(sdk.APP, appId, uuid, eventV3, commonParams, itemList)
}

func (example *SDKExample) SendEventsWithHeader() {
	appId := example.AppId
	did := example.Did
	uuid := example.Uuid
	// 事件公共属性
	custom := map[string]interface{}{
		"common_param_1": "value_1",
	}
	hd := &sdk.Header{
		Aid:          &appId,
		Custom:       custom,
		DeviceId:     did,
		UserUniqueId: &uuid,
	}
	// 事件发生时间
	localTimeMs1 := time.Now().UnixMilli()
	// 事件属性
	eventParams := map[string]interface{}{
		"param1": "value1",
	}
	event1 := &sdk.EventV3{
		Event:        "go_sdk_send_events_with_header1", //
		LocalTimeMs:  &localTimeMs1,
		AbSdkVersion: &[]string{"11,12"}[0],
		Params:       eventParams,
	}
	localTimeMs2 := time.Now().UnixMilli()
	event2 := &sdk.EventV3{
		Event:        "go_sdk_send_events_with_header2",
		LocalTimeMs:  &localTimeMs2,
		AbSdkVersion: sdk.PtrString("22"),
		Params: map[string]interface{}{
			"param1": "value1",
		},
	}
	events := []*sdk.EventV3{event1, event2}
	sdk.SendEventsWithHeader(sdk.APP, appId, hd, events)
}

func (example *SDKExample) SendEventInfoWithHeader() {
	appId := example.AppId
	did := example.Did
	uuid := example.Uuid
	// 事件公共属性
	custom := map[string]interface{}{
		"common_param_1": "value_1",
	}
	hd := &sdk.Header{
		Aid:          &appId,
		Custom:       custom,
		DeviceId:     did,
		UserUniqueId: &uuid,
	}
	// 事件发生时间
	localTimeMs1 := time.Now().UnixMilli()
	// 事件属性
	eventParams := map[string]interface{}{
		"param1": "value1",
	}
	eventV3 := &sdk.EventV3{
		Event:        "go_sdk_send_with_header", //
		LocalTimeMs:  &localTimeMs1,
		AbSdkVersion: &[]string{"11,12"}[0],
		Params:       eventParams,
	}
	sdk.SendEventWithHeader(sdk.APP, appId, hd, eventV3)
}

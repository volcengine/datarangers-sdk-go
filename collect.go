package datarangers_sdk

import "net/http"

/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"net"
	"strconv"
	"strings"
	"time"
)

const PATH_SDK_LOG = "/sdk/log"
const PATH_SDK_LIST = "/sdk/list"
const PATH_EVENT_JSON = "/v2/event/json"
const PATH_EVENT_LIST = "/v2/event/list"
const PATH_USER_ATTR = "/dataprofile/openapi/v1/%d/users/%s/attributes"
const PATH_ITEM_ATTR = "/dataprofile/openapi/v1/%d/items/%s/%s/attributes"

type mcsCollector struct {
	mscHttpClient     *http.Client
	syncKafkaProducer *SyncKafkaProducer
}

func (m mcsCollector) close() {
	if confIns.SdkConfig.Mode == MODE_KAFKA && m.syncKafkaProducer != nil {
		err := m.syncKafkaProducer.Close()
		if err != nil {
			fatal(fmt.Sprintf("close syncKafkaProducer failed, error: %v", err))
			return
		}
		info("syncKafkaProducer closed")
	}
}

func newMcsCollector() (collector *mcsCollector) {
	var syncKafkaProducer *SyncKafkaProducer
	var mscHttpClient *http.Client
	var err error
	if confIns.SdkConfig.Mode == MODE_KAFKA {
		syncKafkaProducer, err = NewSyncKafkaProducer(confIns.KafkaConfig)
		if err != nil {
			panic("init kafka producer failed" + err.Error())
		}
	} else {
		mscHttpClient = &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, time.Duration(confIns.HttpConfig.SocketTimeOut)*time.Second)
				},
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
				MaxIdleConns:        maxIdleConnsPerHost,
			},
			Timeout: time.Duration(confIns.HttpConfig.SocketTimeOut) * time.Second,
		}
	}

	collector = &mcsCollector{
		mscHttpClient:     mscHttpClient,
		syncKafkaProducer: syncKafkaProducer,
	}
	return
}

//sendBatch
// 批量发送
func (p *mcsCollector) sendBatch(dmsgs []interface{}) error {
	var err error
	// File 模式
	if MODE_FILE == confIns.SdkConfig.Mode {
		return fmt.Errorf("batch not support mode: %s", confIns.SdkConfig.Mode)
	} else if MODE_HTTP == confIns.SdkConfig.Mode {
		switch confIns.SdkConfig.Env {
		case ENV_PRI:
			err = p.priSendEventBatch(dmsgs)
			break
		case ENV_SAAS_NATIVE:
			err = p.saasNativeSendEventBatch(dmsgs)
			break
		default:
			fatal(fmt.Sprintf("not support env: %s", confIns.SdkConfig.Env))
			return fmt.Errorf("not support env: %s", confIns.SdkConfig.Env)
		}

		return err
	} else if MODE_KAFKA == confIns.SdkConfig.Mode {
		// kafka模式，只支持私部环境
		switch confIns.SdkConfig.Env {
		case ENV_PRI:
			err = p.priSendEventKafkaBatch(dmsgs)
			break
		default:
			fatal(fmt.Sprintf("not support env: %s", confIns.SdkConfig.Env))
			return fmt.Errorf("not support env: %s", confIns.SdkConfig.Env)
		}

		return err
	}
	//其余的不支持
	fatal(fmt.Sprintf("batch not support mode: %s", confIns.SdkConfig.Mode))
	return fmt.Errorf("btach not support mode: %s", confIns.SdkConfig.Mode)
}

func (p *mcsCollector) send(dmsg interface{}) error {
	var err error
	// File 模式
	if MODE_FILE == confIns.SdkConfig.Mode {
		data, err := json.Marshal(dmsg)
		if err != nil {
			fatal(err.Error())
			return err
		}
		debug("save success！saved json: " + string(data))
		fileWriter.Println(string(data))
		return nil
	} else if MODE_HTTP == confIns.SdkConfig.Mode {
		switch confIns.SdkConfig.Env {
		case ENV_PRI:
			err = p.priSendEvent(dmsg)
			break
		case ENV_SAAS:
			err = p.saasSend(dmsg.(*ServerSdkEventMessage))
			break
		case ENV_SAAS_NATIVE:
			err = p.saasNativeSendEvent(dmsg.(*ServerSdkEventMessage))
			break
		default:
			fatal(fmt.Sprintf("not support env: %s", confIns.SdkConfig.Env))
			return fmt.Errorf("not support env: %s", confIns.SdkConfig.Env)
		}

		return err
	} else if MODE_KAFKA == confIns.SdkConfig.Mode {
		// kafka模式，只支持私部环境
		switch confIns.SdkConfig.Env {
		case ENV_PRI:
			err = p.priSendEventKafka(dmsg)
			break
		default:
			fatal(fmt.Sprintf("not support env: %s", confIns.SdkConfig.Env))
			return fmt.Errorf("not support env: %s", confIns.SdkConfig.Env)
		}

		return err
	}
	//其余的不支持
	fatal(fmt.Sprintf("not support mode: %s", confIns.SdkConfig.Mode))
	return fmt.Errorf("not support mode: %s", confIns.SdkConfig.Mode)
}

func (p *mcsCollector) priSendEvent(dmsg interface{}) error {
	data, err := json.Marshal(dmsg)
	if err != nil {
		fatal(err.Error())
		return err
	}
	url := confIns.HttpConfig.HttpAddr + PATH_SDK_LOG
	return p.request("POST", url, data, nil)
}

func (p *mcsCollector) priSendEventBatch(dmsg []interface{}) error {
	data, err := json.Marshal(dmsg)
	if err != nil {
		fatal(err.Error())
		return err
	}
	url := confIns.HttpConfig.HttpAddr + PATH_SDK_LIST
	return p.request("POST", url, data, nil)
}

func (p *mcsCollector) priSendEventKafkaBatch(dmsg []interface{}) error {
	return p.syncKafkaProducer.BatchSend(dmsg)
}

func (p *mcsCollector) priSendEventKafka(dmsg interface{}) error {
	return p.syncKafkaProducer.Send(dmsg)
}

func (p *mcsCollector) saasNativeSendEvent(dmsg *ServerSdkEventMessage) error {
	appKey, ok := confIns.AppKeys[*dmsg.AppId]
	if !ok {
		panic("App key cannot be empty. appId: " + strconv.Itoa(int(*dmsg.AppId)))
	}
	snam := CreateSaasNativeAppMessage(*dmsg)
	data, err := json.Marshal(snam)
	if err != nil {
		fatal(err.Error())
		return err
	}
	customHeaders := map[string]string{
		APP_KEY: appKey,
	}
	url := confIns.HttpConfig.HttpAddr + PATH_SDK_LOG
	return p.request("POST", url, data, customHeaders)
}

func (p *mcsCollector) saasNativeSendEventBatch(dmsgs []interface{}) error {
	// 根据appId 进行group by
	ssemMap := make(map[int64][]*ServerSdkEventMessage)
	for _, dmsg := range dmsgs {
		ssem := dmsg.(*ServerSdkEventMessage)
		ssemList, ok := ssemMap[*ssem.AppId]
		if !ok {
			ssemList = make([]*ServerSdkEventMessage, 0)
		}
		ssemList = append(ssemList, ssem)
		ssemMap[*ssem.AppId] = ssemList
	}

	// 处理
	for _, ssemList := range ssemMap {
		err := p.doSaasNativeSendEventBatch(ssemList)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *mcsCollector) doSaasNativeSendEventBatch(dmsgs []*ServerSdkEventMessage) error {
	dmsg := dmsgs[0]
	appKey, ok := confIns.AppKeys[*dmsg.AppId]
	if !ok {
		panic("App key cannot be empty. appId: " + strconv.Itoa(int(*dmsg.AppId)))
	}
	snams := CreateSaasNativeAppMessages(dmsgs)
	data, err := json.Marshal(snams)
	if err != nil {
		fatal(err.Error())
		return err
	}
	customHeaders := map[string]string{
		APP_KEY: appKey,
	}
	url := confIns.HttpConfig.HttpAddr + PATH_SDK_LIST
	return p.request("POST", url, data, customHeaders)
}

func (p *mcsCollector) saasSend(dmsg *ServerSdkEventMessage) error {
	// 判断类型

	err := errors.New("server error")
	switch dmsg.MessageType {
	case MESSAGE_USER:
		err = p.saasSendUserProfile(dmsg)
		break
	case MESSAGE_ITEM:
		err = p.saasSendItemProfile(dmsg)
		break
	default:
		err = p.saasSendEvent(dmsg)
	}
	return err
}

// 上报事件
func (p *mcsCollector) saasSendEvent(dmsg *ServerSdkEventMessage) error {
	appKey, ok := confIns.AppKeys[*dmsg.AppId]
	if !ok {
		panic("App key cannot be empty. appId: " + strconv.Itoa(int(*dmsg.AppId)))
	}
	ssam := CreateSaasServerAppMessage(dmsg)
	data, err := json.Marshal(ssam)
	if err != nil {
		fatal(err.Error())
		return err
	}
	url := confIns.HttpConfig.HttpAddr + PATH_EVENT_JSON
	customHeaders := map[string]string{
		APP_KEY: appKey,
	}
	return p.request("POST", url, data, customHeaders)
}

// 上报用户属性
func (p *mcsCollector) saasSendUserProfile(dmsg *ServerSdkEventMessage) error {
	spam := CreateSaasProfileAppMessage(dmsg)
	data, err := json.Marshal(spam)
	if err != nil {
		fatal(err.Error())
		return err
	}
	uriPath := fmt.Sprintf(PATH_USER_ATTR, *dmsg.AppId, *dmsg.UserUniqueId)
	url := confIns.OpenapiConfig.HttpAddr + uriPath
	method := "PUT"
	authorization := Sign(confIns.OpenapiConfig.Ak, confIns.OpenapiConfig.Sk, 1800, method, uriPath, nil, string(data))
	customHeaders := map[string]string{
		AUTHORIZATION: authorization,
	}
	return p.request(method, url, data, customHeaders)
}

// 上报item
func (p *mcsCollector) saasSendItemProfile(dmsg *ServerSdkEventMessage) error {
	if dmsg.EventV3 == nil {
		return nil
	}
	var err error
	for _, eventV3 := range dmsg.EventV3 {
		siam := CreateSaasItemAppMessage(eventV3)
		data, err := json.Marshal(siam)
		if err != nil {
			fatal(err.Error())
			return err
		}
		mParams := eventV3.Params
		itemName, _ := mParams["item_name"]
		itemId, _ := mParams["item_id"]
		uriPath := fmt.Sprintf(PATH_ITEM_ATTR, *dmsg.AppId, itemName, itemId)
		url := confIns.OpenapiConfig.HttpAddr + uriPath
		method := "PUT"
		authorization := Sign(confIns.OpenapiConfig.Ak, confIns.OpenapiConfig.Sk, 1800, method, uriPath, nil, string(data))
		customHeaders := map[string]string{
			AUTHORIZATION: authorization,
		}
		err = p.request(method, url, data, customHeaders)
		if err != nil {
			fatal(err.Error())
			return err
		}
	}
	return err
}

func (p *mcsCollector) request(method string, url string, data []byte, customHeaders map[string]string) error {
	req, err := http.NewRequest(method, url, strings.NewReader(string(data)))
	if err != nil {
		fatal(err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "DataRangers Golang SDK")
	defer func() {
		req.Body.Close()
	}()
	//补充 Header
	for k, v := range headers {
		if m, okk := v.(string); okk {
			req.Header.Set(k, m)
			if strings.ToLower(k) == "host" {
				req.Host = m
			}
		}
		if m, okk := v.(int); okk {
			req.Header.Set(k, strconv.Itoa(m))
		}
	}
	// 特殊的需要设置的header
	if customHeaders != nil {
		for k, v := range customHeaders {
			req.Header.Set(k, v)
		}
	}
	h, _ := json.Marshal(req.Header)
	debug("request Header : " + string(h))
	var resp *http.Response
	resp, err = p.mscHttpClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fatal(err.Error() + "    send error, will retry again")
		resp, err = p.mscHttpClient.Do(req)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			fatal(err.Error() + "    retry send error ")
		}
	} else {
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			rBody := fmt.Sprintf("send error, method: %s, url: %s, status code: %d, resp: %s, req:\r\n%s",
				method, url, resp.StatusCode, string(body), string(data))
			fmt.Println(rBody)
			fatal(rBody)
			return fmt.Errorf(rBody)
		} else {
			if resp != nil && resp.Body != nil {
				body, _ := ioutil.ReadAll(resp.Body)
				responseMsg := map[string]interface{}{}
				err2 := json.Unmarshal(body, &responseMsg)
				if err2 != nil {
					rBody := fmt.Sprintf("parse body error, method: %s, url: %s, status code: %d, resp: %s, req:\r\n%s",
						method, url, resp.StatusCode, string(body), string(data))
					fatal(rBody)
					return fmt.Errorf(rBody)
				}
				if msg, ok := responseMsg["message"]; ok && msg == "success" {
					debug(fmt.Sprintf("send success！message json : %s, resp: %s", string(data), string(body)))
				} else {
					rBody := fmt.Sprintf("send error, method: %s, url: %s, status code: %d, resp: %s, req:\r\n%s",
						method, url, resp.StatusCode, string(body), string(data))
					fatal(rBody)
					return fmt.Errorf(rBody)
				}
			}
		}
	}
	return err
}

func getServerSdkEventWithAbMessage(appid int64, uuid string, abSdkVersionList []string, eventnameList []string, eventParam []map[string]interface{}, custom map[string]interface{}, appType AppType, device_id ...int64) *ServerSdkEventMessage {
	if !isInit {
		panic("Sdk must be init first")
	}
	var webid int64
	if len(device_id) != 0 {
		webid = device_id[0]
	}
	if custom == nil {
		custom = map[string]interface{}{}
	}
	custom["__sdk_platform"] = SDK_VERSION
	hd := &Header{
		Aid:          &appid,
		Custom:       custom,
		DeviceId:     &webid,
		UserUniqueId: &uuid,
		Timezone:     &timezone,
	}
	if confIns.SdkConfig.Mode == MODE_KAFKA {
		hd.Source = PtrString(SOURCE__KAFKA_SDK_SERVER)
	}

	appTypeStr := string(appType)
	dmg := &ServerSdkEventMessage{
		AppId:        &appid,
		UserUniqueId: &uuid,
		AppType:      &appTypeStr,
		DeviceId:     &webid,
	}
	timeObj := time.Unix(time.Now().Unix(), 0)
	var sendEventV3 []*EventV3
	for i, eventname := range eventnameList {
		itm := &EventV3{
			Datetime:    PtrString(timeObj.Format(DATE_TIME_LAYOUT)),
			Event:       eventname,
			LocalTimeMs: PtrInt64(time.Now().UnixMilli()),
			Params:      eventParam[i],
		}
		if abSdkVersionList != nil {
			itm.AbSdkVersion = PtrString(abSdkVersionList[i])
		}
		sendEventV3 = append(sendEventV3, itm)
	}
	dmg.EventV3 = sendEventV3
	dmg.Header = hd
	return dmg
}

func getEventsWithHeader(appId int64, appType AppType, hd *Header, events []*EventV3) *ServerSdkEventMessage {
	if !isInit {
		panic("Sdk must be init first")
	}
	var custom = hd.Custom
	if custom == nil {
		custom = map[string]interface{}{}
	}
	if hd.Timezone == nil {
		hd.Timezone = &timezone
	}
	if confIns.SdkConfig.Mode == MODE_KAFKA {
		hd.Source = PtrString(SOURCE__KAFKA_SDK_SERVER)
	}
	for _, eventV3 := range events {
		if eventV3.LocalTimeMs == nil {
			eventV3.LocalTimeMs = PtrInt64(time.Now().UnixMilli())
			eventV3.Datetime = PtrString(time.UnixMilli(*eventV3.LocalTimeMs).Format(DATE_TIME_LAYOUT))
		}
		if eventV3.Datetime == nil {
			eventV3.Datetime = PtrString(time.UnixMilli(*eventV3.LocalTimeMs).Format(DATE_TIME_LAYOUT))
		}
	}
	appTypeStr := string(appType)
	custom["__sdk_platform"] = SDK_VERSION
	dmg := &ServerSdkEventMessage{
		AppId:        &appId,
		UserUniqueId: hd.UserUniqueId,
		AppType:      &appTypeStr,
		DeviceId:     hd.DeviceId,
		EventV3:      events,
		Header:       hd,
	}
	return dmg
}

func getServerSdkEventMessage(appid int64, uuid string, eventnameList []string, eventParam []map[string]interface{}, custom map[string]interface{}, apptype1 AppType, device_id ...int64) *ServerSdkEventMessage {
	return getServerSdkEventWithAbMessage(appid, uuid, nil, eventnameList, eventParam, custom, apptype1, device_id...)
}

func getTimezone() int32 {
	utc, _ := strconv.Atoi(time.Now().UTC().Format(DATE_TIME_LAYOUT)[11:13])
	cur, _ := strconv.Atoi(time.Now().Format(DATE_TIME_LAYOUT)[11:13])
	ans := cur - utc
	if ans <= -12 {
		ans = ans + 24
	}
	return int32(ans)
}

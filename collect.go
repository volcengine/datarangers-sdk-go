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
	"io/ioutil"

	//"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type mcsCollector struct {
	mscUrl        string
	mscHttpClient *http.Client
}

func newMcsCollector(mcsUrl string) (collector *mcsCollector) {
	collector = &mcsCollector{
		mscUrl: mcsUrl,
		mscHttpClient: &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, time.Duration(confIns.HttpConfig.SocketTimeOut)*time.Second)
				},
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
				MaxIdleConns:        maxIdleConnsPerHost,
			},
			Timeout: time.Duration(confIns.HttpConfig.SocketTimeOut) * time.Second,
		},
	}
	return
}

func (p *mcsCollector) send(dmsg interface{}) error {
	data, err := json.Marshal(dmsg)
	if err != nil {
		fatal(err.Error())
		return err
	}
	req, err := http.NewRequest("POST", p.mscUrl, strings.NewReader(string(data)))
	if err != nil {
		fatal(err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	defer func() {
		req.Body.Close()
	}()
	//补充 header
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
	h, _ := json.Marshal(req.Header)
	debug("request header : " + string(h))

	tmp, _ := json.Marshal(dmsg)
	//
	if !confIns.EventlogConfig.EventSendEnable {
		debug("保存成功！保存的json 为 -> : " + string(tmp))
		logger.Println(string(data))
		return nil
	}
	//其余的都上报。
	var resp *http.Response
	resp, err = p.mscHttpClient.Do(req)
	if err != nil {
		warn(err.Error() + "    消息未发送成功, 重试一次")
		resp, err = p.mscHttpClient.Do(req)
		defer resp.Body.Close()
		if err != nil {
			warn(err.Error() + "    重试时 未发送成功 ")
		}
	} else {
		if resp.StatusCode != 200 {
			fatal("信息发送失败，错误码为: " + strconv.Itoa(resp.StatusCode))
			fmt.Println(resp.Body)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
			return fmt.Errorf("信息发送失败，错误码为: " + strconv.Itoa(resp.StatusCode))
		} else {
			if resp != nil && resp.Body != nil {
				body, _ := ioutil.ReadAll(resp.Body)
				responseMsg := map[string]interface{}{}
				err2 := json.Unmarshal(body, &responseMsg)
				if err2 != nil {
					fatal("解析body出错, body 为 " + string(body))
					return fmt.Errorf("解析body出错")
				}
				if msg, ok := responseMsg["message"]; ok && msg == "success" {
					debug("上报成功！上报的json 为 -> : " + string(tmp))
				} else {
					fatal("数据未上报到applog, 返回body为" + string(body))
					return fmt.Errorf("数据未上报到applog")
				}
				resp.Body.Close()
			}
		}
	}
	if resp != nil && resp.Body != nil {
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	return err
}

//func collectsync(apptype Apptype, appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) error {
//	dmg := generate(appid, uuid, eventname, eventParam, custom, apptype)
//	if err := appcollector.send(dmg); err != nil {
//		data, a := json.Marshal(dmg)
//		if a != nil {
//			warn(err.Error() + "消息未发送成功")
//			return a
//		}
//		errlogger.Println(string(data))
//		return err
//	}
//	return nil
//}


func getServerSdkEventMessage(appid int64, uuid string, eventnameList []string, eventParam []map[string]interface{}, custom map[string]interface{}, apptype1 Apptype) *ServerSdkEventMessage {
	hd := &Header{
		Aid:            proto.Int64(appid),
		Custom:         custom,
		Device_id:      proto.Int64(0),
		User_unique_id: proto.String(uuid),
		Timezone:       proto.Int32(int32(timezone)),
	}
	dmg := &ServerSdkEventMessage{
		App_id:         proto.Int64(appid),
		User_unique_id: proto.String(uuid),
		App_type:       proto.String(string(apptype1)),
		Device_id:      proto.Int64(0),
	}
	timeObj := time.Unix(time.Now().Unix(), 0)
	var sendEventV3 []*Event_v3
	for i, eventname := range eventnameList{
		itm := &Event_v3{
			Datetime:    proto.String(timeObj.Format("2006-01-02 15:04:05")),
			Event:       proto.String(eventname),
			LocalTimeMs: proto.Int64(time.Now().UnixNano() / 1e6),
			Params:      eventParam[i],
		}
		sendEventV3 = append(sendEventV3, itm)
	}
	dmg.Event_v3 = sendEventV3
	dmg.Header = hd
	return dmg
}

func getTimezone() int {
	utc, _ := strconv.Atoi(time.Now().UTC().Format("2006-01-02 15:04:05")[11:13])
	cur, _ := strconv.Atoi(time.Now().Format("2006-01-02 15:04:05")[11:13])
	ans := cur - utc
	if ans <= -12 {
		ans = ans + 24
	}
	return ans
}

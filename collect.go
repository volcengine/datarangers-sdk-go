package datarangers_sdk

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

func (p *mcsCollector) send(dmsg *dancemsg) error {
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
	//req.Header.Add("Host", "snssdk.vpc.com")
	//req.Host = "snssdk.vpc.com"
	//defer req.Body.Close()
	defer func() {
		//ioutil.ReadAll(req.Body)
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
	//fmt.Println(confIns.EventlogConfig.Islog)
	if confIns.EventlogConfig.Islog {
		//fmt.Println("保存成功！保存的json 为 -> : " + string(tmp))
		debug("保存成功！保存的json 为 -> : " + string(tmp))
		logger.Println(string(data))
	}
	var resp *http.Response
	if confIns.EventlogConfig.Iscollect {
		resp, err = p.mscHttpClient.Do(req)
		if err != nil {
			warn(err.Error() + "    消息未发送成功, 重试一次")
			resp, err = p.mscHttpClient.Do(req)
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
					a := map[string]interface{}{}
					err2 := json.Unmarshal(body, &a)
					if err2 != nil {
						fatal("解析body出错 ")
						return fmt.Errorf("解析body出错")
					}
					if msg, ok := a["message"]; ok && msg == "success" {
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
		//resp.Body.Close()
	}
	return err
}

func collectsync(apptype Apptype, appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) error {
	//if err := firstInit(); err != nil {
	//	return err
	//}
	dmg := generate(appid, uuid, eventname, eventParam, custom, apptype)
	if err := appcollector.send(dmg); err != nil {
		data, a := json.Marshal(dmg)
		if a != nil {
			warn(err.Error() + "消息未发送成功")
			return a
		}
		errlogger.Println(string(data))
		return err
	}
	return nil
}


func generate(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, apptype1 Apptype) *dancemsg {
	hd := &header{
		Aid:            proto.Int64(appid),
		Custom:         custom,
		Device_id:      proto.Int64(1),
		User_unique_id: proto.String(uuid),
		Timezone:       proto.Int(timezone),
	}

	timeObj := time.Unix(time.Now().Unix(), 0)
	itm := &items{
		Datetime:    proto.String(timeObj.Format("2006-01-02 15:04:05")),
		Event:       proto.String(eventname),
		LocalTimeMs: proto.Int64(time.Now().UnixNano() / 1e6),
		Params:      eventParam,
	}
	//time := &timeSync{
	//	Local_time: proto.Int64(1),
	//	Server_time: proto.Int64(1),
	//}
	dmg := &dancemsg{
		App_id:         proto.Int64(appid),
		User_unique_id: proto.String(uuid),
		App_type:       proto.String(string(apptype1)),
		Device_id:      proto.Int64(1),
	}
	dmg.Event_v3 = []*items{itm}
	dmg.Header = hd
	//dmg.TimeSync = time
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

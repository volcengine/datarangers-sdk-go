package event_collect

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
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
					return net.DialTimeout(network, addr, dialTimeout)
				},
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
			},
			Timeout: totalTimeout,
		},
	}
	return
}

func (p *mcsCollector) send(dmsg *dancemsg) error {
	data, err := json.Marshal(dmsg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", p.mscUrl, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	//补充 header
	for k, v := range headers {
		if m, okk := v.(string); okk {
			req.Header.Set(k, m)
		}
		if m, okk := v.(int); okk {
			req.Header.Set(k, strconv.Itoa(m))
		}
	}
	if confIns.EventlogConfig.Islog {
		logger.Println(string(data))
	}
	fmt.Println(string(data))
	_, err = p.mscHttpClient.Do(req)
	return err
}

func AppCollect(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) error {
	if err := firstInit(); err != nil {
		return err
	}
	dmg := generate(appid, uuid, eventname, eventParam, custom, "app")
	if err := appcollector.send(dmg); err != nil {
		data, a := json.Marshal(dmg)
		if a != nil {
			return a
		}
		errlogger.Println(string(data))
		return err
	}
	return nil
}

func WebCollect(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}) error {
	if err := firstInit(); err != nil {
		return err
	}
	dmg := generate(appid, uuid, eventname, eventParam, custom, "web")
	if err := webcollector.send(dmg); err != nil {
		data, a := json.Marshal(dmg)
		if a != nil {
			return a
		}
		errlogger.Println(string(data))
		return err
	}
	return nil
}

func firstInit() error {
	if isFirst {
		firstLock.Lock()
		if isFirst {
			if err := initConfig(); err != nil {
				return err
			}
			isFirst = false
			firstLock.Unlock()
		}
	}
	return nil
}

func generate(appid int64, uuid string, eventname string, eventParam map[string]interface{}, custom map[string]interface{}, apptype string) *dancemsg {
	hd := &header{
		Aid:            proto.Int64(appid),
		Custom:         custom,
		Device_id:      proto.Int64(1),
		User_unique_id: proto.String(uuid),
	}

	timeObj := time.Unix(time.Now().Unix(), 0)
	itm := &items{
		Datetime:    proto.String(timeObj.Format("2006-01-02 15:04:05")),
		Event:       proto.String(eventname),
		LocalTimeMs: proto.Int64(time.Now().Unix() * 1000),
		Params:      eventParam,
	}
	//time := &timeSync{
	//	Local_time: proto.Int64(1),
	//	Server_time: proto.Int64(1),
	//}
	dmg := &dancemsg{
		App_id:         proto.Int64(appid),
		User_unique_id: proto.String(uuid),
		App_type:       proto.String(apptype),
		Device_id:      proto.Int64(1),
	}
	dmg.Event_v3 = []*items{itm}
	dmg.Header = hd
	//dmg.TimeSync = time
	return dmg
}

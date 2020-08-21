package mario_collector

import (
	"code.byted.org/data/datarangers-sdk-go/pb_event"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	dialTimeout         = 1 * time.Second
	totalTimeout        = 2 * time.Second
	maxIdleConnsPerHost = 4

	logger *log.Logger
)

func initLog() error {
	if ISLOG {
		logFile, err := os.OpenFile(LOGPATH, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err!=nil { return err}
		logger = &log.Logger{}
		logger.SetOutput(logFile)
	}
	return nil
}

type McsCollector struct {
	mscUrl        string
	appKey        string
	mscHttpClient *http.Client
}

type AppCollector struct {
	McsCollector
}

type WebCollector struct {
	McsCollector
}

func NewMcsCollector(mcsUrl string, appKey string) (collector *McsCollector) {
	collector = &McsCollector{
		mscUrl: mcsUrl,
		appKey: appKey,
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

//App上报的接口
func NewAppCollector() (collector *AppCollector, err error) {
	if err = initLog(); err!=nil{ return nil, err}
	mcsurl := "http://"+HttpAddr+":31081"+AppURL
	collector = &AppCollector{
		*NewMcsCollector(mcsurl,""),
	}
	return
}

//Web小程序上报的接口。
func NewWebMpCollector()(collector *WebCollector, err error){
	if err = initLog(); err!=nil{ return nil, err}
	mcsurl := "http://"+HttpAddr+":31081"+WebURL
	collector = &WebCollector{
		*NewMcsCollector(mcsurl,""),
	}
	return
}
//事件上报
func (this *WebCollector) collect(user *pb_event.User, header *pb_event.Header, ee interface{}) (*http.Response, error) {

	if event,ok:=ee.(*pb_event.Event); ok{
		return this.collectEvent(user, header, event)
	}
	if events,ok:=ee.([]*pb_event.Event); ok{
		return this.collectEvents(user, header, events)
	}
	return nil,fmt.Errorf("事件类型为[]*pb_event.Event或*pb_event.Event")
}

//单个事件上报
func (this *WebCollector) collectEvent(user *pb_event.User, header *pb_event.Header, event *pb_event.Event) (*http.Response, error) {
	events:=[]*pb_event.Event{}
	return this.collectEvents(user, header, append(events, event))
}
//多个事件上报
func (this *WebCollector) collectEvents(user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) (*http.Response, error) {
	caller :=""
	return this.McsCollectEvents(caller, user, header, events);
}


//事件上报
func (this *AppCollector) collect(user *pb_event.User, header *pb_event.Header, ee interface{}) (*http.Response, error) {
	if event,ok:=ee.(*pb_event.Event); ok{
		return this.collectEvent(user, header, event)
	}
	if events,ok:=ee.([]*pb_event.Event); ok{
		return this.collectEvents(user, header, events)
	}
	return nil,fmt.Errorf("事件类型为[]*pb_event.Event或*pb_event.Event")
}
//单个事件上报
func (this *AppCollector) collectEvent( user *pb_event.User, header *pb_event.Header, event *pb_event.Event) (_ *http.Response, err error) {
	events:=[]*pb_event.Event{}
	return this.collectEvents(user, header, append(events, event))
}


func (this *AppCollector) collectEvents( user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) (_ *http.Response, err error) {
	caller:= ""
	//
	if user.DeviceId==nil {
		return nil, fmt.Errorf("APP上报DeviceId不能为空")
	}
	//1.修改pras
	//并增加datetime字段
	if err:=MotifyMatchFormatForApp(user, header,events); err!=nil{
		return nil, err
	}

	//5  -> event_v3修改
	//增加launch
	ts := uint32(time.Now().Unix())
	message := &pb_event.MarioEvents{
		Caller:     &caller,
		ServerTime: &ts,
		User:       user,
		Header:     header,
		AppEvents:  events,
	}
	data, err := json.Marshal(message)
	if err != nil {
		//logger(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", this.mscUrl, strings.NewReader(string(data)))
	fmt.Println(strings.NewReader(string(data)))
	if err != nil {
		//logger(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X-MCS-AppKey", this.appKey)
	if ISLOG { logger.Println(string(data))}
	resp, err := this.mscHttpClient.Do(req)
	return resp, err
}


func MotifyMatchFormatForApp(user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) error{

	for _, event:=range events{
		if event.SessionId == nil{
			return fmt.Errorf("SessionID 为空， 不能写入到离线表")
		}
		var par = map[string]interface{}{}
		data,ok := event.Params.(*string)
		//ok== false 意味着 空指针
		if ok{
			err:=json.Unmarshal([]byte(*data), &par)
			if err==nil{
				event.Params = &par;
			}else {
				fmt.Errorf("params 出错")
			}
		}
		//根据time字段修改时间
		if event.LocalTimeMs == nil {
			timeObj := time.Unix(time.Now().Unix(), 0)
			if event.Datetime == nil {
				event.Datetime = proto.String(timeObj.Format("2006-01-02 15:04:05"));
			}
			event.LocalTimeMs = proto.Uint64(uint64(time.Now().Unix()*1000))
			event.Localtime_ms = proto.Uint64(uint64(time.Now().Unix()*1000))
		} else{
			timeObj := time.Unix(int64(*event.Time), 0)
			if event.Datetime==nil {
				event.Datetime = proto.String(timeObj.Format("2006-01-02 15:04:05"));
			}
			event.Localtime_ms = event.LocalTimeMs
		}
	}

	if header.AppId == nil{
		err:=fmt.Errorf("appid is nil")
		return err
	}
	//2. header的修改,增加 deviceID, 增加aid
	//增加uuID
	header.DeviceId = user.DeviceId
	if header.AppId == nil && header.AppAppId ==nil{
		err:=fmt.Errorf("appid is nil")
		return err
	}
	if header.AppAppId==nil {
		header.AppAppId = header.AppId
	}
	if header.User_unique_id==nil {
		header.User_unique_id = (user.UserUniqueId)
	}
	return nil
}


func (this *McsCollector) McsCollectEvents(caller string, user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) (*http.Response, error) {
	for _,event := range events{
		if event.SessionId == nil {
			return nil,fmt.Errorf("SessionID 为空， 不能写入到离线表")
		}
	}
	ts := uint32(time.Now().Unix())
	message := &pb_event.MarioEvents{
		Caller:     &caller,
		ServerTime: &ts,
		User:       user,
		Header:     header,
		Events:     events,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", this.mscUrl, strings.NewReader(string(data)))
	fmt.Println(strings.NewReader(string(data)))
	if err != nil {
		//logger(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := this.mscHttpClient.Do(req)
	return resp, err
}



/**
 *自动补全 ssid
 * 1. 有UUID ->查询
 * 2. 无UUID -> DeviceId、WebID -> 查询
 */
//func checkSsid(user *pb_event.User, header *pb_event.Header)(error){
//	if user.UserUniqueId==nil && user.WebId ==nil && user.DeviceId== nil{
//		//匿名用户
//		return fmt.Errorf(UUID_IS_EMPTY)
//	}
//	if header.SsId == nil {
//		var webid *string;
//		if user.Ssid == nil{
//			if user.WebId==nil {
//				webid = proto.String(strconv.FormatUint(*user.DeviceId, 10))
//			}
//			tt := &map[string]interface{}{
//				"app_id":header.AppId,
//				"user_unique_id":user.UserUniqueId,
//				"web_id":webid,
//			}
//			client := &http.Client{}
//			data, err := json.Marshal(tt)
//			if err != nil {
//				return  err
//			}
//			req, err := http.NewRequest("POST", "http://"+HttpAddr+":31081/query/ssidinfo", strings.NewReader(string(data)))
//			if err != nil {
//				return err
//			}
//			req.Header.Set("Content-Type", "application/json")
//			resp, err := client.Do(req)
//			if err!=nil {
//				return err
//			}
//			//处理 resp
//			map1:= map[string]interface{}{}
//			body, _ := ioutil.ReadAll(resp.Body)
//
//			json.Unmarshal(body, &map1)
//			if map1["e"]!=0.0 {
//				return fmt.Errorf(SSID_NOT_EXIST)
//			}else{
//				//不补UUID
//				if ssid, ok := map1["ssid"].(string); ok {
//					user.Ssid = &ssid
//				}
//				if did, ok := map1["device_id"].(int); ok {
//					a:= uint64(did)
//					user.DeviceId = &a
//				}
//			}
//		}
//		header.SsId = user.Ssid;
//		if user.WebId == nil{
//			header.Web_id = proto.Uint64(0)
//		}
//	}
//	return nil
//}
//

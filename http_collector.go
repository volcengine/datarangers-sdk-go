package mario_collector

import (
	"code.byted.org/data/mario_collector/pb_event"
	"code.byted.org/dp/mario_common/traceid"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	dialTimeout         = 1 * time.Second
	totalTimeout        = 2 * time.Second
	maxIdleConnsPerHost = 4
	logFileName = flag.String("log", "cServer.log", "Log file name")
	logFile, _ = os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	Info = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

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
func NewAppCollector() (collector *AppCollector) {
	mcsurl := "http://"+HttpAddr+":31081"+AppURL
	collector = &AppCollector{
		*NewMcsCollector(mcsurl,""),
	}
	return
}

//Web小程序上报的接口。
func NewWebMpCollector()(collector *WebCollector){
	mcsurl := "http://"+HttpAddr+":31081"+WebURL
	collector = &WebCollector{
		*NewMcsCollector(mcsurl,""),
	}
	return
}


////App上报的接口
//func NewAppCollector()(collector *McsCollector){
//	mcsurl := "http://"+HttpAddr+":31081"+AppURL
//	return NewMcsCollector(mcsurl, "")
//}

////Web小程序上报的接口。
//func NewWebMpCollector()(collector *McsCollector){
//	mcsurl := "http://"+HttpAddr+":31081"+WebURL
//	return NewMcsCollector(mcsurl, "")
//}


func (this *WebCollector) Collect(user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) (*http.Response, error) {
	err:=checkSsid(user,header)
	if err!=nil {
		return nil,err
	}
	caller :=""
	return this.McsCollectEvents(caller, user, header, events);
}


func (this *AppCollector) Collect( user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) (_ *http.Response, err error) {
	//app_id->aid.
	//paras修改.
	//headers  + deviceId
	//var err error
	//defer logger(err)
	caller:= ""
	//1.修改pras
	//并增加datetime字段
	if err:=MotifyMatchFormatForApp(user, header,events); err!=nil{
		//Error.Println(err)
		//logger(err)
		return nil, err
	}
	//3 补全ssid
	if err:=checkSsid(user, header); err!=nil{
		//logger(err)
		return nil,err
	}

	//4. 构造对应的 launch 和 terminal 。
	//launchtimeObj := time.Unix(time.Now().Unix()-5, 0)
	//launchdatetime := proto.String(launchtimeObj.Format("2006-01-02 15:04:05"));
	//launchlocalTimeMs := proto.Uint64(uint64((time.Now().Unix()-5)*1000))
	//launch:=&pb_event.Launch{
	//	Datetime: launchdatetime,
	//	Localtime_ms: launchlocalTimeMs,
	//	SessionId: events[0].SessionId,
	//	TeaEvent: proto.Uint64(153507),
	//}
	//launchevents:=[]*pb_event.Launch{launch}
	//
	//
	//timeObj := time.Unix(time.Now().Unix()+1000, 0)
	//datetime := proto.String(timeObj.Format("2006-01-02 15:04:05"));
	//localTimeMs := proto.Uint64(uint64((time.Now().Unix()+1000)*1000))
	//terminal:=&pb_event.Terminate{
	//	Datetime: datetime,
	//	Localtime_ms: localTimeMs,
	//	SessionId: events[0].SessionId,
	//	TeaEvent: proto.Uint64(23432),
	//	Duration: proto.Uint64(1213),
	//}
	//terminalevents:=[]*pb_event.Terminate{terminal}

	//5  -> event_v3修改
	//增加launch
	ts := uint32(time.Now().Unix())
	message := &pb_event.MarioEvents{
		Caller:     &caller,
		ServerTime: &ts,
		User:       user,
		Header:     header,
		AppEvents:  events,
		TraceId:    proto.String(traceid.NewTraceId()),
		//Launchs: 	launchevents,
		//Terminates: terminalevents,
	}
	data, err := json.Marshal(message)
	if err != nil {
		//logger(err)
		return nil, err
	}
	//fmt.Println(user.WebId)
	//fmt.Println(strings.NewReader(string(data)))
	req, err := http.NewRequest("POST", this.mscUrl, strings.NewReader(string(data)))
	if err != nil {
		//logger(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X-MCS-AppKey", this.appKey)
	resp, err := this.mscHttpClient.Do(req)
	//logger(err)
	return resp, err
}


func MotifyMatchFormatForApp(user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) error{

	//if user.WebId == nil {
	//	user.WebId = proto.Uint64(123456789)
	//}

	for _, event:=range events{
		var par = make(map[string]interface{})
		data,ok:=event.Params.(*string)
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
			event.Datetime = proto.String(timeObj.Format("2006-01-02 15:04:05"));
			event.LocalTimeMs = proto.Uint64(uint64(time.Now().Unix()*1000))
			event.Localtime_ms = proto.Uint64(uint64(time.Now().Unix()*1000))
		} else{
			timeObj := time.Unix(int64(*event.Time), 0)
			event.Datetime = proto.String(timeObj.Format("2006-01-02 15:04:05"));
			event.Localtime_ms = event.LocalTimeMs
		}
		//添加user ID
		event.UserId = proto.String(strconv.FormatUint(*user.UserId, 10))
	}

	if header.AppId == nil{
		err:=fmt.Errorf("appId 不能为空")
		return err
	}
	//2. header的修改,增加 deviceID, 增加aid
	//增加uuID
	if user.DeviceId==nil{
		err:=fmt.Errorf("user DeviceId 不能为空")
		return err
	}
	header.DeviceId = user.DeviceId
	if header.AppId == nil && header.AppAppId ==nil{
		err:=fmt.Errorf("Header DeviceId 不能为空")
		return err
	}
	if header.AppAppId==nil {
		header.AppAppId = header.AppId
	}
	if header.User_unique_id==nil {
		header.User_unique_id = user.UserUniqueId
	}
	if header.Install_id ==nil {
		header.Install_id = proto.Uint64(98765432123)
	}
	return nil
}


/**
自动补全 ssid
 */
func checkSsid(user *pb_event.User, header *pb_event.Header)(error){
	if user.UserUniqueId==nil {
		return fmt.Errorf("UUID不能为空")
	}
	if header.SsId == nil {
		var webid *string;
		//if user.WebId == nil{
		//	a := "123456789";
		//	webid = &a;
		//}else{
		//	a := strconv.FormatUint(*user.WebId, 10);
		//	webid = &a;
		//}
		if user.Ssid == nil{
			tt := &map[string]interface{}{
				"app_id":header.AppId,
				"user_unique_id":user.UserUniqueId,
				"web_id":webid,
			}
			client := &http.Client{}
			data, err := json.Marshal(tt)
			if err != nil {
				return  err
			}
			req, err := http.NewRequest("POST", "http://"+HttpAddr+":31081/query/ssidinfo", strings.NewReader(string(data)))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err!=nil {
				return err
			}
			//处理 resp
			map1:= map[string]interface{}{}
			body, _ := ioutil.ReadAll(resp.Body)

			json.Unmarshal(body, &map1)
			if map1["e"]!=0.0 {
				err := fmt.Errorf("SSid 不存在")
				return err
			}else{
				if ssid, ok := map1["ssid"].(string); ok {
					user.Ssid = &ssid
				}
				if did, ok := map1["device_id"].(int); ok {
					a:= uint64(did)
					user.DeviceId = &a
				}
			}
		}
		header.SsId = user.Ssid;
		if user.WebId == nil{
			header.Web_id = proto.Uint64(0)
		}
	}
	//fmt.Print("SSID::: ")
	//println(*header.SsId)
	return nil
}




func (this *McsCollector) McsCollectEvents(caller string, user *pb_event.User, header *pb_event.Header, events []*pb_event.Event) (*http.Response, error) {
	ts := uint32(time.Now().Unix())
	message := &pb_event.MarioEvents{
		Caller:     &caller,
		ServerTime: &ts,
		User:       user,
		Header:     header,
		Events:     events,
		TraceId:    proto.String(traceid.NewTraceId()),
	}
	data, err := json.Marshal(message)
	if err != nil {
		logger(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", this.mscUrl, strings.NewReader(string(data)))
	if err != nil {
		//logger(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("X-MCS-AppKey", this.appKey)
	resp, err := this.mscHttpClient.Do(req)
	//logger(err)
	return resp, err
}

//func init() {
//	runtime.GOMAXPROCS(runtime.NumCPU())
//	flag.Parse()
//	//set logfile Stdout
//	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
//	if logErr != nil {
//		fmt.Println("Fail to find", *logFile, "cServer start Failed")
//		os.Exit(1)
//	}
//	log.SetOutput(logFile)
//	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
//	//Info := log.New(logFile,
//	//	"INFO: ",
//	//	log.Ldate|log.Ltime|log.Lshortfile)
//	//Info.Println("ceshiyixia")
//	//write log
//	//log.Printf("1111111111Server abort! Cause:%v \n", "test log file")
//}

func logger(err error){
	print(err)
	if err != nil{
		Error.Println(err)
	}else {
		//Info.Println("没有错误")
	}
}

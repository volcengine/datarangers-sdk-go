package datarangers_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAppCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          true,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/sensors1",
		Log_errlogpath:     "sdklogs1/errlog1",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://10.225.129.5",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})
	//for i := 0; i < 1; i++ {
	//	err := SendEvent(APP, 10000013, "2020_11_22", "old uuid", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
	//	if err != nil {
	//		//fmt.Println(err.Error())
	//	}
	//}
	for i := 0; i < 1; i++ {
		err := SendEventWithDevice(APP, 10000013, "2020_11_22", "old uuid", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1}, IOS, "121321212")
		if err != nil {
			//fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
}


func TestItemCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          false,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/sensors1",
		Log_errlogpath:     "sdklogs1/errlog1",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://10.225.129.5",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})

	for i := 0; i < 1; i++ {
		item1 := &Item{
			ItemName: proto.String("phone"),
			ItemId: proto.String("123"),
		}
		item2 := &Item{
			ItemName: proto.String("book"),
			ItemId: proto.String("124"),
		}
		item3 := &Item{
			ItemName: proto.String("book"),
			ItemId: proto.String("125"),
		}

		itemList := []*Item{}
		itemList = append(itemList, item1 )
		itemList = append(itemList, item2 )
		itemList = append(itemList, item3 )
		err := SendItem(APP, 10000013, "lxy", "buy", map[string]interface{}{"time": 1}, map[string]interface{}{}, itemList)
		if err != nil {
			//fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
}

func TestItemSetCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          false,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/sensors",
		Log_errlogpath:     "sdklogs1/errlog",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://10.225.129.5",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})

	item2 := &Item{
		ItemName: proto.String("book"),
		ItemId: proto.String("124"),
	}
	for i := 0; i < 1; i++ {
		err := ItemSet(APP, 10000000,  map[string]interface{}{"list7":43, "list8":[]string{"bh1"}}, *item2)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
}



func TestProfileCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          true,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/sensors",
		Log_errlogpath:     "sdklogs1/errlog",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://xxxx",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})

	for i := 0; i < 1; i++ {
		err := ProfileAppend(APP, 10000000, "lxy", map[string]interface{}{"list7":43, "list8":[]string{"bh1"}})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
}

func TestSsid(a *testing.T){
	InitByFile("sdkconf.yml")

	threadnum := 128
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Duration(confIns.HttpConfig.SocketTimeOut)*time.Second)
			},
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: 1024,
			MaxIdleConns:        1024,
		},
		Timeout: time.Duration(5) * time.Second,
	}
	var openudid = []int{}
	for i := 0; i< 100001; i++{
		openudid = append(openudid, rand.Int())
	}
	fmt.Println("openudid init over")

	wg := sync.WaitGroup{}
	wg.Add(threadnum)
	for i := 0; i<threadnum; i++{
		var k = i
		go func() {
			for j :=0 ;j< 1000; j++{
				newJ := j+1
				if newJ >= 999 {
					newJ = 999
				}
				head := &H{
					Aid: "2020122014",
					//Openudid: strconv.Itoa(openudid[rand.Intn(100000)]),
					Openudid: strconv.Itoa(k)+"_"+strconv.Itoa(j)+"_",
					User_unique_id: strconv.Itoa(k)+"_"+strconv.Itoa(newJ),
					//Openudid: "2020121410",
					//User_unique_id: "202012167",
				}

				reqestMessage := Req{
					Header: *head,
					//History_register_time: time.Now().Unix(),
				}

				data, _ := json.Marshal(reqestMessage)

				req, _ := http.NewRequest("POST", "http://10.225.129.5:31010/service/2/device_register", strings.NewReader(string(data)))
				req.Header.Add("Content-Type", "application/json")
				//var resp *http.Response
				resp, err := client.Do(req)
				if err!=nil {
					fmt.Println("err", err.Error())
					//for j:=0;j<5; j++{
					//	fmt.Println("err", err.Error())
					//	time.Sleep(1*time.Second)
					//	req, _ := http.NewRequest("POST", "http://10.225.129.79:31010/service/2/device_register", strings.NewReader(string(data)))
					//	req.Header.Add("Content-Type", "application/json")
					//	//var resp *http.Response
					//	if resp, err = client.Do(req); err == nil {
					//		break;
					//	}
					//
					//}
				}
				if err == nil {
					body, _ := ioutil.ReadAll(resp.Body)
					responseMsg := map[string]interface{}{}
					json.Unmarshal(body, &responseMsg)
					responseMsg["user_unique_id"] = head.User_unique_id
					responseMsg["Openuid"] = head.Openudid
					body ,_ = json.Marshal(responseMsg)
					//if responseMsg["ssid"] == "" {
					//	fmt.Errorf("ssid is nil " + head.User_unique_id)
					//	req, _ := http.NewRequest("POST", "http://10.225.129.79:31010/service/2/device_register", strings.NewReader(string(data)))
					//	req.Header.Add("Content-Type", "application/json")
					//	resp, err = client.Do(req)
					//	body, _ = ioutil.ReadAll(resp.Body)
					//}
					errlogger.Println(string(body))
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("-- ok --")
	time.Sleep(1 * time.Second)
}


type Req struct {
	Header 					H   	`json:"header,omitempty"`
	History_register_time 	int64 	`json:"history_register_time,omitempty"`
}


type H struct{
	Aid	   			string   `json:"aid,omitempty"`
	Openudid        string   `json:"openudid,omitempty"`
	User_unique_id  string   `json:"user_unique_id,omitempty"`
}

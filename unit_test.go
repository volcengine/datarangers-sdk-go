package datarangers_sdk

import (
	"fmt"
	"testing"
	"time"
)

func TestAppCollect(a *testing.T) {

	InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          false,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/sensors1",
		Log_errlogpath:     "sdklogs1/errlog1",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://xxxxx",
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

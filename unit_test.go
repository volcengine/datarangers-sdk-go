package datarangers_sdk

import (
	"fmt"
	"testing"
	"time"
)

func TestAppCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		Log_islog:          true,
		Log_loglevel:       "debug", //log level
		Log_iscollect:      true,
		Log_path:           "sdklogs1/sensors",
		Log_errlogpath:     "sdklogs1/errlog",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://10.225.129.3",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})
	for i := 0; i < 1; i++ {
		err := SendEvent(APP, 10000013, "uuidwjx2", "old uuid", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
		if err != nil {
			//fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
}

func TestProfileCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		Log_islog:          true,
		Log_loglevel:       "debug", //log level
		Log_iscollect:      true,
		Log_path:           "sdklogs1/sensors",
		Log_errlogpath:     "sdklogs1/errlog",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://10.225.129.3",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})

	for i := 0; i < 1; i++ {
		err := SendProfile(APP, 10000013, "wjx", SET, map[string]interface{}{"list_test": []string{"a","b"}})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
}

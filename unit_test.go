package datarangers_sdk

import (
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
		Http_addr:          "http://10.225.130.127",
		Http_socketTimeOut: 1,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
	})
	for i := 0; i < 1; i++ {
		err := SendEvent(APP, 10000013, "uuidwjx2", "old uuid", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
		if err != nil {
			//fmt.Println(err.Error())
		}
	}
	//fmt.Println("都塞到队列中了")
	//err := SendEventWithDevice(APP, 10000013, "uuidwjx2", "withdid", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1}, ANDROID, "asdsasdsd112")
	//if err != nil {
	//	//fmt.Println(err.Error())
	//}
	time.Sleep(100 * time.Second)
}

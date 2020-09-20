package datarangers_sdk

import (
	"testing"
	"time"
)

func TestAppCollect(a *testing.T) {

	InitByFile("test/sdkconf1.yml")
	//InitByProperty(&Property{
	//	Log_islog:          true,
	//	Log_loglevel:       "debug", //log level
	//	Log_iscollect:      true,
	//	Log_path:           "sdklogs/sensors",
	//	Log_errlogpath:     "sdklogs/errlog",
	//	Log_maxsize:        30,  //Mb
	//	Log_maxage:         100, //days
	//	Log_maxsbackup:     100, //count
	//	Http_addr:          "http://10.225.130.127",
	//	Http_socketTimeOut: 1,
	//	Asyn_mqlen:         150000,
	//	Asyn_routine:       128,
	//})
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

//
//func TestAppCollectAsyn(t *testing.T) {
//	for i := 0; i < 10000; i++ {
//		err := AppCollectAsyn(10000013, "uuid", "newAppTest", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//	}
//	fmt.Println("加完了")
//	time.Sleep(30 * time.Second)
//}

//func TestWebCollect(t *testing.T) {
//	err := WebCollect(10000013, "uuid", "newWebTest", map[string]interface{}{"param": 2}, map[string]interface{}{"cuns": 1})
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
//func TestWebCollectAsyn(t *testing.T) {
//	for i := 0; i < 1; i++ {
//		err := WebCollectAsyn(10000013, "uuid", "newWebTest", map[string]interface{}{"param": 2}, map[string]interface{}{"cuns": 1})
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//	}
//	fmt.Println("加完了")
//	time.Sleep(30 * time.Second)
//}
//
//func TestInit(t *testing.T) {
//	initConfig()
//	mqlxy = newMq()
//	instance = newExecpool(confIns.Asynconf.Routine)
//	dmg := generate(10000013, "uuid", "newWebTest", map[string]interface{}{"param": 2}, map[string]interface{}{"cuns": 1}, "web")
//	mqlxy.push(dmg)
//	instance.exec()
//
//}

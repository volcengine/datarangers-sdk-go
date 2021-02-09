package datarangers_sdk

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
	"time"
)

func TestAppCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          false,
		Log_loglevel:       "warn", //log level
		Log_path:           "log/rangers.log",
		Log_errlogpath:     "log/errlog1",
		Log_maxsize:        3000,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://xx.xxx.xxx.xx",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         10000000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})
	for i := 0; i < 10000000; i++ {
		//err := SendEventWithDevice(APP, 10000013, "2020_11_22", "old uuid", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1}, IOS, "121321212")
		//if err != nil {
		//	fmt.Println(err.Error())
		//}

		SendEvents(APP, 10000013, "2020_11_22", []string{"event1", "event2"}, []map[string]interface{}{{"event1param": 1}, {"event2param":2}}, map[string]interface{}{"cuns": 1})
		fmt.Println(i)
	}
	time.Sleep(1 * time.Second)
}


func TestItemCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          true,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/sensors1",
		Log_errlogpath:     "sdklogs1/errlog1",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://xx.xxx.xxx.xx",
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
		err := SendItem(10000034, "lxy", "buy", map[string]interface{}{"money": 100}, map[string]interface{}{}, itemList)
		if err != nil {
			//fmt.Println(err.Error())
		}
	}
	time.Sleep(1 * time.Second)
}

func TestItemSetCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          true,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/sensors",
		Log_errlogpath:     "sdklogs1/errlog",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://xx.xxx.xxx.xx",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})

	for i := 0; i < 1; i++ {
		err := ItemSet( 10000034, "book", []map[string]interface{}{{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"},{"id":124, "time11":"死亡之舞2"}})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	//ItemUnset( 10000034, "book", "124", []string{"time"})
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
		Http_addr:          "http://xx.xxx.xxx.xx",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       128,
		Headers:     		map[string]interface{}{},
	})

	//for i := 0; i < 1; i++ {
	//	err := ProfileAppend(10000000, "lxy", map[string]interface{}{"list7":413, "list8":[]string{"b1h1"}})
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//}
	ProfileUnset(10000000, "lxy", []string{"list7"})

	time.Sleep(1 * time.Second)
}

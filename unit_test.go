package datarangers_sdk

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"testing"
	"time"
)

func TestAppCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          true,
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
	for i := 0; i < 1; i++ {

		SendEvents(WEB, 10000001, "", []string{"event1", "event2"}, []map[string]interface{}{{"event1param": 1}, {"event2param":2}}, map[string]interface{}{"cuns": 1}, 6936017802184688384)
	//	SendEvent(WEB, 10000001, "", "event1", map[string]interface{}{}, map[string]interface{}{}, 6936017852940059136)
	//	ProfileSet(WEB, 10000001, "", map[string]interface{}{"list": []string{"a"}},6936017802184688384)
	//	ProfileIncrement(WEB, 10000001, "", map[string]interface{}{"int": 3},6936017802184688384)
	//	ProfileAppend(WEB, 10000001, "", map[string]interface{}{"list": []string{"b"}},6936017802184688384)
	//	ProfileSetOnce(WEB, 10000001, "", map[string]interface{}{"list": []string{"cc"}},6936017802184688384)
	//	ProfileUnset(WEB, 10000001, "", []string{"list"},6936017802184688384)
	}
	time.Sleep(1 * time.Second)
}


func TestItemCollect(a *testing.T) {

	//InitByFile("sdkconf.yml")
	InitByProperty(&Property{
		EventSendEnable:          true,
		Log_loglevel:       "debug", //log level
		Log_path:           "sdklogs1/rangers",
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
		Log_path:           "sdklogs1/rangers",
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
		Log_path:           "sdklogs1/rangers.log",
		Log_errlogpath:     "sdklogs1/errlog",
		Log_maxsize:        30,  //Mb
		Log_maxage:         100, //days
		Log_maxsbackup:     100, //count
		Http_addr:          "http://xx.xxx.xxx.xx",
		Http_socketTimeOut: 10,
		Asyn_mqlen:         150000,
		Asyn_routine:       10,
		Headers:     		map[string]interface{}{},
	})

	for i := 0; i < 10000; i++ {
		err := ProfileAppend(APP,10000004, "lxy_"+strconv.Itoa(i), map[string]interface{}{"list":413})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	//ProfileUnset(10000000, "lxy", []string{"list7"})

	time.Sleep(100 * time.Second)
}


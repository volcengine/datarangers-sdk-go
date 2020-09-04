package main

import (
	sdk "code.byted.org/data/datarangers-sdk-go/event-collect"
	"fmt"
	"time"
)

func main() {
	//err := sdk.AppCollect(10000013, "uuid", "newAppTest", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	aaa()
	//bbb()
}

func aaa() {
	for i := 0; i < 2; i++ {
		err := sdk.AppCollectAsyn(10000013, "uuid", "newAppTest", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	time.Sleep(2 * time.Second)
}

func bbb() {
	err := sdk.AppCollect(10000013, "uuid", "newAppTest", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
	if err != nil {
		fmt.Println(err.Error())
	}
}

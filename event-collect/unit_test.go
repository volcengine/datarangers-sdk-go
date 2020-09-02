package event_collect

import (
	"fmt"
	"testing"
	"time"
)

func TestAppCollect(t *testing.T) {
	err := AppCollect(10000013, "uuid", "newAppTest", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestAppCollectAsyn(t *testing.T) {
	for i := 0; i < 10000; i++ {
		err := AppCollectAsyn(10000013, "uuid", "newAppTest", map[string]interface{}{"param": 1}, map[string]interface{}{"cuns": 1})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("加完了")
	time.Sleep(30 * time.Second)
}

func TestWebCollect(t *testing.T) {
	err := WebCollect(10000013, "uuid", "newWebTest", map[string]interface{}{"param": 2}, map[string]interface{}{"cuns": 1})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestWebCollectAsyn(t *testing.T) {
	for i := 0; i < 1; i++ {
		err := WebCollectAsyn(10000013, "uuid", "newWebTest", map[string]interface{}{"param": 2}, map[string]interface{}{"cuns": 1})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println("加完了")
	time.Sleep(30 * time.Second)
}

func TestInit(t *testing.T) {
	initConfig()
	mqlxy = newMq()
	instance = newExecpool(confIns.Asynconf.Routine)
	dmg := generate(10000013, "uuid", "newWebTest", map[string]interface{}{"param": 2}, map[string]interface{}{"cuns": 1}, "web")
	mqlxy.push(dmg)

	instance.exec()

}

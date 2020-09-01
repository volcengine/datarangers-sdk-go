package test

import (
	"fmt"
	"os"
	"testing"
)
import sycnSdk "code.byted.org/data/datarangers-sdk-go/Synchronize"

func TestAppCollector(t *testing.T) {
	param := map[string]interface{}{}
	param["a"] = "app"
	if _, err := sycnSdk.AppCollect(10000013, "uuid", "testAppCollect20202", param, param); err != nil {
		println(err.Error())
	}
}

func TestWebCollector(t *testing.T) {
	param := map[string]interface{}{}
	param["a"] = "web"
	if _, err := sycnSdk.WebCollect(10000013, "uuid", "testAppCollect20202", param, param); err != nil {
		println(err.Error())
	}
}

func TestWebCollector2(t *testing.T) {

	if err := createFile("aa/bb/cc/dd.log"); err != nil {
		fmt.Println(err.Error())
	}
	if _, err := os.OpenFile("aa/bb/c.logg", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666); err != nil {
		fmt.Println(err)
	}
	//fmt.Println(err.Error())
}

func createFile(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, 0777)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

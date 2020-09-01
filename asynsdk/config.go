package asynsdk

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

var singelConfig *config

// APP
type config struct {
	A asynconf `yaml:"asyn"`
}

type asynconf struct {
	Routine    int    `yaml:"routine"`
	Errlogpath string `yaml:"errlogpath"`
	Mqlen      int    `yaml:"mqlen"`
	Mqwait     int    `yaml:"mqwait"`
}

func initAsynConf() error {
	yamlFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		return err
	}
	singelConfig = &config{}
	err = yaml.Unmarshal(yamlFile, singelConfig)
	if err != nil {
		return err
	}
	fmt.Println(singelConfig.A)
	pathslice := strings.Split(singelConfig.A.Errlogpath, "/")
	pathstr := ""
	for i := 0; i < len(pathslice)-1; i++ {
		pathstr += pathslice[i]
		if i < len(pathslice)-2 {
			pathstr += "/"
		}
	}
	//文件夹
	if err := createFile(pathstr); err != nil {
		return err
	}
	return nil
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

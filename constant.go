package mario_collector

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const (
	AppURL = "/service/2/app_log"
	WebURL = "/v2/event/json"
)

var (
	//net setting
	//HttpAddr = "10.225.130.127"
	HttpAddr string
	//log setting
	ISLOG bool
	LOGPATH string
)

type conf struct {
	Islog 		bool   		`yaml:"islog"` //yaml：yaml格式 enabled：属性的为enabled
	Path    	string 		`yaml:"path"`
	HttpAddr    string 		`yaml:"httpaddr"`
	//AppURL    string `yaml:"path"`
	//WebURL    string `yaml:"path"`
}


func init(){
	yamlFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	c := &conf{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	ISLOG = c.Islog
	HttpAddr = c.HttpAddr
	LOGPATH = c.Path
}
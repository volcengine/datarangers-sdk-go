package event_collect

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

const (
	appURL = "/service/2/app_log"
	webURL = "/v2/event/json"
)

var (
	confIns   *syncConf
	headers   map[string]interface{}
	logger    *log.Logger
	errlogger *log.Logger

	dialTimeout         = 1 * time.Second
	totalTimeout        = 2 * time.Second
	maxIdleConnsPerHost = 4

	isFirst      = true
	webcollector *mcsCollector
	appcollector *mcsCollector

	firstLock = sync.Mutex{}

	mqlxy    *mq
	instance *execpool
)

type syncConf struct {
	EventlogConfig eventlogConfig `yaml:"log"`
	HttpConfig     httpConfig     `yaml:"http"`
	Asynconf       asynconf       `yaml:"asyn"`
}

type eventlogConfig struct {
	Islog      bool   `yaml:"islog"` //yaml：yaml格式 enabled：属性的为enabled
	Iscollect  bool   `yaml:"iscollect"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"maxsize"` // megabytes
	MaxBackups int    `yaml:"maxsbackup"`
	MaxAge     int    `yaml:"maxage"` // days
	ErrPath    string `yaml:"errlogpath"`
}

type httpConfig struct {
	HttpAddr string `yaml:"addr"`
}

type asynconf struct {
	Routine    int    `yaml:"routine"`
	Errlogpath string `yaml:"errlogpath"`
	Mqlen      int    `yaml:"mqlen"`
	Mqwait     int    `yaml:"mqwait"`
}

//同步 异步 都会第一个初始化 此代码。
func initConfig() error {
	yamlFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		return err
	}
	confIns = &syncConf{}
	if err = yaml.Unmarshal(yamlFile, confIns); err != nil {
		return err
	}
	maps := map[string]map[string]interface{}{}
	if err = yaml.Unmarshal(yamlFile, &maps); err != nil {
		return err
	}
	headers = maps["headers"]
	newLog()

	appcollector = newMcsCollector("http://" + confIns.HttpConfig.HttpAddr + ":31081" + appURL)
	webcollector = newMcsCollector("http://" + confIns.HttpConfig.HttpAddr + ":31081" + webURL)
	return nil
}

func newLog() {
	logger = &log.Logger{}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   confIns.EventlogConfig.Path,
		MaxSize:    confIns.EventlogConfig.MaxSize, // megabytes
		MaxBackups: confIns.EventlogConfig.MaxBackups,
		MaxAge:     confIns.EventlogConfig.MaxAge, // days
		LocalTime:  true,
	})
	logger.SetOutput(w)

	if errlogger == nil {
		errlogger = &log.Logger{}
		q := zapcore.AddSync(&lumberjack.Logger{
			Filename:   confIns.EventlogConfig.ErrPath,
			MaxSize:    confIns.EventlogConfig.MaxSize, // megabytes
			MaxBackups: confIns.EventlogConfig.MaxBackups,
			MaxAge:     confIns.EventlogConfig.MaxAge, // days
			LocalTime:  true,
		})
		errlogger.SetOutput(q)
	}
}

func initAsyn() {
	mqlxy = newMq()
	instance = newExecpool(confIns.Asynconf.Routine)
	go instance.exec()
}

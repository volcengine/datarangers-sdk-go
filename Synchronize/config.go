package Synchronize

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var (
	syncConfIns *syncConf
	headers     interface{}
	logger      *log.Logger

	dialTimeout         = 1 * time.Second
	totalTimeout        = 2 * time.Second
	maxIdleConnsPerHost = 4

	isFirst      = true
	webcollector *webCollector
	appcollector *appCollector

	firstLock = sync.Mutex{}
	once      sync.Once
)

type syncConf struct {
	EventlogConfig eventlogConfig `yaml:"log"`
	HttpConfig     httpConfig     `yaml:"http"`
}

type eventlogConfig struct {
	Islog      bool   `yaml:"islog"` //yaml：yaml格式 enabled：属性的为enabled
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"maxsize"` // megabytes
	MaxBackups int    `yaml:"maxsbackup"`
	MaxAge     int    `yaml:"maxage"` // days
	ErrPath    string `yaml:"errlogpath"`
}

type httpConfig struct {
	HttpAddr string `yaml:"addr"`
}

func initConfig() error {
	yamlFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		return err
	}
	syncConfIns = &syncConf{}
	if err = yaml.Unmarshal(yamlFile, syncConfIns); err != nil {
		return err
	}
	maps := map[string]interface{}{}
	if err = yaml.Unmarshal(yamlFile, &maps); err != nil {
		return err
	}
	headers = maps["headers"]
	fmt.Println(headers)
	newLog()
	return nil
}

func newLog() {

	logger = &log.Logger{}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   syncConfIns.EventlogConfig.Path,
		MaxSize:    syncConfIns.EventlogConfig.MaxSize, // megabytes
		MaxBackups: syncConfIns.EventlogConfig.MaxBackups,
		MaxAge:     syncConfIns.EventlogConfig.MaxAge, // days
		LocalTime:  true,
	})
	logger.SetOutput(w)
}

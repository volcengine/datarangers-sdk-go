package datarangers_sdk

/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type AppType string
type deviceType string
type ProfileActionType string
type ItemActionType string

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal

	MP  AppType = "mp"
	APP AppType = "app"
	WEB AppType = "web"
	ALL AppType = "*"

	IOS     deviceType = "IOS"
	ANDROID deviceType = "ANDROID"

	APP_KEY       string = "X-MCS-AppKey"
	AUTHORIZATION string = "Authorization"

	SET        ProfileActionType = "__profile_set"
	UNSET      ProfileActionType = "__profile_unset"
	APPEND     ProfileActionType = "__profile_append"
	SET_ONCE   ProfileActionType = "__profile_set_once"
	INCREAMENT ProfileActionType = "__profile_increment"

	ITEM_SET    ItemActionType = "__item_set"
	ITEM_UNSET  ItemActionType = "__item_unset"
	ITEM_DELETE ItemActionType = "__item_delete"

	EVENT_MESSAGE string = "EVENT"
	USER_MESSAGE  string = "USER"
	ITEM_MESSAGE  string = "ITEM"

	DATE_TIME_LAYOUT string = "2006-01-02 15:04:05"
	SDK_VERSION      string = "datarangers_sdk_go_v2.0.1"

	ENV_PRI         string = "pri"
	ENV_SAAS        string = "saas"
	ENV_SAAS_NATIVE string = "saas_native"

	MODE_HTTP string = "http"
	MODE_FILE string = "file"

	MESSAGE_EVENT string = "event"
	MESSAGE_USER  string = "user"
	MESSAGE_ITEM  string = "item"
)

const (
	DEFAULT_ROUTINE    int = 20
	DEFAULT_QUEUE_SIZE int = 10240

	DEFAULT_FILE_MAX_SIZE   int = 100
	DEFAULT_FILE_MAX_BACKUP int = 0
	DEFAULT_FILE_MAX_AGE    int = 0

	DEFAULT_BATCH_SIZE         int   = 20
	DEFAULT_BATCH_WAIT_TIME_MS int64 = 100

	DEFAULT_SOCKET_TIME_OUT int = 30
)

var (
	//Deprecated:
	defaultConf *Property
	confIns     = &SysConf{}
	headers     map[string]interface{}

	fileWriter    *log.Logger
	errFileWriter *log.Logger

	debugLog    = log.New(os.Stdout, "[DEBUG:]", log.Ldate|log.Ltime)
	infoLog     = log.New(os.Stdout, "[INFO:]", log.Ldate|log.Ltime)
	warnLog     = log.New(os.Stdout, "[WARN:]", log.Ldate|log.Ltime)
	errorLog    = log.New(os.Stdout, "[ERROR:]", log.Ldate|log.Ltime)
	logLevel    int
	logLevelMap = map[string]int{
		"TRACE": 0,
		"DEBUG": 1,
		"INFO":  2,
		"WARN":  3,
		"ERROR": 5,
		"FATAL": 6,
	}

	maxIdleConnsPerHost = 1024

	isInit   = false
	initLock = sync.Mutex{}

	appCollector *mcsCollector
	timezone     int32

	mq       *messageQueue
	instance *execPool

	profileOperationMap = map[string]string{
		string(SET):        "SET",
		string(SET_ONCE):   "SET_ONE",
		string(APPEND):     "APPEND",
		string(INCREAMENT): "INCREASE",
		string(UNSET):      "UNSET",
	}
	itemOperationMap = map[string]string{
		string(ITEM_SET):    "SET",
		string(ITEM_DELETE): "DELETE",
		string(ITEM_UNSET):  "UNSET",
	}

	sasSDomainUrls = map[string]bool{
		"https://mcs.ctobsnssdk.com": true,
		"https://mcs.tobsnssdk.com":  true,
	}

	sasSNativeDomainUrls = map[string]bool{
		"https://gator.volces.com": true,
	}
)

type SysConf struct {
	SdkConfig     SdkConfig        `yaml:"sdk"`
	BatchConfig   BatchConfig      `yaml:"batch"`
	FileConfig    FileConfig       `yaml:"file"`
	HttpConfig    HttpConfig       `yaml:"http"`
	AsynConfig    AsynConfig       `yaml:"asyn"`
	OpenapiConfig OpenapiConfig    `yaml:"openapi"`
	VerifyConfig  VerifyConfig     `yaml:"verify"`
	AppKeys       map[int64]string `yaml:"appKeys"`
}

type BatchConfig struct {
	Enable     bool  `yaml:"enable"`
	Size       int   `yaml:"size"`
	WaitTimeMs int64 `yaml:"waitTimeMs"`
}

type SdkConfig struct {
	Mode     string `yaml:"mode"`
	Env      string `yaml:"env"`
	LogLevel string `yaml:"logLevel"`
}

type FileConfig struct {
	EventSendEnable bool   `yaml:"eventSendEnable"` //yaml：yaml格式 enabled：属性的为enabled，已经过期废弃，不建议使用
	Path            string `yaml:"path"`
	MaxSize         int    `yaml:"maxSize"` // megabytes
	MaxBackup       int    `yaml:"maxBackup"`
	MaxAge          int    `yaml:"maxAge"` // days
	ErrPath         string `yaml:"errPath"`
}

type HttpConfig struct {
	HttpAddr      string                 `yaml:"addr"`
	SocketTimeOut int                    `yaml:"timeout"`
	Headers       map[string]interface{} `yaml:"headers"`
}

type AsynConfig struct {
	Routine   int `yaml:"routine"`
	QueueSize int `yaml:"queueSize"`
}

type OpenapiConfig struct {
	HttpAddr string `yaml:"addr"`
	Ak       string `yaml:"ak"`
	Sk       string `yaml:"sk"`
}

type VerifyConfig struct {
	Url string `yaml:"url"`
}

func initFile() {
	fileWriter = &log.Logger{}
	if confIns.FileConfig.Path == "" {
		confIns.FileConfig.Path = "logs/datarangers.log"
	}
	if confIns.FileConfig.MaxSize < 1 {
		confIns.FileConfig.MaxSize = DEFAULT_FILE_MAX_SIZE
	}
	if confIns.FileConfig.MaxBackup < 0 {
		confIns.FileConfig.MaxBackup = DEFAULT_FILE_MAX_BACKUP
	}
	if confIns.FileConfig.MaxAge < 0 {
		confIns.FileConfig.MaxAge = DEFAULT_FILE_MAX_AGE
	}
	if confIns.FileConfig.ErrPath == "" {
		confIns.FileConfig.ErrPath = "logs/error-datarangers.log"
	}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   confIns.FileConfig.Path,
		MaxSize:    confIns.FileConfig.MaxSize, // megabytes
		MaxBackups: confIns.FileConfig.MaxBackup,
		MaxAge:     confIns.FileConfig.MaxAge, // days
		LocalTime:  true,
	})
	fileWriter.SetOutput(w)

	errFileWriter = &log.Logger{}
	q := zapcore.AddSync(&lumberjack.Logger{
		Filename:   confIns.FileConfig.ErrPath,
		MaxSize:    confIns.FileConfig.MaxSize, // megabytes
		MaxBackups: confIns.FileConfig.MaxBackup,
		MaxAge:     confIns.FileConfig.MaxAge, // days
		LocalTime:  true,
	})
	errFileWriter.SetOutput(q)

}

func debug(s interface{}) {
	if logLevel > LevelDebug {
		return
	}
	debugLog.Println(s)
}

func info(s interface{}) {
	if logLevel > LevelInfo {
		return
	}
	infoLog.Println(s)
}

func warn(s interface{}) {
	if logLevel > LevelWarn {
		return
	}
	warnLog.Println(s)
}

func fatal(s interface{}) {
	if logLevel > LevelError {
		return
	}
	errorLog.Println(s)
}

func initAsyn() {
	if confIns.AsynConfig.Routine < 1 {
		confIns.AsynConfig.Routine = DEFAULT_ROUTINE
	}
	if confIns.AsynConfig.QueueSize < 1 {
		confIns.AsynConfig.QueueSize = DEFAULT_QUEUE_SIZE
	}
	mq = newMq()
	instance = newExecpool(confIns.AsynConfig.Routine)
	go instance.exec()
	debug("init goroutine pool success")
}

func InitBySysConf(conf *SysConf) error {
	if !isInit {
		initLock.Lock()
		defer initLock.Unlock()
		confIns = conf
		maps := confIns.HttpConfig.Headers
		if maps != nil {
			//覆盖参数，而不是替换。
			if headers == nil {
				headers = map[string]interface{}{}
			}
			for k, v := range maps {
				headers[k] = v
			}
		}
		initEnvMode()
		initFile()
		initLogLevel()
		appCollector = newMcsCollector()
		timezone = getTimezone()
		initAsyn()
		initBatch()
		initHook()

		if a, err := json.Marshal(confIns); a != nil {
			debug("user config :" + string(a))
			println("user config :" + string(a))
		} else {
			panic("user config error : " + err.Error())
			return err
		}
		isInit = true
	}
	return nil
}

func InitByFile(path string) error {
	if !isInit {
		initLock.Lock()
		defer initLock.Unlock()
		if !isInit {
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				panic("init config fail->  : " + err.Error())
				return err
			}
			if err = yaml.Unmarshal(yamlFile, confIns); err != nil {
				panic("init config fail->  : " + err.Error())
				return err
			}
			maps := confIns.HttpConfig.Headers
			if maps != nil {
				//覆盖参数，而不是替换。
				if headers == nil {
					headers = map[string]interface{}{}
				}
				for k, v := range maps {
					headers[k] = v
				}
			}
			initEnvMode()
			initFile()
			initLogLevel()
			appCollector = newMcsCollector()
			timezone = getTimezone()
			initAsyn()
			initBatch()
			initHook()

			if a, err := json.Marshal(confIns); a != nil {
				debug("user config :" + string(a))
				println("user config :" + string(a))
			} else {
				panic("user config error : " + err.Error())
				return err
			}
			isInit = true
		}
	}

	return nil
}

func initLogLevel() {
	ll, ok := logLevelMap[strings.ToUpper(confIns.SdkConfig.LogLevel)]
	if !ok {
		logLevel = 0
		warn("logLevel not set, use default 0: trace")
	} else {
		logLevel = ll
	}
}

//InitByProperty
//Deprecated: instead of InitByFile, or InitByConf
func InitByProperty(p *Property) error {
	if !isInit {
		initLock.Lock()
		defer initLock.Unlock()
		if !isInit {
			defaultConf = &Property{
				EventSendEnable:    true,
				Log_path:           "logs/datarangers.log",
				Log_errlogpath:     "logs/error-datarangers.log",
				Log_loglevel:       "trace",
				Log_maxage:         uint32(DEFAULT_FILE_MAX_AGE),
				Log_maxsize:        uint32(DEFAULT_FILE_MAX_SIZE),
				Log_maxsbackup:     uint32(DEFAULT_FILE_MAX_BACKUP),
				Http_addr:          "http://default_http_addr",
				Http_socketTimeOut: DEFAULT_SOCKET_TIME_OUT,
				Asyn_mqlen:         uint32(DEFAULT_QUEUE_SIZE),
				Asyn_routine:       uint32(DEFAULT_ROUTINE),
			}
			defaultConf.EventSendEnable = p.EventSendEnable
			if !p.EventSendEnable && p.Log_path != "" {
				defaultConf.Log_path = p.Log_path
			}
			if p.EventSendEnable && p.Http_addr != "" {
				defaultConf.Http_addr = p.Http_addr
			}
			if p.Log_errlogpath != "" {
				defaultConf.Log_errlogpath = p.Log_errlogpath
			}
			if p.Log_maxsbackup != 0 {
				defaultConf.Log_maxsbackup = p.Log_maxsbackup
			}
			if p.Log_maxsize != 0 {
				defaultConf.Log_maxsize = p.Log_maxsize
			}
			if p.Log_maxage != 0 {
				defaultConf.Log_maxage = p.Log_maxage
			}
			defaultConf.Log_loglevel = p.Log_loglevel
			defaultConf.Headers = p.Headers
			if p.Asyn_routine != 0 {
				defaultConf.Asyn_routine = p.Asyn_routine
			}
			if p.Asyn_mqlen != 0 {
				defaultConf.Asyn_mqlen = p.Asyn_mqlen
			}
			if p.Http_socketTimeOut > 0 {
				defaultConf.Http_socketTimeOut = p.Http_socketTimeOut
			}
			//根据这个初始化 confins。
			confIns.FileConfig.EventSendEnable = defaultConf.EventSendEnable
			confIns.SdkConfig.LogLevel = defaultConf.Log_loglevel
			confIns.FileConfig.ErrPath = defaultConf.Log_errlogpath
			confIns.FileConfig.Path = defaultConf.Log_path
			confIns.FileConfig.MaxAge = int(defaultConf.Log_maxage)
			confIns.FileConfig.MaxBackup = int(defaultConf.Log_maxsbackup)
			confIns.FileConfig.MaxSize = int(defaultConf.Log_maxsize)

			confIns.HttpConfig.HttpAddr = defaultConf.Http_addr
			confIns.HttpConfig.SocketTimeOut = defaultConf.Http_socketTimeOut

			confIns.AsynConfig.Routine = int(defaultConf.Asyn_routine)
			confIns.AsynConfig.QueueSize = int(defaultConf.Asyn_mqlen)

			//初始化其他参数
			if p.Headers != nil {
				//覆盖参数，而不是替换。
				if headers == nil {
					headers = map[string]interface{}{}
				}
				for k, v := range p.Headers {
					headers[k] = v
				}
			}
			initEnvMode()
			initFile()
			initLogLevel()
			appCollector = newMcsCollector()
			timezone = getTimezone()
			initAsyn()
			initBatch()
			initHook()

			if a, err := json.Marshal(confIns); a != nil {
				debug("user config :" + string(a))
				println("user config :" + string(a))
			} else {
				fatal("user config error : " + err.Error())
				return err
			}
			isInit = true
		}
	}
	return nil
}

// Property
//Deprecated: instead of using SysConf
type Property struct {
	//Log_islog          bool
	//Log_iscollect      bool
	EventSendEnable    bool
	Log_path           string
	Log_errlogpath     string
	Log_maxsize        uint32
	Log_maxsbackup     uint32
	Log_maxage         uint32
	Log_loglevel       string
	Http_addr          string //#上报的IP 或 域名, http://10.225.130.127
	Http_socketTimeOut int    //
	Headers            map[string]interface{}
	Asyn_routine       uint32
	Asyn_mqlen         uint32
}

func initHook() {
	info("init hook")
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case sig := <-c:
			{
				length := len(mq.queue)
				info(fmt.Sprintf("try to handle queue before shutdown, length: %d", length))
				for len(mq.queue) != 0 {
					instance.Send()
				}
				println(fmt.Sprintf("sig: %s", sig))
				os.Exit(1)
			}
		}

	}()
}

func initBatch() {
	if !confIns.BatchConfig.Enable {
		return
	}
	if confIns.BatchConfig.Size < 1 {
		confIns.BatchConfig.Size = DEFAULT_BATCH_SIZE
	}
	if confIns.BatchConfig.WaitTimeMs < 1 {
		confIns.BatchConfig.WaitTimeMs = DEFAULT_BATCH_WAIT_TIME_MS
	}
}

func initEnvMode() {
	if confIns.SdkConfig.Env == "" {
		confIns.SdkConfig.Env = ENV_PRI
	}
	confIns.SdkConfig.Env = strings.ToLower(confIns.SdkConfig.Env)
	if sasSDomainUrls[confIns.HttpConfig.HttpAddr] {
		confIns.SdkConfig.Env = ENV_SAAS
	} else if sasSNativeDomainUrls[confIns.HttpConfig.HttpAddr] {
		confIns.SdkConfig.Env = ENV_SAAS_NATIVE
	}

	// mode
	if confIns.SdkConfig.Mode == "" {
		if confIns.FileConfig.EventSendEnable {
			confIns.SdkConfig.Mode = MODE_HTTP
		} else {
			confIns.SdkConfig.Mode = MODE_FILE
		}
	}
	confIns.SdkConfig.Mode = strings.ToLower(confIns.SdkConfig.Mode)
}

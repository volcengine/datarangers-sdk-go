package datarangers_sdk

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Apptype string
type devicetype string
type ProfileActionType string

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal

	MP  Apptype = "mp"
	APP Apptype = "app"
	WEB Apptype = "web"

	IOS     devicetype = "IOS"
	ANDROID devicetype = "ANDROID"

	SET ProfileActionType = "__profile_set"
	UNSET ProfileActionType = "__profile_unset"
	APPEND ProfileActionType = "__profile_append"
	SET_ONCE ProfileActionType = "__profile_set_once"
	INCREAMENT ProfileActionType = "__profile_increment"
)

var (
	defaultconf *Property
	confIns     = &syncConf{}
	headers     map[string]interface{}

	logger    *log.Logger
	errlogger *log.Logger

	debuglog  = log.New(os.Stdout, "[DEBUG:]", log.Ldate|log.Ltime)
	warnlog   = log.New(os.Stdout, "[WARN:]", log.Ldate|log.Ltime)
	errstdlog = log.New(os.Stdout, "[ERROR:]", log.Ldate|log.Ltime)
	loglevel  int

	maxIdleConnsPerHost = 1024

	isFirst               = true
	isFirstConfByProperty = true
	isFirstConfByFile     = true
	isInit                = true

	appcollector *mcsCollector
	timezone     int

	firstLock = sync.Mutex{}
	mqlxy     *mq
	instance  *execpool
)

type syncConf struct {
	EventlogConfig eventlogConfig `yaml:"Log"`
	HttpConfig     httpConfig     `yaml:"Http"`
	Asynconf       asynconf       `yaml:"Asyn"`
}

type eventlogConfig struct {
	Islog      bool   `yaml:"islog"` //yaml：yaml格式 enabled：属性的为enabled
	Iscollect  bool   `yaml:"iscollect"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"maxsize"` // megabytes
	MaxBackups int    `yaml:"maxsbackup"`
	MaxAge     int    `yaml:"maxage"` // days
	ErrPath    string `yaml:"errlogpath"`
	LogLevel   string `yaml:"loglevel"`
}

type httpConfig struct {
	HttpAddr      string `yaml:"addr"`
	SocketTimeOut int    `yaml:"timeout"`
}

type asynconf struct {
	Routine    int    `yaml:"routine"`
	Errlogpath string `yaml:"errlogpath"`
	Mqlen      int    `yaml:"mqlen"`
	Mqwait     int    `yaml:"mqwait"`
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

func debug(s interface{}) {
	if loglevel > LevelDebug {
		return
	}
	debuglog.Println(s)
}

func warn(s interface{}) {
	if loglevel > LevelWarn {
		return
	}
	warnlog.Println(s)
}

func fatal(s interface{}) {
	if loglevel > LevelError {
		return
	}
	errstdlog.Println(s)
}

func initAsyn() {
	mqlxy = newMq()
	instance = newExecpool(confIns.Asynconf.Routine)
	go instance.exec()
}

func InitByFile(path string) error {
	if isFirstConfByFile {
		firstLock.Lock()
		defer firstLock.Unlock()
		if isFirstConfByFile {
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				fatal("初始化文件配置fail->  : " + err.Error() + ", 将使用默认配置")
				return err
			}
			if err = yaml.Unmarshal(yamlFile, confIns); err != nil {
				fatal("初始化文件配置fail->  : " + err.Error() + ", 将使用默认配置")
				return err
			}
			maps := map[string]map[string]interface{}{}
			if err = yaml.Unmarshal(yamlFile, &maps); err != nil {
				fatal("headers配置错误： " + err.Error() + ", 将使用默认headers")
				return err
			}

			if maps["Headers"]!= nil {
				//覆盖参数，而不是替换。
				if headers == nil {
					headers = map[string]interface{}{}
				}
				for k,v:= range maps["Headers"]{
					headers[k] = v;
				}
			}
			newLog()

			appcollector = newMcsCollector(confIns.HttpConfig.HttpAddr + "/sdk/log")
			timezone = getTimezone()
			loglevel, _ = strconv.Atoi(strings.ToUpper(confIns.EventlogConfig.LogLevel))
			switch strings.ToUpper(confIns.EventlogConfig.LogLevel) {
			case "TRACE":
				loglevel = 0
				break
			case "DEBUG":
				loglevel = 1
				break
			case "INFO":
				loglevel = 2
				break
			case "WARN":
				loglevel = 3
				break
			case "ERROR":
				loglevel = 5
				break
			case "FATAL":
				loglevel = 6
				break
			default:
				warn("loglevel 没有设置，默认为 0")
			}
			isFirstConfByFile = false
			if a, err := json.Marshal(confIns); a != nil {
				debug("自定义 File 配置 ：" + string(a))
			} else {
				fatal("自定义 File 配置 error : " + err.Error())
				return err
			}
		}
	}

	return nil
}

func InitByProperty(p *Property) error {

	if isInit || isFirstConfByProperty {
		firstLock.Lock()
		defer firstLock.Unlock()
		if isInit || isFirstConfByProperty {
			defaultconf.Log_islog = p.Log_islog
			defaultconf.Log_iscollect = p.Log_iscollect
			if p.Log_islog && p.Log_path != "" {
				defaultconf.Log_path = p.Log_path
			}
			if p.Log_iscollect && p.Http_addr != "" {
				defaultconf.Http_addr = p.Http_addr
			}
			if p.Log_errlogpath != "" {
				defaultconf.Log_errlogpath = p.Log_errlogpath
			}
			if p.Log_maxsbackup != 0 {
				defaultconf.Log_maxsbackup = p.Log_maxsbackup
			}
			if p.Log_maxsize != 0 {
				defaultconf.Log_maxsize = p.Log_maxsize
			}
			if p.Log_maxage != 0 {
				defaultconf.Log_maxage = p.Log_maxage
			}
			defaultconf.Log_loglevel = p.Log_loglevel
			defaultconf.Headers = p.Headers
			if p.Asyn_routine != 0 {
				defaultconf.Asyn_routine = p.Asyn_routine
			}
			if p.Asyn_mqlen != 0 {
				defaultconf.Asyn_mqlen = p.Asyn_mqlen
			}
			if p.Http_socketTimeOut > 0 {
				defaultconf.Http_socketTimeOut = p.Http_socketTimeOut
			}
			//根据这个初始化 confins。
			confIns.EventlogConfig.Islog = defaultconf.Log_islog
			confIns.EventlogConfig.Iscollect = defaultconf.Log_iscollect
			confIns.EventlogConfig.LogLevel = defaultconf.Log_loglevel
			confIns.EventlogConfig.ErrPath = defaultconf.Log_errlogpath
			confIns.EventlogConfig.Path = defaultconf.Log_path
			confIns.EventlogConfig.MaxAge = int(defaultconf.Log_maxage)
			confIns.EventlogConfig.MaxBackups = int(defaultconf.Log_maxsbackup)
			confIns.EventlogConfig.MaxSize = int(defaultconf.Log_maxsize)

			confIns.HttpConfig.HttpAddr = defaultconf.Http_addr
			confIns.HttpConfig.SocketTimeOut = defaultconf.Http_socketTimeOut

			confIns.Asynconf.Routine = int(defaultconf.Asyn_routine)
			confIns.Asynconf.Mqlen = int(defaultconf.Asyn_mqlen)

			if isInit {
				if a, err := json.Marshal(confIns); a != nil {
					debug("初始化默认配置 :" + string(a))
				} else {
					fatal("初始化默认配置error : " + err.Error())
					return err
				}
			} else {
				if a, err := json.Marshal(confIns); a != nil {
					debug("自定义 Prooerty 配置 :" + string(a))
				} else {
					fatal("自定义 Prooerty 配置 error : " + err.Error())
					return err
				}
				isFirstConfByProperty = false
			}

			//初始化其他参数
			if p.Headers != nil {
				//覆盖参数，而不是替换。
				if headers == nil {
					headers = map[string]interface{}{}
				}
				for k,v:= range p.Headers{
					headers[k] = v;
				}
			}
			newLog()
			appcollector = newMcsCollector(confIns.HttpConfig.HttpAddr + "/sdk/log")
			timezone = getTimezone()
			loglevel, _ = strconv.Atoi(strings.ToUpper(confIns.EventlogConfig.LogLevel))

			switch strings.ToUpper(confIns.EventlogConfig.LogLevel) {
			case "TRACE":
				loglevel = 0
				break
			case "DEBUG":
				loglevel = 1
				break
			case "INFO":
				loglevel = 2
				break
			case "WARN":
				loglevel = 3
				break
			case "ERROR":
				loglevel = 5
				break
			case "FATAL":
				loglevel = 6
				break
			default:
				warn("loglevel 没有设置，默认为 0")
			}
		}
	}
	return nil
}

type Property struct {
	Log_islog          bool
	Log_iscollect      bool
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

//不初始化 按照默认配置走。
func init() {
	defaultconf = &Property{
		Log_islog:          true,
		Log_iscollect:      true,
		Log_path:           "sdklogs/sensors",
		Log_errlogpath:     "sdklogs/errlog",
		Log_loglevel:       "debug",
		Log_maxage:         60,
		Log_maxsize:        10,
		Log_maxsbackup:     10,
		Http_addr:          "http://default_http_addr",
		Http_socketTimeOut: 30,
		Headers:            map[string]interface{}{"host": "snssdk.vpc.com","User-Agent":"GoSDK"},
		Asyn_mqlen:         200000,
		Asyn_routine:       1024,
	}
	InitByProperty(defaultconf)
	isInit = false
}

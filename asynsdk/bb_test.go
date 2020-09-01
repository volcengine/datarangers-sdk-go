package asynsdk

import (
	"fmt"
	loglxy "github.com/donnie4w/go-logger/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"testing"
	"time"
)

func Test_bbb(t *testing.T) {

	//execp := newExecpool(8)
	for i := 0; i < 100; i++ {
		AppCollect(10000013, "uuid", "testApp", map[string]interface{}{"app": 1})
	}

	time.Sleep(100 * time.Second)
	fmt.Println("strat ----")
	//execp.exec()
}

func TestAppCollect(t *testing.T) {
	err := initAsynConf()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(singelConfig.A.Errlogpath)
}

func TestLog(t *testing.T) {
	loglxy.SetRollingFile("logt", "aa.log", 10, 200, loglxy.KB)
	lxy := loglxy.GetLogger()
	for i := 0; i < 100000; i++ {
		lxy.Info("go: inconsistent vendoring in /Users/bytedance/golang/src/code.byted.org/data/datarangers-sdk-go:\n\tgithub.com/donnie4w/go-logger@v0.0.0-20170827050443-4740c51383f4: is explicitly required in go.mod, but not marked as explicit in vendor/modules.txt")
	}
}

func TestRota(t *testing.T) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "var/log/myapp/foo.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	lxy := zap.New(core)
	for i := 0; i < 100000; i++ {
		lxy.Info("go: inconsistent vendoring in /Users/bytedance/golang/src/code.byted.org/data/datarangers-sdk-go:\n\tgithub.com/donnie4w/go-logger@v0.0.0-20170827050443-4740c51383f4: is explicitly required in go.mod, but not marked as explicit in vendor/modules.txt")
	}
}

func TestTtt(t *testing.T) {
	logger := &log.Logger{}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "aaa/log",
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		LocalTime:  true,
	})
	//writer, _ := rotatelogs.New(
	//	lOGPATH+".%Y%m%d%H%M",
	//	//rotatelogs.WithLinkName(lOGPATH),
	//	rotatelogs.WithRotationTime(time.Duration(splitTime)*time.Minute),
	//)
	logger.SetOutput(w)
	for i := 0; i < 100000; i++ {
		logger.Println("go: inconsistent vendoring in /Users/bytedance/golang/src/code.byted.org/data/datarangers-sdk-go:\n\tgithub.com/donnie4w/go-logger@v0.0.0-20170827050443-4740c51383f4: is explicitly required in go.mod, but not marked as explicit in vendor/modules.txt")
	}
}

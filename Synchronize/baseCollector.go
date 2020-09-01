package Synchronize

import (
	"net"
	"net/http"
)

type mcsCollector struct {
	mscUrl        string
	appKey        string
	mscHttpClient *http.Client
}

func newMcsCollector(mcsUrl string, appKey string) (collector *mcsCollector) {
	collector = &mcsCollector{
		mscUrl: mcsUrl,
		appKey: appKey,
		mscHttpClient: &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, dialTimeout)
				},
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
			},
			Timeout: totalTimeout,
		},
	}
	return
}

//
//func initL() error {
//	if err := initConf(); err != nil {
//		return err
//	}
//	if err := initLog(); err != nil {
//		return err
//	}
//	return nil
//}
//
//func initConf() error {
//	yamlFile, err := ioutil.ReadFile("conf.yml")
//	if err != nil {
//		return err
//	}
//	c := &conf{}
//	err = yaml.Unmarshal(yamlFile, c)
//	if err != nil {
//		return err
//	}
//
//	iSLOG = c.A.Islog
//	splitTime = c.A.Splittime
//	//路径这要看一下。
//	lOGPATH = c.A.Path
//	pathslice := strings.Split(lOGPATH, "/")
//	pathstr := ""
//	for i := 0; i < len(pathslice)-1; i++ {
//		pathstr += pathslice[i]
//		if i < len(pathslice)-2 {
//			pathstr += "/"
//		}
//	}
//	//文件夹
//	if err := createFile(pathstr); err != nil {
//		return err
//	}
//	//创建log
//	//if _, err := os.OpenFile(lOGPATH, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666); err != nil {
//	//	return err
//	//}
//	//HttpAddr 有问题。
//	httpAddr = c.B.HttpAddr
//	url := "http://" + httpAddr + "/healthz"
//	response, err := http.Head(url)
//	if err != nil {
//		return err
//	}
//	if response.StatusCode != http.StatusOK {
//		bytea, _ := ioutil.ReadAll(response.Body)
//		return fmt.Errorf(string(bytea))
//	}
//	return nil
//}
//
//func initLog() error {
//	if true {
//		logger = &log.Logger{}
//		writer, _ := rotatelogs.New(
//			lOGPATH+".%Y%m%d%H%M",
//			//rotatelogs.WithLinkName(lOGPATH),
//			rotatelogs.WithRotationTime(time.Duration(splitTime)*time.Minute),
//		)
//		logger.SetOutput(writer)
//	}
//	return nil
//}
//
//func createFile(filePath string) error {
//	if !isExist(filePath) {
//		err := os.MkdirAll(filePath, 0777)
//		return err
//	}
//	return nil
//}
//
//// 判断所给路径文件/文件夹是否存在(返回true是存在)
//func isExist(path string) bool {
//	_, err := os.Stat(path) //os.Stat获取文件信息
//	if err != nil {
//		if os.IsExist(err) {
//			return true
//		}
//		return false
//	}
//	return true
//}

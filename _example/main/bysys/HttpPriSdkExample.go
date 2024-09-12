/*
 * Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package main

import (
	sdk "github.com/volcengine/datarangers-sdk-go"
	m "github.com/volcengine/datarangers-sdk-go/_example"
	"os"
	"time"
)

func main() {
	// 初始化
	sdk.InitBySysConf(&sdk.SysConf{
		// 设置模式
		SdkConfig: sdk.SdkConfig{
			Mode: sdk.MODE_HTTP,
		},
		// 设置sdk domain，这里注意替换成真实的参数
		HttpConfig: sdk.HttpConfig{
			HttpAddr: os.Getenv("SDK_DOMAIN"),
		},
	})

	appId := 10000000
	sdkExample := m.SDKExample{
		AppId: int64(appId),
		Uuid:  "test_go_sdk_user1",
	}

	// 上报事件
	sdkExample.SendEventInfo()
	sdkExample.SendEventInfos()
	sdkExample.SendEventInfoWithHeader()

	// 避免程序立刻退出
	time.Sleep(60 * time.Second)
}

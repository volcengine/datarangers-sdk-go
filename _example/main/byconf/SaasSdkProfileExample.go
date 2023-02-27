/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package main

import (
	sdk "github.com/volcengine/datarangers-sdk-go"
	m "github.com/volcengine/datarangers-sdk-go/_example"
	"os"
	"strconv"
	"time"
)

func main() {
	// 初始化
	sdk.InitByFile("_example/conf/sdkconf_saas.yml")

	appId, _ := strconv.Atoi(os.Getenv("SDK_APP_1"))
	sdkExample := m.SDKExample{
		AppId: int64(appId),
		Uuid:  "test_go_sdk_user2",
	}

	// 上报事件
	sdkExample.SendEventInfo()
	sdkExample.SendEventInfos()
	sdkExample.SendEventInfoWithHeader()

	//上报用户属性，需要保证先在系统新增用户属性
	sdkExample.SetProfile()

	//item 需要先在管理页面进行创建
	sdkExample.SetItem()

	// 创建完成之后，再在事件中上报关联的item
	sdkExample.SendEventWithItem()

	// 避免程序立刻退出
	time.Sleep(60 * time.Second)
}

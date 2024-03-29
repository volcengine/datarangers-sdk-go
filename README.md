#简介
此版本为数据流优化后服务端sdk, profile事件通过事件流统一接口上报。
主要接口及使用方法见_example目录下文件。

#改动说明
##2023.01.28
统一配置方式
支持saas、saas云原生、私有化多种环境上报
支持批量上报
重构代码，初始化会有不兼容的地方
添加_example

##2020.9.18
增加了多种初始化配置方式。
增加了对statuscode的验证。
增加了vendor字段，提供设备信息接口

##2020.9.27
增加了body体的验证。
修复errlog-path不正确问题。

##2020.10.13
增加了sendprofile接口，增加了5种profile操作类型。
增加了User-Agent，防止Applog反爬虫。

##2020.10.26
修复了headers设置覆盖失败
问题。

##2020.11.02
修复了InitByFile后UA覆盖失败问题。 
增加了5种Profile方法

##2020.12.20
将内部数据格式设置为public。
增加了Item上报事件接口。
增加了设置ItemSet ItemUnset接口。

##2020.12.23
1. 增加批量上报Event的接口 sendEvents
2. 修改ProfileUnset的接口定义。
3. 修改ItemSet与ItemUnset的定义。

##2021.02.24
1. 修改了did为1的问题

##2021.04.06
1. 增加了 绑定 did 设置接口
2. profile 接口增加 apptype
3. 增加custom中的版本号


## License
Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. 
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

#简介
此版本为数据流优化后服务端sdk, profile事件通过事件流统一接口上报。
主要接口及使用方法见unit_test.go文件。

#改动说明
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
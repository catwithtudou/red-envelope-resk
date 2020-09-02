package main

import (
	"github.com/catwithtudou/red-envelope-infra"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	_ "red-envelope"
)

func main() {
	//获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("config.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	app := infra.New(conf)
	app.Start()
}

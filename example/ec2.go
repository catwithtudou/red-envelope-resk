package main

import (
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/tietang/props/ini"
)

/**
 *@Author tudou
 *@Date 2020/8/30
 **/


func main() {

	conf := ini.NewIniFileConfigSource("ec.ini")
	client := eureka.NewClient(conf)
	client.Start()
	c := make(chan int, 1)
	<-c
}

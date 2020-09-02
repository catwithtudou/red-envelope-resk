package main


import (
	"github.com/tietang/go-eureka-client/eureka"
	"time"
)

/**
 *@Author tudou
 *@Date 2020/8/30
 **/


func main(){
	cfg := eureka.Config{
		DialTimeout: time.Second * 10,
	}
	client := eureka.NewClientByConfig([]string{
		"http://127.0.0.1:8761/eureka",
	}, cfg)
	appName := "Go-Example"
	instance := eureka.NewInstanceInfo(
		"test.com", appName,
		"127.0.0.2",
		8080, 30,
		false)
	client.RegisterInstance(appName, instance)
	client.Start()
	c := make(chan int, 1)
	<-c



}
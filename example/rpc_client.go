package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"red-envelope/services"
)

/**
 *@Author tudou
 *@Date 2020/7/28
 **/


func main(){
	c,err:=rpc.Dial("tcp",":8082")
	if err!=nil{
		logrus.Panic(err)
	}
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(5),
		UserId:       "1fFn6yMPqwe21WwYvgLzMkfhOso",
		Username:     "测试资金账户",
		EnvelopeType: services.GeneralEnvelopeType,
		Quantity:     2,
		Blessing:     "",
	}
	out := &services.RedEnvelopeActivity{}
	c.Call("EnvelopeRpc.SendOut",in,&out)
	logrus.Infof("%+v",out)
}
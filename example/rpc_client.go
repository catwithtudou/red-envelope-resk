package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"github.com/catwithtudou/red-envelope-resk/services"
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
	sendOut(c)
	receive(c)

}

func receive(c *rpc.Client){
	in := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   "",
		RecvUserId:   "",
		RecvUsername: "",
		AccountNo:    "",
	}
	out := &services.RedEnvelopeItemDTO{}
	err := c.Call("Envelope.Receive", in, &out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("%+v", out)
}

func sendOut(c *rpc.Client){
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(100),
		UserId:       "1fM21rA58Nlm954VXFDZ1oZQsLI",
		Username:     "测试账户10",
		EnvelopeType: services.GeneralEnvelopeType,
		Quantity:     2,
		Blessing:     "",
	}
	out := &services.RedEnvelopeActivity{}
	err:=c.Call("EnvelopeRpc.SendOut",in,&out)
	if err!=nil{
		logrus.Panic(err)
	}
	logrus.Infof("%+v",out)
}
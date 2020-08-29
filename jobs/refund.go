package jobs

import (
	"github.com/sirupsen/logrus"
	"red-envelope/infra"
	"time"
)

/**
 *@Author tudou
 *@Date 2020/8/29
 **/



type RefundExpiredJobStarter struct{
	infra.BaseStarter
	ticker *time.Ticker
}


func (r *RefundExpiredJobStarter)Init(ctx infra.StarterContext){
	d:=ctx.Props().GetDurationDefault("jobs.refund.interval",time.Minute)
	r.ticker = time.NewTicker(d)
}


func (r *RefundExpiredJobStarter)Start(ctx infra.StarterContext){
	go func() {

		for{
			c:=<-r.ticker.C
			logrus.Debug("过期红包退款开始...",c)
			//红包过期退款的业务逻辑代码
		}


	}()
}


func (r *RefundExpiredJobStarter)Stop(ctx infra.StarterContext){
	r.ticker.Stop()
}
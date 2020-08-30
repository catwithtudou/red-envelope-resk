package jobs

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"red-envelope/core/envelopes"
	"red-envelope/infra"
	"time"
	"github.com/go-redsync/redsync"
)

/**
 *@Author tudou
 *@Date 2020/8/29
 **/



type RefundExpiredJobStarter struct{
	infra.BaseStarter
	ticker *time.Ticker
	mutex *redsync.Mutex
}


func (r *RefundExpiredJobStarter)Init(ctx infra.StarterContext){
	d:=ctx.Props().GetDurationDefault("jobs.refund.interval",time.Minute)
	r.ticker = time.NewTicker(d)

	//构建分布式锁连接池
	maxIdle:=ctx.Props().GetIntDefault("redis.maxIdle",2)
	maxActive:= ctx.Props().GetIntDefault("redis.maxActive",5)
	timeout:=ctx.Props().GetDurationDefault("redis.timeout",20*time.Second)
	addr:=ctx.Props().GetDefault("redis.addr","127.0.0.1:6379")
	pools:=make([]redsync.Pool,0)
	pool:=&redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",addr)
		},
		TestOnBorrow:    nil,
		MaxIdle:         maxIdle,
		MaxActive:       maxActive,
		IdleTimeout:     timeout,
		Wait:            false,
		MaxConnLifetime: 0,
	}
	pools=append(pools,pool)
	rsync:=redsync.New(pools)
	ip,err:=utils.GetExternalIP()
	if err!=nil{
		ip="127.0.0.1"
	}
	r.mutex=rsync.NewMutex("lock:RefundExpired",
		redsync.SetExpiry(50 * time.Second),
		redsync.SetRetryDelay(3),
		redsync.SetGenValueFunc(func() (string, error) {
			now:=time.Now()
			logrus.Infof("节点%s正在执行过期红包的退款任务",ip)
			return fmt.Sprintf("%d:%s",now.Unix(),ip),nil
		}))


}


func (r *RefundExpiredJobStarter)Start(ctx infra.StarterContext){
	go func() {

		for{
			c:=<-r.ticker.C
			err:=r.mutex.Lock()
			if err!=nil{
				//如果拿到分布式锁
				//红包过期退款的业务逻辑代码
				logrus.Debug("过期红包退款开始...",c)
				refundDomain:=envelopes.ExpiredEnvelopeDomain{}
				_, err = refundDomain.Expired()
				if err!=nil{
					logrus.Error(err)
					return
				}
			}else{
				logrus.Info("已经有节点在运行该任务")
			}
			_, _ = r.mutex.Unlock()
		}


	}()
}


func (r *RefundExpiredJobStarter)Stop(ctx infra.StarterContext){
	r.ticker.Stop()
}
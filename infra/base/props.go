package base

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"red-envelope/infra"
)

var props kvs.ConfigSource

//Props 配置文件获取客户端
func Props() kvs.ConfigSource {
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	log.Info("初始化配置.")
}

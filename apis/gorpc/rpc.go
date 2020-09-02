package gorpc

import (
	"github.com/catwithtudou/red-envelope-infra"
	"github.com/catwithtudou/red-envelope-infra/base"
)

/**
 *@Author tudou
 *@Date 2020/7/28
 **/

type GoRpcApiStarter struct {
	infra.BaseStarter
}

func (g *GoRpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}

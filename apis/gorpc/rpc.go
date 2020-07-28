package gorpc

import (
	"red-envelope/infra"
	"red-envelope/infra/base"
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
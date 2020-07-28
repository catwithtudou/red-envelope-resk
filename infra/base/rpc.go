package base

import (
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"red-envelope/infra"
	"reflect"
)

/**
 *@Author tudou
 *@Date 2020/7/28
 **/

var rpcServer *rpc.Server


func RpcServer()*rpc.Server{
	Check(rpcServer)
	return rpcServer
}

func RpcRegister(ri interface{}){
	typ:=reflect.TypeOf(ri)
	logrus.Infof("goRPC Register : %s",typ.String())
	RpcServer().Register(ri)

}


type GoRpcStarter struct{
	infra.BaseStarter
	server *rpc.Server
}

func (s *GoRpcStarter)Init(ctx infra.StarterContext){
	s.server=rpc.NewServer()
	rpcServer = s.server
}

func (s *GoRpcStarter)Start(ctx infra.StarterContext){
	port:=ctx.Props().GetDefault("app.rpc.port","8082")
	//监听网络端口
	listener,err:=net.Listen("tcp",":"+port)
	if err!=nil{
		logrus.Panic(err)
	}
	logrus.Info("tcp port listened for rpc:",port)
	//处理网络连接和请求
	go s.server.Accept(listener)

}
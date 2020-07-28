package web

import (
	"github.com/kataras/iris"
	"red-envelope/infra"
	"red-envelope/infra/base"
	"red-envelope/services"
)

func init(){
	infra.RegisterApi(new(EnvelopeApi))
}


type EnvelopeApi struct{
	service services.RedEnvelopeService
}

func(e *EnvelopeApi)Init(){
	e.service=services.GetRedEnvelopeService()
	gorupRouter:=base.Iris().Party("/v1/envelope")
	gorupRouter.Post("/sendout",e.sendOutHandler)
}




func (e *EnvelopeApi)sendOutHandler(ctx iris.Context){
	dto:=services.RedEnvelopeSendingDTO{}
	err:=ctx.ReadJSON(&dto)
	r:=base.Res{
		Code:base.ResCodeOk,
	}
	if err!=nil{
		r.Code=base.ResCodeRequestParamsError
		r.Message=err.Error()
		ctx.JSON(&r)
		return
	}
	activity,err:= e.service.SendOut(dto)
	if err!=nil{
		r.Code=base.ResCodeInnerServerError
		r.Message=err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data=activity
	ctx.JSON(r)
	return

}
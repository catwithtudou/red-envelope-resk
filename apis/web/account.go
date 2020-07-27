package web

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"red-envelope/infra"
	"red-envelope/infra/base"
	"red-envelope/services"
)

//统一前缀

func init(){
	infra.RegisterApi(new(AccountApi))
}

type AccountApi struct{

}

func ( a *AccountApi) Init(){
	groupRouter:=base.Iris().Party("/v1/account")
	create(groupRouter)
}



//账户创建：/v1/account/create
//POST body json
/*
{
	"UserId": "w123456",
	"Username": "测试用户1",
	"AccountName": "测试账户1",
	"AccountType": 0,
	"CurrencyCode": "CNY",
	"Amount": "100.11"
}

{
    "code": 1000,
    "message": "",
    "data": {
        "AccountNo": "1K1hrG0sQw7lDuF6KOQbMBe2o3n",
        "AccountName": "测试账户1",
        "AccountType": 0,
        "CurrencyCode": "CNY",
        "UserId": "w123456",
        "Username": "测试用户1",
        "Balance": "100.11",
        "Status": 1,
        "CreatedAt": "2019-04-18T13:26:51.895+08:00",
        "UpdatedAt": "2019-04-18T13:26:51.895+08:00"
    }
}
*/
func create(groupRouter iris.Party){
	groupRouter.Post("/create", func(context iris.Context) {
		account:=services.AccountCreatedDTO{}
		err:=context.ReadJSON(&account)
		r:=base.Res{
			Code:base.ResCodeOk,
		}
		if err!=nil{
			r.Code=base.ResCodeRequestParamsError
			r.Message=err.Error()
			context.JSON(&r)
			logrus.Error(err)
			return
		}

		service:=services.GetAccountService()
		dto,err:=service.CreateAccount(account)
		if err!=nil{
			r.Code = base.ResCodeInnerServerError
			r.Message = err.Error()
			logrus.Error(err)
		}
		r.Data = dto
		context.JSON(&r)
		return
	})
}

//转账：/v1/account/transfer

//查询红包账户：/v1/account/envelope/get

//查询账户信息：/v11/account/get
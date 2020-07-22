package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"red-envelope/services"
	"testing"
)

/**
 *@Author tudou
 *@Date 2020/7/22
 **/


func TestAccountDomain_Create(t *testing.T) {
	dto:=services.AccountDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Balance:      decimal.NewFromFloat(0),
		Status:       1,
	}
	domain := new(accountDomain)
	Convey("账户创建",t, func() {
		rdto,err:=domain.Create(dto)
		So(err,ShouldBeNil)
		So(rdto,ShouldNotBeNil)
		So(rdto.Balance.String(),ShouldEqual,dto.Balance.String())
		So(rdto.UserId,ShouldEqual,dto.UserId)
		So(rdto.Username,ShouldEqual,dto.Username)
		So(rdto.Status,ShouldEqual,dto.Status)
	})
}


func TestAccountDomain_Transfer(t *testing.T) {
	//2个账户，交易主体账户要有
	adto1:=&services.AccountDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户1",
		Balance:      decimal.NewFromFloat(100),
		Status:       1,
	}
	adto2:=&services.AccountDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户2",
		Balance:      decimal.NewFromFloat(100),
		Status:       1,
	}
	domain := accountDomain{}
	Convey("转账测试",t, func() {
		dot1,err:=domain.Create(*adto1)
		So(err,ShouldBeNil)
		So(dot1,ShouldNotBeNil)
		So(dot1.Balance.String(),ShouldEqual,dot1.Balance.String())
		So(dot1.UserId,ShouldEqual,dot1.UserId)
		So(dot1.Username,ShouldEqual,dot1.Username)
		So(dot1.Status,ShouldEqual,dot1.Status)
		adto1=dot1

		dot2,err:=domain.Create(*adto2)
		So(err,ShouldBeNil)
		So(dot2,ShouldNotBeNil)
		So(dot2.Balance.String(),ShouldEqual,dot2.Balance.String())
		So(dot2.UserId,ShouldEqual,dot2.UserId)
		So(dot2.Username,ShouldEqual,dot2.Username)
		So(dot2.Status,ShouldEqual,dot2.Status)
		adto2=dot2

		//转账操作验证
		//1. 余额充足，金额转入其他账户
		Convey("余额充足，金额转入其他账户", func() {
			amount := decimal.NewFromFloat(1)
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId: adto1.UserId,
				Username: adto1.Username,
			}
			target:=services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId: adto2.UserId,
				Username: adto2.Username,
			}
			dto:=services.AccountTransferDTO{
				TradeBody: body,
				TradeTarget: target,
				TradeNo: ksuid.New().Next().String(),
				Amount: amount,
				ChangeType: services.ChangeType(-1),
				ChangeFlag: services.FlagTransferOut,
				Decs: "转账",
			}
			status,err:=domain.Transfer(dto)
			So(err,ShouldBeNil)
			So(status,ShouldEqual,services.TransferedStatusSuccess)

			//实际余额更新过后的预期值

			a2:=domain.GetAccount(adto1.AccountNo)
			So(a2,ShouldNotBeNil)
			So(a2.Balance.String(),ShouldEqual,adto1.Balance.Sub(amount).String())
		})

		//2. 余额不足，金额转出
		Convey("余额不足，金额转出", func() {
			amount := decimal.NewFromFloat(200)
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId: adto1.UserId,
				Username: adto1.Username,
			}
			target:=services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId: adto2.UserId,
				Username: adto2.Username,
			}
			dto:=services.AccountTransferDTO{
				TradeBody: body,
				TradeTarget: target,
				TradeNo: ksuid.New().Next().String(),
				Amount: amount,
				ChangeType: services.ChangeType(-1),
				ChangeFlag: services.FlagTransferOut,
				Decs: "转账",
			}
			status,err:=domain.Transfer(dto)
			So(err,ShouldNotBeNil)
			So(status,ShouldEqual,services.TransferedStatusSufficientFunds)
			//实际余额更新过后的预期值

			a2:=domain.GetAccount(adto1.AccountNo)
			So(a2,ShouldNotBeNil)
			So(a2.Balance.String(),ShouldEqual,adto1.Balance.String())

		})

		//3. 充值
		Convey("充值", func() {
			amount := decimal.NewFromFloat(1.1)
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId: adto1.UserId,
				Username: adto1.Username,
			}
			target:=services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId: adto2.UserId,
				Username: adto2.Username,
			}
			dto:=services.AccountTransferDTO{
				TradeBody: body,
				TradeTarget: target,
				TradeNo: ksuid.New().Next().String(),
				Amount: amount,
				ChangeType: services.AccountStoreValue,
				ChangeFlag: services.FlagTransferIn,
				Decs: "充值",
			}
			status,err:=domain.Transfer(dto)
			So(err,ShouldBeNil)
			So(status,ShouldEqual,services.TransferedStatusSuccess)
			//实际余额更新过后的预期值


			a2:=domain.GetAccount(adto1.AccountNo)
			So(a2,ShouldNotBeNil)
			So(a2.Balance.String(),ShouldEqual,adto1.Balance.Add(amount).String())

		})

	})
}
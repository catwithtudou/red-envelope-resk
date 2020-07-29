package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"red-envelope/services"
	"strconv"
	"testing"
)

/**
 *@Author tudou
 *@Date 2020/7/28
 **/

func TestRedEnvelopeService_SendOut(t *testing.T) {
	//发红包人的红包资金账户
	ac:=services.GetAccountService()
	account := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "1000",
		AccountName:  "测试账户",
		AccountType:  int(services.EnvelopeAccountType),
		CurrencyCode: "CNY",
	}
	re:=services.GetRedEnvelopeService()
	Convey("准备资金账户", t, func() {
		//准备资金账户
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO, ShouldNotBeNil)
	})
	Convey("发送红包",t, func() {
		goods:=services.RedEnvelopeSendingDTO{
			UserId: account.UserId,
			Username: account.Username,
			EnvelopeType: services.GeneralEnvelopeType,
			Amount: decimal.NewFromFloat(8.88),
			Quantity: 10,
			Blessing: services.DefaultBlessing,
		}

		Convey("发普通红包", func() {
			at,err:=re.SendOut(goods)
			So(err,ShouldBeNil)
			So(at,ShouldNotBeNil)
			So(at.Link,ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO,ShouldNotBeNil)

			//验证每一个属性
			dto:=at.RedEnvelopeGoodsDTO
			So(dto.Username,ShouldEqual,goods.Username)
			So(dto.UserId,ShouldEqual,goods.UserId)
			So(dto.Quantity,ShouldEqual,goods.Quantity)
			q := decimal.NewFromFloat(float64(dto.Quantity))
			So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(q).String())

		})

		goods.EnvelopeType=services.LuckyEnvelopeType
		goods.Amount=decimal.NewFromFloat(88.8)

		Convey("发碰运气红包", func() {
			at,err:=re.SendOut(goods)
			So(err,ShouldBeNil)
			So(at,ShouldNotBeNil)
			So(at.Link,ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO,ShouldNotBeNil)

			//验证每一个属性
			dto:=at.RedEnvelopeGoodsDTO
			So(dto.Username,ShouldEqual,goods.Username)
			So(dto.UserId,ShouldEqual,goods.UserId)
			So(dto.Quantity,ShouldEqual,goods.Quantity)
			So(dto.Amount.String(), ShouldEqual, goods.Amount.String())

		})

	})
}


func TestRedEnvelopeService_Receive(t *testing.T) {
	//1.准备几个红包资金账户，用户发红包和收红包
	accountService:=services.GetAccountService()

	Convey("收红包测试用例",t, func() {
		accounts:=make([]*services.AccountDTO,0)
		size:=10
		for i:=0;i<size;i++{
			account := services.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户"+strconv.Itoa(i+1),
				Amount:       "2000",
				AccountName:  "测试账户"+strconv.Itoa(i+1),
				AccountType:  int(services.EnvelopeAccountType),
				CurrencyCode: "CNY",
			}
			//账户创建
			accountDto,err:=accountService.CreateAccount(account)
			So(err,ShouldBeNil)
			So(accountDto,ShouldNotBeNil)

			accounts=append(accounts,accountDto)
		}

		acDto:=accounts[0]
		So(len(accounts), ShouldEqual, size)
		//2. 使用其中一个用户发送一个红包
		re := services.GetRedEnvelopeService()
		//发送普通红包
		goods := services.RedEnvelopeSendingDTO{
			UserId:       acDto.UserId,
			Username:     acDto.Username,
			EnvelopeType: services.GeneralEnvelopeType,
			Amount:       decimal.NewFromFloat(1.88),
			Quantity:     size,
			Blessing:     services.DefaultBlessing,
		}

		at, err := re.SendOut(goods)
		So(err, ShouldBeNil)
		So(at, ShouldNotBeNil)
		So(at.Link, ShouldNotBeEmpty)
		So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
		//验证每一个属性
		dto := at.RedEnvelopeGoodsDTO
		So(dto.Username, ShouldEqual, goods.Username)
		So(dto.UserId, ShouldEqual, goods.UserId)
		So(dto.Quantity, ShouldEqual, goods.Quantity)
		q := decimal.NewFromFloat(float64(dto.Quantity))
		So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(q).String())


		//3.使用发送红包数量的人收红包
		remainAmount := at.Amount
		Convey("收普通红包", func() {
			for _,account:=range accounts{
				rcv:=services.RedEnvelopeReceiveDTO{
					EnvelopeNo: at.EnvelopeNo,
					RecvUserId: account.UserId,
					RecvUsername: account.Username,
					AccountNo: account.AccountNo,
				}
				item,err:=re.Receive(rcv)
				So(err,ShouldBeNil)
				So(item,ShouldNotBeNil)
				So(item.Amount,ShouldEqual,at.AmountOne)
				remainAmount=remainAmount.Sub(at.AmountOne)
				So(item.RemainAmount.String(),ShouldEqual,remainAmount.String())
			}
		})

	})
}
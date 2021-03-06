package envelopes

import (
	accountService "github.com/catwithtudou/red-envelope-account/services"
	"github.com/catwithtudou/red-envelope-resk/services"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

/**
 *@Author tudou
 *@Date 2020/7/28
 **/

func TestRedEnvelopeService_SendOut(t *testing.T) {
	//发红包人的红包资金账户
	ac := accountService.GetAccountService()
	account := accountService.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "1000",
		AccountName:  "测试账户",
		AccountType:  int(accountService.EnvelopeAccountType),
		CurrencyCode: "CNY",
	}
	re := services.GetRedEnvelopeService()
	Convey("准备资金账户", t, func() {
		//准备资金账户
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO, ShouldNotBeNil)
	})
	Convey("发送红包", t, func() {
		goods := services.RedEnvelopeSendingDTO{
			UserId:       account.UserId,
			Username:     account.Username,
			EnvelopeType: services.GeneralEnvelopeType,
			Amount:       decimal.NewFromFloat(8.88),
			Quantity:     10,
			Blessing:     services.DefaultBlessing,
		}

		Convey("发普通红包", func() {
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

		})

		goods.EnvelopeType = services.LuckyEnvelopeType
		goods.Amount = decimal.NewFromFloat(88.8)

		Convey("发碰运气红包", func() {
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
			So(dto.Amount.String(), ShouldEqual, goods.Amount.String())

		})

	})
}

func TestRedEnvelopeService_SendOut_Failure(t *testing.T) {
	//发红包人的红包资金账户
	ac := accountService.GetAccountService()
	account := accountService.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户A",
		Amount:       "10",
		AccountName:  "测试账户A",
		AccountType:  int(accountService.EnvelopeAccountType),
		CurrencyCode: "CNY",
	}
	re := services.GetRedEnvelopeService()
	Convey("准备资金账户", t, func() {
		//准备资金账户
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO, ShouldNotBeNil)
	})
	Convey("发红红包	", t, func() {
		Convey("发碰运气红包", func() {
			goods := services.RedEnvelopeSendingDTO{
				UserId:       account.UserId,
				Username:     account.Username,
				EnvelopeType: services.LuckyEnvelopeType,
				Amount:       decimal.NewFromFloat(11),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			at, err := re.SendOut(goods)
			So(err, ShouldNotBeNil)
			So(at, ShouldBeNil)
			a := ac.GetEnvelopeAccountByUserId(account.UserId)
			So(a, ShouldNotBeNil)
			So(a.Balance.String(), ShouldEqual, account.Amount)

		})

	})

}

func TestRedEnvelopeService_Receive(t *testing.T) {
	//1.准备几个红包资金账户，用户发红包和收红包
	reAccountService := accountService.GetAccountService()

	Convey("收红包测试用例", t, func() {
		accounts := make([]*accountService.AccountDTO, 0)
		size := 10
		for i := 0; i < size; i++ {
			account := accountService.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i+1),
				Amount:       "2000",
				AccountName:  "测试账户" + strconv.Itoa(i+1),
				AccountType:  int(accountService.EnvelopeAccountType),
				CurrencyCode: "CNY",
			}
			//账户创建
			accountDto, err := reAccountService.CreateAccount(account)
			So(err, ShouldBeNil)
			So(accountDto, ShouldNotBeNil)

			accounts = append(accounts, accountDto)
		}

		acDto := accounts[0]
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
			for _, account := range accounts {
				rcv := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUserId:   account.UserId,
					RecvUsername: account.Username,
					AccountNo:    account.AccountNo,
				}
				item, err := re.Receive(rcv)
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				So(item.Amount, ShouldEqual, at.AmountOne)
				remainAmount = remainAmount.Sub(at.AmountOne)
				So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
			}
		})

	})
}

func TestRedEnvelopeService_Receive_Failure(t *testing.T) {
	//1. 准备几个红包资金账户，用于发红包和收红包
	reAccountService := accountService.GetAccountService()

	Convey("收红包测试用例", t, func() {
		accounts := make([]*accountService.AccountDTO, 0)
		size := 5
		for i := 0; i < size; i++ {
			account := accountService.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i+1),
				Amount:       "100",
				AccountName:  "测试账户" + strconv.Itoa(i+1),
				AccountType:  int(accountService.EnvelopeAccountType),
				CurrencyCode: "CNY",
			}
			//账户创建
			acDto, err := reAccountService.CreateAccount(account)
			So(err, ShouldBeNil)
			So(acDto, ShouldNotBeNil)
			accounts = append(accounts, acDto)
		}
		//2. 使用其中一个用户发送一个红包
		acDto := accounts[0]
		So(len(accounts), ShouldEqual, size)
		re := services.GetRedEnvelopeService()
		//发送普通红包
		goods := services.RedEnvelopeSendingDTO{
			UserId:       acDto.UserId,
			Username:     acDto.Username,
			EnvelopeType: services.LuckyEnvelopeType,
			Amount:       decimal.NewFromFloat(10),
			Quantity:     3,
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
		So(dto.Amount.String(), ShouldEqual, goods.Amount.String())

		re = services.GetRedEnvelopeService()
		Convey("收碰运气红包", func() {
			So(len(accounts), ShouldEqual, size)
			total := decimal.NewFromFloat(0)
			remainAmount := goods.Amount
			sendingAmount := decimal.NewFromFloat(0)

			for i, account := range accounts {
				rcv := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUserId:   account.UserId,
					RecvUsername: account.Username,
					AccountNo:    account.AccountNo,
				}
				if i <= 2 {
					item, err := re.Receive(rcv)
					if item != nil {
						total = total.Add(item.Amount)
					}
					logrus.Info(i+1, " ", total.String(), " ", item.Amount.String())
					So(err, ShouldBeNil)
					So(item, ShouldNotBeNil)
					remainAmount = remainAmount.Sub(item.Amount)
					So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
					a := reAccountService.GetEnvelopeAccountByUserId(rcv.RecvUserId)
					So(a, ShouldNotBeNil)
					if item.RecvUserId == goods.UserId {
						b := decimal.NewFromFloat(100)
						b = b.Sub(decimal.NewFromFloat(10))
						b = b.Add(item.Amount)
						So(a.Balance.String(), ShouldEqual, b.String())
						sendingAmount = item.Amount
					} else {
						So(a.Balance.String(), ShouldEqual, item.Amount.Add(decimal.NewFromFloat(100)).String())
					}

				} else {
					item, err := re.Receive(rcv)
					So(err, ShouldNotBeNil)
					So(item, ShouldBeNil)
				}

			}
			So(total.String(), ShouldEqual, goods.Amount.String())

			order := re.Get(at.EnvelopeNo)
			So(order, ShouldNotBeNil)
			So(order.RemainAmount.String(), ShouldEqual, "0")
			So(order.RemainQuantity, ShouldEqual, 0)
			a := reAccountService.GetEnvelopeAccountByUserId(order.UserId)
			So(a, ShouldNotBeNil)
			So(a.Balance.String(), ShouldEqual, sendingAmount.Add(decimal.NewFromFloat(90)).String())
		})

	})

}

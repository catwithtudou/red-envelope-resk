package envelopes

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"red-envelope/infra/base"
	"red-envelope/services"
	_ "red-envelope/testx"
	"testing"
	"time"
)

/**
 *@Author tudou
 *@Date 2020/7/27
 **/

//红包商品数据写入
func TestRedEnvelopeGoodsDao_GetOne(t *testing.T) {
	err:=base.Tx(func(runner *dbx.TxRunner)error {
		dao:=&RedEnvelopeGoodsDao{
			runner: runner,
		}
		now:=time.Now()
		Convey("通过编号查询普通红包商品数据",t, func() {
			goods:=&RedEnvelopeGoods{
				EnvelopeNo:     ksuid.New().Next().String(),
				EnvelopeType:   services.GeneralEnvelopeType,
				Username:       sql.NullString{String: "测试用户", Valid: true},
				UserId:         ksuid.New().Next().String(),
				Blessing:       sql.NullString{String: services.DefaultBlessing,Valid: true},
				Amount:         decimal.NewFromFloat(50),
				AmountOne:      decimal.NewFromFloat(5),
				Quantity:       50 / 5,
				RemainAmount:   decimal.NewFromFloat(50),
				RemainQuantity: 50 / 5,
				ExpiredAt:      now.Add(24 * time.Hour),
				Status:         services.OrderCreate,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.PayNothing,
			}
			id,err:=dao.Insert(goods)
			So(err,ShouldBeNil)
			So(id,ShouldBeGreaterThan,0)
			good:=dao.GetOne(goods.EnvelopeNo)
			So(good,ShouldNotBeNil)
			So(good.Amount.String(),ShouldEqual,goods.Amount.String())
			So(good.AmountOne.String(),ShouldEqual,goods.AmountOne.String())
			So(good.CreatedAt,ShouldNotBeNil)
			So(good.UpdatedAt,ShouldNotBeNil)
		})
		return nil
	})
	if err!=nil{
		logrus.Error(err)
	}
}

//更新红包剩余金额和数量
func TestRedEnvelopeGoodsDao_UpdateBalance(t *testing.T) {
	err:=base.Tx(func(runner *dbx.TxRunner)error {
		dao:=&RedEnvelopeGoodsDao{
			runner: runner,
		}
		now:=time.Now()
		Convey("通过编号查询普通红包商品数据",t, func() {
			goods:=&RedEnvelopeGoods{
				EnvelopeNo:     ksuid.New().Next().String(),
				EnvelopeType:   services.GeneralEnvelopeType,
				Username:       sql.NullString{String: "测试用户", Valid: true},
				UserId:         ksuid.New().Next().String(),
				Blessing:       sql.NullString{String: services.DefaultBlessing,Valid: true},
				Amount:         decimal.NewFromFloat(50),
				AmountOne:      decimal.NewFromFloat(5),
				Quantity:       50 / 5,
				RemainAmount:   decimal.NewFromFloat(50),
				RemainQuantity: 50 / 5,
				ExpiredAt:      now.Add(24 * time.Hour),
				Status:         services.OrderCreate,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.PayNothing,
			}
			id,err:=dao.Insert(goods)
			So(err,ShouldBeNil)
			So(id,ShouldBeGreaterThan,0)
			good:=dao.GetOne(goods.EnvelopeNo)
			So(good,ShouldNotBeNil)
			So(good.Amount.String(),ShouldEqual,goods.Amount.String())
			So(good.AmountOne.String(),ShouldEqual,goods.AmountOne.String())
			So(good.CreatedAt,ShouldNotBeNil)
			So(good.UpdatedAt,ShouldNotBeNil)
		})
		return nil
	})
	if err!=nil{
		logrus.Error(err)
	}
}
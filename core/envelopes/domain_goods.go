package envelopes

import (
	"context"
	"github.com/catwithtudou/red-envelope-infra/base"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"red-envelope/services"
	"time"
)

type goodsDomain struct {
	RedEnvelopeGoods
	item itemDomain
}

//生成一个红包编号
func (d *goodsDomain) createEnvelopeNo() {
	d.EnvelopeNo = ksuid.New().Next().String()
}

//创建一个红包商品对象
func (d *goodsDomain) Create(
	goods services.RedEnvelopeGoodsDTO) {
	d.RedEnvelopeGoods.FromDTO(&goods)
	d.RemainQuantity = goods.Quantity
	d.Username.Valid = true
	d.Blessing.Valid = true

	//根据类型区分进行计算
	if d.EnvelopeType == services.GeneralEnvelopeType {
		d.Amount = goods.AmountOne.Mul(decimal.NewFromFloat(float64(goods.Quantity)))
	} else if d.EnvelopeType == services.LuckyEnvelopeType {
		d.AmountOne = decimal.NewFromFloat(0)
	}
	d.RemainAmount = d.Amount
	//记录生成和更新时间
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	//计算过期时间
	d.ExpiredAt = now.Add(24 * time.Hour)
	//改变状态
	d.Status = services.OrderCreate
	//TODO:没有考虑OrderType和PayStatus

	//生成红包编号
	d.createEnvelopeNo()

}

//保存到红包商品表
func (d *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		id, err = dao.Insert(&d.RedEnvelopeGoods)
		return err
	})
	return id, err
}

//创建并保存红包商品
func (d *goodsDomain) CreateAndSave(ctx context.Context, goods services.RedEnvelopeGoodsDTO) (id int64, err error) {
	//创建红包商品
	d.Create(goods)
	//保存红包商品
	return d.Save(ctx)
}

//查询红包商品信息
func (d *goodsDomain) Get(envelopeNo string) (goods *RedEnvelopeGoods) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
	return
}

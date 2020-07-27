package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"red-envelope/infra/base"
	"red-envelope/services"
	"time"
)

type goodsDomain struct{
	RedEnvelopeGoods
}

//生成一个红包编号
func (d *goodsDomain) createEnvelopeNo(){
	d.EnvelopeNo = ksuid.New().Next().String()
}


//创建一个红包商品对象
func (d *goodsDomain) Create(
	goods services.RedEnvelopeGoodsDTO){
	d.RedEnvelopeGoods.FromDTO(&goods)
	d.RemainQuantity=goods.Quantity
	d.Username.Valid=true
	d.Blessing.Valid=true

	//根据类型区分进行计算
	if d.EnvelopeType == services.GeneralEnvelopeType{
		d.Amount=goods.AmountOne.Mul(decimal.NewFromFloat(float64(goods.Quantity)))
	}
	if d.EnvelopeType == services.LuckyEnvelopeType{
		d.AmountOne = decimal.NewFromFloat(0)
	}
	d.RemainAmount=goods.Amount
	//计算过期时间
	d.ExpiredAt=time.Now().Add(24 * time.Hour)
	//改变状态
	d.Status = services.OrderCreate
	//生成红包编号
	d.createEnvelopeNo()
}

//保存到红包商品表
func (d *goodsDomain)Save(ctx context.Context)(id int64,err error){
	err=base.ExecuteContext(ctx,func(runner *dbx.TxRunner) error {
		dao:=RedEnvelopeGoodsDao{runner: runner}
		id,err=dao.Insert(&d.RedEnvelopeGoods)
		return err
	})
	return id,err
}

//创建并保存红包商品
func (d *goodsDomain)CreateAndSave(ctx context.Context,goods services.RedEnvelopeGoodsDTO)(id int64,err error){
	//创建红包商品
	d.Create(goods)
	//保存红包商品
	return d.Save(ctx)


}
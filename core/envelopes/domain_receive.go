package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"red-envelope/core/accounts"
	"red-envelope/infra/algo"
	"red-envelope/infra/base"
	"red-envelope/services"
)

/**
 *@Author tudou
 *@Date 2020/7/28
 **/

var multiple = decimal.NewFromFloat(100.0)


func (d *goodsDomain)Receive(ctx context.Context,dto services.RedEnvelopeReceiveDTO)(item *services.RedEnvelopeItemDTO,err error){
	//1.创建收红包的订单明细 preCreateItem
	d.preCreateItem(dto)

	//2.查询除当前红包的剩余数量和剩余金额信息
	goods:=d.Get(dto.EnvelopeNo)

	//3. 校验剩余红包剩余金额 ：
	// 如果没有剩余，直接返回无可用红包金额
	if goods.RemainQuantity <=0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0))<=0{
		return nil,errors.New("没有足够的红包和金额")
	}


	//4.使用红包算法计算红包金额
	//nextAmount
	nextAmount,_:=d.nextAmount(goods)
	err =base.Tx(func(runner *dbx.TxRunner) error {
		dao:=RedEnvelopeGoodsDao{runner: runner}

		//5.使用乐观锁更新语句，尝试更新剩余数量和剩余金额：
		//若更新成功，也就是返回1，表示抢到红包
		//若更新失败，也就是返回0，表示无可用红包数量和金额，抢红包则失败
		rows,err:=dao.UpdateBalance(goods.EnvelopeNo,nextAmount)
		if rows<=0||err!=nil{
			return errors.New("没有足够的红包和金额了")
		}

		//6.保存订单明细数据
		d.item.Quantity=1
		//TODO:考虑PayStatus合理性
		d.item.PayStatus=int(services.Paying)
		d.item.AccountNo= dto.AccountNo
		d.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		d.item.Amount = nextAmount
		txCtx:=base.WithValueContext(ctx,runner)
		_,err=d.item.Save(txCtx)
		if err!=nil{
			return err
		}


		//7.将抢到的红包金额从系统红包中间账户转入当前用户的资金账户 : transfer
		status,err:=d.transfer(txCtx,dto)
		if status==services.TransferedStatusSuccess{
			return nil
		}

		return err
	})


	return d.item.ToDTO(),err
}

//红包转账
func (d *goodsDomain)transfer(ctx context.Context,dto services.RedEnvelopeReceiveDTO)(status services.TransferedStatus, err error){
	//获取红包中间商户
	systemAccount := base.GetSystemAccount()
	//交易主体
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}
	transfer := services.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		TradeNo:     dto.EnvelopeNo,
		Amount:      d.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包收入",
	}
	reDomain := accounts.NewAccountDomain()
	return reDomain.TransferWithContextTx(ctx, transfer)
}

//创建收红包的订单明细
func (d *goodsDomain)preCreateItem(dto services.RedEnvelopeReceiveDTO){
	d.item.AccountNo = dto.AccountNo
	d.item.EnvelopeNo = dto.EnvelopeNo
	d.item.RecvUsername = sql.NullString{String: dto.RecvUsername}
	d.item.RecvUserId = dto.RecvUserId
	d.item.createItemNo()
}

//计算红包金额
func (d *goodsDomain)nextAmount(goods *RedEnvelopeGoods)(amount decimal.Decimal,err error){
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount,nil
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		return goods.AmountOne,nil
	} else if goods.EnvelopeType == services.LuckyEnvelopeType {
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	} else {
		log.Error("不支持的红包类型")
		err=errors.New("不支持的红包类型")
	}
	return
}
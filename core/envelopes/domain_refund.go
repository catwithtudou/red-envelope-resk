package envelopes

import (
	"context"
	"errors"
	"github.com/catwithtudou/red-envelope-account/core/accounts"
	accountService "github.com/catwithtudou/red-envelope-account/services"
	"github.com/catwithtudou/red-envelope-infra/base"
	"github.com/catwithtudou/red-envelope-resk/services"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

/**
 *@Author tudou
 *@Date 2020/8/29
 **/

const (
	pageSize = 100
)

type ExpiredEnvelopeDomain struct {
	expiredGoods []RedEnvelopeGoods
	offset       int
}

//查询过期红包
func (e *ExpiredEnvelopeDomain) Next() (ok bool) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		e.expiredGoods = dao.FindExpired(e.offset, pageSize)
		if len(e.expiredGoods) >= 0 {
			e.offset += len(e.expiredGoods)
			ok = true
		}
		return nil
	})
	return
}

//处理过期红包退款
func (e *ExpiredEnvelopeDomain) Expired() (reFundEnvelopeGoods []*services.RedEnvelopeGoodsDTO, err error) {
	if e.Next() {
		reFundEnvelopeGoods = make([]*services.RedEnvelopeGoodsDTO, 0, len(e.expiredGoods))
		for _, g := range e.expiredGoods {
			logrus.Debugf("过期红包退款开始：%+v", g)
			refund, err := e.ExpiredOne(g)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			reFundEnvelopeGoods = append(reFundEnvelopeGoods, refund)
			logrus.Debugf("过期红包退款结束：%+v", g)
		}

	}
	return reFundEnvelopeGoods, err
}

//发起退款流程
func (e *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (reFundGoodsDTO *services.RedEnvelopeGoodsDTO, err error) {
	//创建一个退款订单
	refund := goods
	refund.OrderType = services.OrderTypeRefund
	//无符号导致不能使用负数
	refund.RemainAmount = goods.RemainAmount
	refund.RemainQuantity = goods.RemainQuantity
	refund.Status = services.OrderExpired
	refund.OriginEnvelopeNo = goods.EnvelopeNo
	refund.EnvelopeNo = ""
	domain := goodsDomain{RedEnvelopeGoods: refund}
	domain.createEnvelopeNo()
	refund.EnvelopeNo = domain.EnvelopeNo

	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		id, err := domain.Save(txCtx)
		if err != nil || id == 0 {
			return errors.New("创建退款订单失败")
		}

		//修改原订单订单状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpired)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}

		//调用资金账户接口进行转账
		systemAccount := base.GetSystemAccount()
		body := accountService.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}

		accountDomain := accounts.NewAccountDomain()
		account := accountDomain.GetEnvelopeAccountByUserId(goods.UserId)
		if account == nil {
			return errors.New("没有找到该用户的红包资金账户:" + goods.UserId)
		}
		target := accountService.TradeParticipator{
			AccountNo: account.AccountNo,
			UserId:    account.UserId,
			Username:  account.Username,
		}

		transfer := accountService.AccountTransferDTO{
			TradeNo:     refund.EnvelopeNo,
			TradeBody:   body,
			TradeTarget: target,
			Amount:      goods.RemainAmount,
			ChangeType:  accountService.EnvelopExpiredRefund,
			ChangeFlag:  accountService.FlagTransferOut,
			Decs:        "红包过期退款支出:" + goods.EnvelopeNo,
		}
		status, err := accountDomain.TransferWithContextTx(txCtx, transfer)
		if status != accountService.TransferedStatusSuccess {
			return errors.New("转账失败")
		}

		transfer = accountService.AccountTransferDTO{
			TradeNo:     refund.EnvelopeNo,
			TradeBody:   target,
			TradeTarget: body,
			Amount:      goods.RemainAmount,
			ChangeType:  accountService.EnvelopExpiredRefund,
			ChangeFlag:  accountService.FlagTransferIn,
			Decs:        "红包过期退款收入:" + goods.EnvelopeNo,
		}
		status, err = accountDomain.TransferWithContextTx(txCtx, transfer)
		if status != accountService.TransferedStatusSuccess {
			return errors.New("转账失败")
		}

		dao = RedEnvelopeGoodsDao{runner: runner}
		//修改原订单状态
		rows, err = dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		//修改退款订单状态
		rows, err = dao.UpdateOrderStatus(refund.EnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil || rows == 0 {
			return errors.New("更新退款订单状态失败")
		}

		//查询退款订单信息

		dao = RedEnvelopeGoodsDao{runner: runner}
		reFundGoods := dao.GetOne(refund.EnvelopeNo)
		if reFundGoods == nil {
			return errors.New("退款订单查询失败")
		}
		reFundGoodsDTO = reFundGoods.ToDTO()

		return nil
	})

	if err != nil {
		logrus.Error(err)
		return
	}

	return reFundGoodsDTO, nil
}

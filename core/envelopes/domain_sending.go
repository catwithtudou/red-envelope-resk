package envelopes

import (
	"context"
	"github.com/catwithtudou/red-envelope-account/core/accounts"
	accountService "github.com/catwithtudou/red-envelope-account/services"
	"github.com/catwithtudou/red-envelope-infra/base"
	"github.com/catwithtudou/red-envelope-resk/services"
	"github.com/tietang/dbx"
	"path"
)

/**
 *@Author tudou
 *@Date 2020/7/27
 **/

//发送红包业务领域代码
func (d *goodsDomain) SendOut(goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {

	//创建红包商品
	d.Create(goods)
	//创建活动
	activity = new(services.RedEnvelopeActivity)
	//红包链接
	//http://localhost/v1/envelope/{id}/link/
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)

	accountDomain := accounts.NewAccountDomain()

	err = base.Tx(func(runner *dbx.TxRunner) (err error) {

		ctx := base.WithValueContext(context.Background(), runner)

		//事务逻辑问题：
		//保存红包商品和红包金额的支付必须要保证全部成功或者全部失败

		//保存红包商品
		id, err := d.Save(ctx)
		if id <= 0 || err != nil {
			return err
		}
		//红包金额支付
		//1.需要红包中间商的红包资金账户，定义在配置文件中，事先初始化到资金账户表中
		//2.从红包发送人的资金账户中扣减红包金额
		body := accountService.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := accountService.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			Username:  systemAccount.Username,
			UserId:    systemAccount.UserId,
		}

		transfer := accountService.AccountTransferDTO{
			TradeBody:   body,
			TradeTarget: target,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  accountService.EnvelopeOutgoing,
			ChangeFlag:  accountService.FlagTransferOut,
			Decs:        "红包金额支付",
		}

		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status != accountService.TransferedStatusSuccess {
			return err
		}

		//3.将扣减的红包总金额转入红包中间商的红包资金账户
		//入账
		transfer = accountService.AccountTransferDTO{
			TradeBody:   target,
			TradeTarget: body,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  accountService.EnvelopeIncoming,
			ChangeFlag:  accountService.FlagTransferIn,
			Decs:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if status == accountService.TransferedStatusSuccess {
			return err
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	//旧件金额无问题则返回活动
	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()

	return activity, err
}

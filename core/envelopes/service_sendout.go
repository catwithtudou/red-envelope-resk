package envelopes

import (
	"errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"red-envelope/infra/base"
	"red-envelope/services"
	"sync"
)

var _ services.RedEnvelopeService = new(redEnvelopeService)

var once sync.Once


func init(){
	once.Do(func() {
		services.IRedEnvelopeService=new(redEnvelopeService)
	})
	
}

type redEnvelopeService struct{
	
}


//发红包
func (r *redEnvelopeService) SendOut(dto services.RedEnvelopeSendingDTO) (activity *services.RedEnvelopeActivity, err error) {
	//验证输入参数
	if err = base.ValidateStruct(&dto); err != nil {
		return activity, err
	}

	//获取红包发送人的资金账户信息
	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.UserId)
	if account == nil {
		return nil, errors.New("用户账户不存在：" + dto.UserId)
	}
	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo

	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}
	//执行发送红包的逻辑
	domain := new(goodsDomain)
	activity, err = domain.SendOut(*goods)
	if err != nil {
		log.Error(err)
	}

	return activity,err
}

//
func (r *redEnvelopeService) Receive(dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	return
}

func (r *redEnvelopeService) Refund(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	return
}

func (r *redEnvelopeService) Get(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	return
}

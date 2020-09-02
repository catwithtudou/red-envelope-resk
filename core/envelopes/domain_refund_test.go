package envelopes

import (
	"fmt"
	"github.com/catwithtudou/red-envelope-account/core/accounts"
	"github.com/catwithtudou/red-envelope-resk/services"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/**
 *@Author tudou
 *@Date 2020/8/29
 **/

func TestExpiredEnvelopeDomain_Next(t *testing.T) {
	expiredEnvelopes := ExpiredEnvelopeDomain{}
	expiredEnvelopes.Next()
	for _, v := range expiredEnvelopes.expiredGoods {
		fmt.Println(v)
	}
}

func TestExpiredEnvelopeDomain_Expired(t *testing.T) {
	expiredEnvelopes := ExpiredEnvelopeDomain{}
	expiredEnvelopes.Next()
	reExpiredEnvelopes := ExpiredEnvelopeDomain{}
	Convey("过期红包退款", t, func() {
		reFundGoods, err := reExpiredEnvelopes.Expired()
		So(err, ShouldBeNil)
		//查询资金账户
		account := accounts.NewAccountDomain()
		goodsDomain := goodsDomain{}
		for k, v := range reExpiredEnvelopes.expiredGoods {
			//验证原过期红包
			goods := goodsDomain.Get(v.EnvelopeNo)
			So(goods, ShouldNotBeNil)
			So(goods.Status, ShouldEqual, services.OrderExpiredRefundSuccessful)
			//验证退款订单
			So(reFundGoods[k], ShouldNotBeNil)
			So(reFundGoods[k].Status, ShouldEqual, services.OrderExpiredRefundSuccessful)
			So(reFundGoods[k].OrderType, ShouldEqual, services.OrderTypeRefund)
			So(reFundGoods[k].OriginEnvelope, ShouldEqual, v.EnvelopeNo)
			accountDto := account.GetEnvelopeAccountByUserId(v.UserId)
			So(accountDto, ShouldNotBeNil)
			So(accountDto.Balance, ShouldEqual, accountDto.Balance.Add(v.RemainAmount))
		}
	})

}

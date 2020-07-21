package accounts

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"red-envelope/infra/base"
	_ "red-envelope/testx"
	"testing"
)

/**
 *@Author tudou
 *@Date 2020/7/21
 **/

func TestAccountDao_GetOne(t *testing.T) {

	err:=base.Tx(func(runner *dbx.TxRunner) error {
		dao:=&AccountDao{
			runner:runner,
		}
		Convey("通过编号查询账户数据",t, func() {
			a:=&Account{
				AccountNo:   	 ksuid.New().Next().String(),
				AccountName: 	"测试资金账户",
				UserId: 		ksuid.New().Next().String(),
				Username: 		sql.NullString{String: "测试用户",Valid: true},
				Balance:      	decimal.NewFromFloat(100),
				Status:       	1,
			}
			id,err:=dao.Insert(a)
			So(err,ShouldBeNil)
			So(id,ShouldBeGreaterThan,0)
			na := dao.GetOne(a.AccountNo)
			So(na,ShouldNotBeNil)
			So(na.Balance.String(),ShouldEqual,a.Balance.String())
			So(na.CreatedAt,ShouldNotBeNil)
			So(na.UpdatedAt,ShouldNotBeNil)
		})
		return nil
	})
	if err!=nil {
		logrus.Error(err)
	}

}
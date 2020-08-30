package envelopes

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct{
	runner *dbx.TxRunner
}



//查询
func (dao *RedEnvelopeItemDao)GetOne(itemNo string)*RedEnvelopeItem{
	form:=&RedEnvelopeItem{ItemNo: itemNo}
	ok,err:=dao.runner.GetOne(form)
	if err!=nil||!ok{
		return nil
	}
	return form
}


//插入红包订单详情数据的写入
func (dao *RedEnvelopeItemDao)Insert(form *RedEnvelopeItem)(int64,error){
	rs,err:=dao.runner.Insert(form)
	if err!=nil{
		return 0,err
	}
	return rs.LastInsertId()
}

//根据envelopeNo查询详情数据
func (dao *RedEnvelopeItemDao)FindItems(envelopeNo string)(res []*RedEnvelopeItem){
	sql:="select * from red_envelope_item where envelope_no=?"
	err:=dao.runner.Find(&res,sql,envelopeNo)
	if err!=nil{
		logrus.Error(err)
	}
	return 
}

//根据envelope和userId获取订单详情信息
func (dao *RedEnvelopeItemDao)GetByUser(envelopeNo,userId string)(res *RedEnvelopeItem){
	sql:="select * from red_envelope_item where envelope_no=? and user_id = ?"
	ok,err:=dao.runner.Get(&res,sql,envelopeNo,userId)
	if !ok{
		return 
	}
	if err!=nil{
		logrus.Error(err)
		return
	}
	return
}


//根据接收者的user_id查询订单详情数据
func (dao *RedEnvelopeItemDao) ListReceivedItems(userId string, page, size int)(res []*RedEnvelopeItem){
	items := make([]*RedEnvelopeItem, 0)
	sql := "select * from red_envelope_item where recv_user_id=? order by created_at desc limit ?,?"
	err := dao.runner.Find(&items, sql, userId, page, size)
	if err != nil {
		logrus.Error(err)
		return items
	}
	return items
}


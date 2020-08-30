package users

import "time"

/**
 *@Author tudou
 *@Date 2020/8/30
 **/

type User struct {
	Id        int64     `db:"id,omitempty"`
	UserId    string    `db:"-"`
	Mobile    string    `db:"mobile,uni"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at,omitempty"` //创建时间
	UpdatedAt time.Time `db:"updated_at,omitempty"` //更新时间
}

// Code generated from user.sql by gen_sql_model. DO NOT EDIT.

package model

import (
	"time"
)

func (User) TableName() string {
	return "users"
}

type User struct {
	Id          int64     `gorm:"column:id;primary_key" json:"id"`         // 自增id
	UserName    string    `gorm:"column:user_name" json:"user_name"`       // 用户名
	GmtCreated  time.Time `gorm:"column:gmt_created" json:"gmt_created"`   // 创建时间
	GmtModified time.Time `gorm:"column:gmt_modified" json:"gmt_modified"` // 修改时间
}

func NewUser() *User {
	return &User{}
}

func (TbUser) TableName() string {
	return "tb_user"
}

type TbUser struct {
	UserId      int64     `gorm:"column:user_id;primary_key" json:"user_id"`
	Email       string    `gorm:"column:email" json:"email"`
	Password    string    `gorm:"column:password" json:"password"`
	UserName    string    `gorm:"column:user_name" json:"user_name"`
	GmtCreated  time.Time `gorm:"column:gmt_created" json:"gmt_created"`   // 创建时间
	GmtModified time.Time `gorm:"column:gmt_modified" json:"gmt_modified"` // 修改时间
}

func NewTbUser() *TbUser {
	return &TbUser{}
}

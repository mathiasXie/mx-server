// Code generated from ai_users.sql by gen_sql_model. DO NOT EDIT.

package model

import (
	"time"
)

func (AiUser) TableName() string {
	return "ai_users"
}

type AiUser struct {
	Id        int32     `gorm:"column:id;primary_key" json:"id"`   // Primary Key
	Email     string    `gorm:"column:email" json:"email"`         // 邮箱
	Password  string    `gorm:"column:password" json:"password"`   // 密码
	Phone     string    `gorm:"column:phone" json:"phone"`         // 手机号
	UserName  string    `gorm:"column:user_name" json:"user_name"` // 用户名
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewAiUser() *AiUser {
	return &AiUser{}
}

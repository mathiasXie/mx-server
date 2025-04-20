// Code generated from ai_messages.sql by gen_sql_model. DO NOT EDIT.

package model

import (
	"time"
)

func (AiMessage) TableName() string {
	return "ai_messages"
}

type AiMessage struct {
	DeviceId  int32     `gorm:"column:device_id" json:"device_id"` // 此消息的设备id
	Id        int32     `gorm:"column:id;primary_key" json:"id"`   // Primary Key
	UserId    int32     `gorm:"column:user_id" json:"user_id"`     // 此消息的用户id
	Messsage  string    `gorm:"column:messsage" json:"messsage"`   // 消息内容
	Role      string    `gorm:"column:role" json:"role"`           // 消息由谁发出,USER,ASSISTANT
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewAiMessage() *AiMessage {
	return &AiMessage{}
}

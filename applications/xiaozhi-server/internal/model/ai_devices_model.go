// Code generated from ai_devices.sql by gen_sql_model. DO NOT EDIT.

package model

import (
	"time"
)

func (AiDevice) TableName() string {
	return "ai_devices"
}

type AiDevice struct {
	BindCode   int32     `gorm:"column:bind_code" json:"bind_code"`     // 绑定码
	Id         int32     `gorm:"column:id;primary_key" json:"id"`       // Primary Key
	RoleId     int32     `gorm:"column:role_id" json:"role_id"`         // 绑定的角色ID
	UserId     int32     `gorm:"column:user_id" json:"user_id"`         // 用户id
	DeviceId   string    `gorm:"column:device_id" json:"device_id"`     // 设备id
	DeviceMac  string    `gorm:"column:device_mac" json:"device_mac"`   // 设备网络mac地址
	DeviceName string    `gorm:"column:device_name" json:"device_name"` // 设备名称
	Token      string    `gorm:"column:token" json:"token"`             // 设备token
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewAiDevice() *AiDevice {
	return &AiDevice{}
}

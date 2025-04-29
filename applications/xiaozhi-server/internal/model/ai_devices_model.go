// Code generated from ai_devices.sql by gen_sql_model. DO NOT EDIT.

package model

import (
	"time"
)

func (AiDevice) TableName() string {
	return "ai_devices"
}

type AiDevice struct {
	BindCode      int32     `gorm:"column:bind_code" json:"bind_code"`             // 绑定码
	Id            int32     `gorm:"column:id;primary_key" json:"id"`               // Primary Key
	RoleId        int32     `gorm:"column:role_id" json:"role_id"`                 // 绑定的角色ID
	UserId        int32     `gorm:"column:user_id" json:"user_id"`                 // 用户id
	BoardIp       string    `gorm:"column:board_ip" json:"board_ip"`               // 设备内网ip
	BoardSsid     string    `gorm:"column:board_ssid" json:"board_ssid"`           // 设备连接的ssid
	BoardType     string    `gorm:"column:board_type" json:"board_type"`           // 设备主板
	ChipModelName string    `gorm:"column:chip_model_name" json:"chip_model_name"` // 设备芯片
	DeviceMac     string    `gorm:"column:device_mac" json:"device_mac"`           // 设备网络mac地址
	Ip            string    `gorm:"column:ip" json:"ip"`                           // 外网ip
	Language      string    `gorm:"column:language" json:"language"`               // 语言
	Version       string    `gorm:"column:version" json:"version"`                 // 设备版本号
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewAiDevice() *AiDevice {
	return &AiDevice{}
}

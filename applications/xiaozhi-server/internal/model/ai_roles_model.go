// Code generated from ai_roles.sql by gen_sql_model. DO NOT EDIT.

package model

import (
	"time"
)

func (AiRole) TableName() string {
	return "ai_roles"
}

type AiRole struct {
	Id         int32     `gorm:"column:id;primary_key" json:"id"`         // Primary Key
	UserId     int32     `gorm:"column:user_id" json:"user_id"`           // 用户id
	Language   string    `gorm:"column:language" json:"language"`         // 角色使用的语言
	Llm        string    `gorm:"column:llm" json:"llm"`                   // 角色使用的大模型提供商
	LlmModelId string    `gorm:"column:llm_model_id" json:"llm_model_id"` // 角色使用的大模型id
	RoleDesc   string    `gorm:"column:role_desc" json:"role_desc"`       // 角色描述
	RoleName   string    `gorm:"column:role_name" json:"role_name"`       // 角色名称
	Tts        string    `gorm:"column:tts" json:"tts"`                   // 角色使用的语音合成提供商
	TtsVoiceId string    `gorm:"column:tts_voice_id" json:"tts_voice_id"` // 角色使用的语音合成声音id
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewAiRole() *AiRole {
	return &AiRole{}
}

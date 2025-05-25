package dto

import llm_proto "github.com/mathiasXie/gin-web/applications/llm-rpc/proto/pb/proto"

type UserInfo struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     *UserRole `json:"role"`

	ChatMessages []*llm_proto.ChatMessage `json:"chat_messages"`
}

type Device struct {
	Id        int    `json:"id"`
	DeviceMac string `json:"device_mac"`
	Language  string `json:"language"`
}

type UserRole struct {
	LLM        string `json:"llm"`
	LLMModelId string `json:"llm_model_id"`
	TTS        string `json:"tts"`
	TTSVoiceId string `json:"tts_voice_id"`
	Language   string `json:"language"`
	RoleDesc   string `json:"role_desc"`
}

type DeviceInfo struct {
	Id         int    `json:"id"`
	DeviceId   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceMac  string `json:"device_mac"`
	Token      string `json:"token"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeat_password" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Phone          string `json:"phone" binding:"required"`
}

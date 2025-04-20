package dto

type UserInfo struct {
	ID       int         `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     *UserRole   `json:"role"`
	Device   *DeviceInfo `json:"device"`
}

type UserRole struct {
	LLM        string `json:"llm"`
	LLMModelId string `json:"llm_model_id"`
	TTS        string `json:"tts"`
	TTSVoiceId string `json:"tts_voice_id"`
	Language   string `json:"language"`
}

type DeviceInfo struct {
	DeviceId   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceMac  string `json:"device_mac"`
	Token      string `json:"token"`
}

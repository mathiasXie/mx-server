package config

type ProviderConfig struct {
	LLM          string `yaml:"llm" json:"llm"`
	ASR          string `yaml:"asr" json:"asr"`
	TTS          string `yaml:"tts" json:"tts"`
	DefaultVoice string `yaml:"default_voice" json:"default_voice"`
}

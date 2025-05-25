package config

type ProviderConfig struct {
	LLM          string       `yaml:"llm" json:"llm"`
	DefaultModel string       `yaml:"default_model" json:"default_model"`
	ASR          string       `yaml:"asr" json:"asr"`
	TTS          string       `yaml:"tts" json:"tts"`
	DefaultVoice string       `yaml:"default_voice" json:"default_voice"`
	PromptPrefix string       `yaml:"prompt_prefix" json:"prompt_prefix"`
	Indent       IndentConfig `yaml:"indent" json:"indent"`
}

type IndentConfig struct {
	LLM   string `yaml:"llm"`
	Model string `yaml:"model"`
}

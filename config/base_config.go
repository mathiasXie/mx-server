package config

// Config 总配置结构
type Config struct {
	AppName           string         `yaml:"app_name" json:"app_name"`
	RunMode           string         `yaml:"run_mode" json:"run_mode"`
	Server            ServerConfig   `yaml:"server" json:"server"`
	Resource          ResourceConfig `yaml:"resource" json:"resource"`
	Log               LoggerConfig   `yaml:"log" json:"log"`
	AccessTokenSecret string         `yaml:"access_token_secret" json:"access_token_secret"`
	TTS               TTSConfig      `yaml:"tts" json:"tts"`
	FunctionRPC       FunctionConfig `yaml:"function_rpc" json:"function_rpc"`
	LLM               LLMConfig      `yaml:"llm" json:"llm"`
	ASR               ASRConfig      `yaml:"asr" json:"asr"`
	Provider          ProviderConfig `yaml:"provider" json:"provider"`
}

// Instance 全局配置实例
var Instance = &Config{}

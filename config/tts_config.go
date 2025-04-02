package config

// TTSConfig TTS服务配置
type TTSConfig struct {
	Microsoft MicrosoftConfig `yaml:"microsoft"`
	Doubao    DoubaoConfig    `yaml:"doubao"`
}

// MicrosoftConfig 微软TTS配置
type MicrosoftConfig struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
}

// DoubaoConfig 豆包TTS配置
type DoubaoConfig struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
}

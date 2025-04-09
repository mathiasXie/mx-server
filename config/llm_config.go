package config

// LLMConfig LLM服务配置
type LLMConfig struct {
	VolcEngine VolcEngineLLMConfig `yaml:"volcengine"`
	Aliyun     AliyunConfig        `yaml:"aliyun"`
}

// AliyunConfig 阿里云LLM配置
type AliyunConfig struct {
	APIKey         string   `yaml:"api_key"`
	BaseURL        string   `yaml:"base_url"`
	DefaultModelID string   `yaml:"default_model_id"`
	Models         []string `yaml:"models"`
}

// VolcengineConfig 火山引擎LLM配置
type VolcEngineLLMConfig struct {
	AuthToken      string   `yaml:"auth_token"`
	BaseURL        string   `yaml:"base_url"`
	DefaultModelID string   `yaml:"default_model_id"`
	Models         []string `yaml:"models"`
}

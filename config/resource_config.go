package config

// ResourceConfig 资源配置
type ResourceConfig struct {
	Mysql       MysqlConfig       `yaml:"mysql"`
	Redis       RedisConfig       `yaml:"redis"`
	TTSRPC      TTSRPCConfig      `yaml:"tts_rpc"`
	FunctionRPC FunctionRPCConfig `yaml:"function_rpc"`
	LLMRPC      LLMRPCConfig      `yaml:"llm_rpc"`
}

// MysqlConfig MySQL配置
type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Charset  string `yaml:"charset"`
	LogLevel int    `yaml:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// TTSRPCConfig TTS RPC服务配置
type TTSRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// FunctionRPCConfig Function RPC服务配置
type FunctionRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// LLMRPCConfig LLM RPC服务配置
type LLMRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

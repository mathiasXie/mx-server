package config

// ResourceConfig 资源配置
type ResourceConfig struct {
	Mysql       MysqlConfig `yaml:"mysql"`
	Redis       RedisConfig `yaml:"redis"`
	TTSRPC      RPCConfig   `yaml:"tts_rpc"`
	FunctionRPC RPCConfig   `yaml:"function_rpc"`
	LLMRPC      RPCConfig   `yaml:"llm_rpc"`
	ASRRPC      RPCConfig   `yaml:"asr_rpc"`
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

type RPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

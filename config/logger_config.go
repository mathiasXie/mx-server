package config

// LoggerConfig 日志配置
type LoggerConfig struct {
	FileDirectory   string   `yaml:"file_directory"`
	MaxSize         int      `yaml:"maxsize"`
	MaxBackups      int      `yaml:"max_backups"`
	MaxAge          int      `yaml:"max_age"`
	Compress        bool     `yaml:"compress"`
	Level           int32    `yaml:"level"`
	LogIDShowHeader bool     `yaml:"log_id_show_header"`
	SkipPaths       []string `yaml:"skip_paths"`
}

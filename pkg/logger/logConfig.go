package logger

func NewLogConfig() *LogConfig {
	return &LogConfig{}
}

func (lcfg *LogConfig) LoadConfigFromEnv() {
	lcfg.Level = getEnvOrDefault("LOG_LEVEL", "debug")
	lcfg.Environment = getEnvOrDefault("ENVIRONMENT", "development")
	lcfg.Version = getEnvOrDefault("SERVICE_VERSION", "unknown")
}

func (lcfg *LogConfig) InitializeLogConfig(serviceName string) *LogConfig {
	lcfg.LoadConfigFromEnv()
	lcfg.ServiceName = serviceName
	return lcfg
}

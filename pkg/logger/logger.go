package logger

type Logger interface {
	Init(cfg *LoggerConfig)
	Debug(msg string, params map[string]interface{})
	Info(msg string, params map[string]interface{})
	Warn(msg string, params map[string]interface{})
	Error(msg string, params map[string]interface{})
}

type LoggerConfig struct {
	JSONFormatter bool
	Level         string
	ShowMethod    bool
}

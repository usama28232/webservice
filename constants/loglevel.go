package constants

type LogLevel string

const (
	Warn  LogLevel = "warn"
	Info  LogLevel = "info"
	Debug LogLevel = "debug"
	Error LogLevel = "error"
)

var logLevelMap = map[string]LogLevel{
	"warn":  Warn,
	"info":  Info,
	"debug": Debug,
	"error": Error,
}

// Get Log level from String
//
// Returns specified or Default Log Level
func GetLogLevelFromString(level string) LogLevel {
	if status, ok := logLevelMap[level]; ok {
		return status
	}
	return Info
}

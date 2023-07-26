package loggers

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"time"

	"webservice/constants"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.SugaredLogger = nil
var accessLogger *zap.SugaredLogger = nil
var debugLogger *zap.SugaredLogger = nil

// GetLoggerByConfigFile is used to get logger type based on config file: helpers.config.json
//
// returns Info or Debug logger
func GetLoggerByConfigFile() *zap.SugaredLogger {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return defaultLogger
	}

	callerFunc := runtime.FuncForPC(pc)
	if callerFunc == nil {
		return defaultLogger
	}

	// ....
	_logger := getLoggerFromConfig(callerFunc.Name())

	return _logger
}

// RestoreDefaultLogger is used to return default logger back
//
// Returns application default logger
func RestoreDefaultLogger() *zap.SugaredLogger {
	return defaultLogger
}

// Gets Zap Logger with specified level
//
// returns new SugaredLogger instance
func GetLogger(level constants.LogLevel) *zap.SugaredLogger {
	if defaultLogger == nil {
		defaultLogger = getLogger(level, constants.LOG_FILE)
	}
	return defaultLogger
}

// Gets Zap Logger for HTTP Requests
//
// returns new SugaredLogger instance
func GetAccessLogger() *zap.SugaredLogger {
	if accessLogger == nil {
		accessLogger = getLogger(constants.LogLevel(constants.Info), constants.ACCES_LOG_FILE)
	}
	return accessLogger
}

func GetDebugLogger() *zap.SugaredLogger {
	if debugLogger == nil {
		debugLogger = getLogger(constants.LogLevel(constants.Debug), constants.LOG_FILE)
	}
	return debugLogger
}

func getLoggerFromConfig(callerName string) *zap.SugaredLogger {

	file, err := os.ReadFile(constants.DEBUG_CONFIG)
	if err != nil {
		defaultLogger.Infow("Error reading config file:", err)
		return defaultLogger
	}

	// Unmarshal the JSON data into the Config struct
	var config LoggerConfig
	if errC := json.Unmarshal(file, &config); errC != nil {
		defaultLogger.Infow("Error reading config file:", errC)
		os.Exit(1)
	}

	lastIndexGroup := strings.LastIndex(callerName, constants.CALLER_DELIM_1)

	if lastIndexGroup >= 0 {
		groupVal := callerName[lastIndexGroup+len(constants.CALLER_DELIM_1):]
		infoGrp := strings.Split(groupVal, constants.CALLER_DELIM_2)

		if len(infoGrp) >= 0 {
			value := infoGrp[0]

			if level, ok := config.LogConfig[value]; ok {
				logLevel := constants.GetLogLevelFromString(strings.ToLower(level))
				switch logLevel {
				case constants.Debug:
					return GetDebugLogger()
				default:
					return defaultLogger
				}
			} else {
				defaultLogger.Infow("Cannot find logger mapping", "v", config)
			}

		} else {
			defaultLogger.Infow("Delimeter not found ", "v:", groupVal)
		}

	} else {
		defaultLogger.Infow("Delimeter not found ", "v:", callerName)
	}

	return defaultLogger
}

func getLogger(level constants.LogLevel, filename string) *zap.SugaredLogger {
	// Configure logger options
	var l zapcore.Level
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(constants.DATETIME_FMT))
	}

	logFile, errLogFile := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if errLogFile != nil {
		defer logFile.Close()
		panic("Failed to open log file " + filename)
	}

	// Set the initial logging level
	switch level {
	case constants.Warn:
		l = zap.WarnLevel
	case constants.Debug:
		l = zap.DebugLevel
	case constants.Error:
		l = zap.ErrorLevel
	default:
		l = zap.InfoLevel
	}

	var err error
	_logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(logFile), l)).Sugar()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	return _logger
}

package loggers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"webservice/constants"
	"webservice/helpers"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.SugaredLogger = nil
var accessLogger *zap.SugaredLogger = nil
var debugLogger *zap.SugaredLogger = nil

var collection = make(map[string]*zap.SugaredLogger)

// SetLoggerFromRequest is used to set logging level Debug/Default
//
// void
func SetLoggerFromRequest(r *http.Request) {
	param, err := helpers.ParseDebugRequest(r)
	if err != nil {
		defaultLogger.Errorw("Failed to parse", "Err", err)
	}
	username := r.Header.Get(constants.USER_HEADER_KEY)
	if username != "" {
		if param.Debug {
			defaultLogger.Infow("SetLoggerFromRequest", "Parsed", param)
			initDebugLoggerbyUsername(username)
		} else {
			RemoveDebugLoggerByUsername(username)
		}
	}
}

// GetLoggerByRequest is used to get logger type based on incoming request
//
// returns Info or Debug logger
func GetLoggerByRequest(r *http.Request) *zap.SugaredLogger {
	param, err := helpers.ParseDebugRequest(r)
	if err != nil {
		defaultLogger.Errorw("Failed to parse", "Err", err)
	}
	username := r.Header.Get(constants.USER_HEADER_KEY)
	if param.Debug && username != "" {
		defaultLogger.Infow("GettingLoggerByRequest", "Parsed", param)
		return GetLoggerbyUsername(username)
	}
	return defaultLogger
}

// RemoveLoggerByRequest is used to remove logger mapping
//
// void
func RemoveLoggerByRequest(r *http.Request) {
	param, err := helpers.ParseDebugRequest(r)
	if err != nil {
		defaultLogger.Errorw("Failed to parse", "Err", err)
	}
	username := r.Header.Get(constants.USER_HEADER_KEY)
	if param.Debug && username != "" {
		RemoveDebugLoggerByUsername(username)
	}
}

func initDebugLoggerbyUsername(uname string) *zap.SugaredLogger {
	if val, ok := collection[strings.ToLower(uname)]; ok {
		return val
	} else {
		_logger := GetDebugLogger()
		_logger.Debugw("Debug logger retrieved", "User", uname)
		collection[strings.ToLower(uname)] = _logger
		return _logger
	}
}

// GetLoggerbyUsername is used to get logger based on incoming Http Header
//
// returns Debug logger or Default Logger
func GetLoggerbyUsername(uname string) *zap.SugaredLogger {
	if val, ok := collection[strings.ToLower(uname)]; ok {
		return val
	} else {
		// return default logger
		return GetLogger(constants.Info)
	}
}

// RemoveDebugLoggerByUsername is used to release logger mapping
//
// void
func RemoveDebugLoggerByUsername(uname string) {
	if uname != "" {
		delete(collection, strings.ToLower(uname))
	}
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

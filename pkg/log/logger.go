package log

import (
	"os"

	"github.com/ory/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger    *zap.Logger
	logLevels map[string]zapcore.Level
)

func init() {
	logLevels = map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"fatal": zap.FatalLevel,
		"error": zap.ErrorLevel,
	}
}

// SetLogLevel sets the log level detected from the env - LOG_LEVEL
// If the LOG_LEVEL is not set/found, the log level will default to INFO
func SetLogLevel() {
	viper.AutomaticEnv()
	level := os.Getenv("LOG_LEVEL")

	lvl, ok := logLevels[level]
	if !ok {
		bootstrap(zapcore.InfoLevel)
		return
	}

	bootstrap(lvl)
}

func bootstrap(level zapcore.Level) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	logger = zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevelAt(level),
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

// Info log entry at info level
func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

// Error log entryat error level
func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

// Fatal log entry at fatal level
func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

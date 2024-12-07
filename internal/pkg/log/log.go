package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var (
	mu  sync.Mutex
	std = NewLogger()
)

var _ Logger = (*zapLogger)(nil)

type zapLogger struct {
	z *zap.Logger
}

func Init() {
	mu.Lock()
	defer mu.Unlock()

	std = NewLogger()
}

func NewLogger() *zapLogger {

	currentDate := time.Now().Format(time.DateOnly)

	logDir := fmt.Sprintf("logs/%s", currentDate)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		Encoding:          "json",
		OutputPaths:       []string{"stdout", fmt.Sprintf("%s/json", logDir)},
		ErrorOutputPaths:  []string{"stderr", fmt.Sprintf("%s/json-err.log", logDir)},
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
		},
	}
	z, err := config.Build()
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{z: z}

	zap.RedirectStdLog(z)

	return logger
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Debugw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Infow(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Warnw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Errorw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Panicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.z.Sugar().Fatalw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Sync() {
	l.z.Sync()
}

func Sync() {
	std.z.Sugar().Sync()
}

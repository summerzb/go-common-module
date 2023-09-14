package logger

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	rotate "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapAdapt struct {
	conf   *Config
	sugar  *zap.SugaredLogger
	logger *zap.Logger
}

func newZapAdapt(conf *Config) (Logger, error) {
	adapt := &zapAdapt{conf: conf}

	if err := adapt.initialize(); err != nil {
		return nil, err
	}

	return adapt, nil
}

func (z *zapAdapt) initialize() error {
	var writers []zapcore.WriteSyncer
	// 打印到命令行
	if z.conf.StdoutPrint {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}
	// 写入文件
	fileWriter, err := output(z.conf)
	if err != nil {
		return err
	}
	if fileWriter != nil {
		writers = append(writers, zapcore.AddSync(fileWriter))
	}

	// 设置日志文件格式
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	// 格式化时间
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	if z.conf.Format == TextFormat {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), z.level())

	z.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	z.sugar = z.logger.Sugar()

	return nil
}

func (z *zapAdapt) level() zapcore.Level {
	var level zapcore.Level

	switch z.conf.Level {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	case PanicLevel:
		level = zapcore.PanicLevel
	case FatalLevel:
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	return level
}

// Debug logs a message at DebugLevel.
func (z *zapAdapt) Debug(ctx context.Context, args ...interface{}) {
	z.sugar.Debug(args...)
}

// Info logs a message at InfoLevel.
func (z *zapAdapt) Info(ctx context.Context, args ...interface{}) {
	z.sugar.Info(args...)
}

// Warn logs a message at WarnLevel.
func (z *zapAdapt) Warn(ctx context.Context, args ...interface{}) {
	z.sugar.Warn(args...)
}

// Error logs a message at ErrorLevel.
func (z *zapAdapt) Error(ctx context.Context, args ...interface{}) {
	z.sugar.Error(args...)
}

// Panic logs a message at PanicLevel.
func (z *zapAdapt) Panic(ctx context.Context, args ...interface{}) {
	z.sugar.Panic(args...)
}

// Fatal logs a message at FatalLevel.
func (z *zapAdapt) Fatal(ctx context.Context, args ...interface{}) {
	z.sugar.Fatal(args...)
}

// Debugf logs a message at DebugLevel.
func (z *zapAdapt) Debugf(ctx context.Context, format string, args ...interface{}) {
	z.sugar.Debugf(format, args...)
}

// Infof logs a message at InfoLevel.
func (z *zapAdapt) Infof(ctx context.Context, format string, args ...interface{}) {
	z.sugar.Infof(format, args...)
}

// Warnf logs a message at WarnLevel.
func (z *zapAdapt) Warnf(ctx context.Context, format string, args ...interface{}) {
	z.sugar.Warnf(format, args...)
}

// Errorf logs a message at ErrorLevel.
func (z *zapAdapt) Errorf(ctx context.Context, format string, args ...interface{}) {
	z.sugar.Errorf(format, args...)
}

// Panicf logs a message at PanicLevel.
func (z *zapAdapt) Panicf(ctx context.Context, format string, args ...interface{}) {
	z.sugar.Panicf(format, args...)
}

// Fatalf logs a message at FatalLevel.
func (z *zapAdapt) Fatalf(ctx context.Context, format string, args ...interface{}) {
	z.sugar.Fatalf(format, args)
}

// WithField adds a field to the logger.
func (z *zapAdapt) WithField(key string, value interface{}) Logger {
	clone := z.logger.With(zap.Any(key, value))

	return &zapAdapt{conf: z.conf, logger: clone, sugar: clone.Sugar()}
}

// WithFields adds multiple fields to the logger.
func (z *zapAdapt) WithFields(fields Fields) Logger {
	var items []zap.Field
	for k, v := range fields {
		items = append(items, zap.Any(k, v))
	}

	clone := z.logger.With(items...)
	return &zapAdapt{conf: z.conf, logger: clone, sugar: clone.Sugar()}
}

// output returns the io.Writer to write logs to.
func output(conf *Config) (io.Writer, error) {
	var writer io.Writer
	var err error

	if conf.Director == "" {
		return nil, err
	}

	fileName := time.Now().Format("2006-01-02") + ".log"
	if conf.RotateType == 2 {
		writer = &lumberjack.Logger{
			Filename:   path.Join(conf.Director, fileName), //日志文件的位置
			MaxSize:    int(conf.MaxSize),                  //在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: int(conf.MaxBackups),               //保留旧文件的最大个数
			MaxAge:     int(conf.MaxAge),                   //保留旧文件的最大天数
			Compress:   conf.Compress,                      //是否压缩/归档旧文件
		}
		return writer, nil
	}

	writer, err = rotate.New(
		path.Join(conf.Director, "%Y-%m-%d.log"),
		rotate.WithMaxAge(time.Duration(conf.MaxAge)*24*time.Hour),
		rotate.WithRotationTime(24*time.Hour),
	)

	return writer, err
}

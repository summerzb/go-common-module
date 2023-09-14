package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// zeroAdapt zerolog适配器
type zeroAdapt struct {
	logger zerolog.Logger
	conf   *Config
}

func newZeroAdapt(conf *Config) (Logger, error) {
	adapt := &zeroAdapt{conf: conf}
	if err := adapt.initialize(); err != nil {
		return nil, err
	}
	return adapt, nil
}

func (l *zeroAdapt) initialize() error {
	// 初始化时间格式
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.TimestampFieldName = "ts"
	// 日志信息字段名称
	zerolog.MessageFieldName = "msg"
	// 设置日志级别
	zerolog.SetGlobalLevel(l.level())

	var writes []io.Writer
	if l.conf.StdoutPrint {
		writes = append(writes, os.Stdout)
	}

	fileWriter, err := output(l.conf)
	if err != nil {
		return err
	}
	if fileWriter != nil {
		writes = append(writes, fileWriter)
	}

	if l.conf.ShowLine {
		l.logger = zerolog.New(zerolog.MultiLevelWriter(writes...)).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	} else {
		l.logger = zerolog.New(zerolog.MultiLevelWriter(writes...)).With().Timestamp().Logger()
	}

	return nil
}

func (l *zeroAdapt) Debug(ctx context.Context, args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprint(args...))
}

func (l *zeroAdapt) Info(ctx context.Context, args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(args...))
}

func (l *zeroAdapt) Warn(ctx context.Context, args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(args...))
}

func (l *zeroAdapt) Error(ctx context.Context, args ...interface{}) {
	l.logger.Error().Msg(fmt.Sprint(args...))
}

func (l *zeroAdapt) Panic(ctx context.Context, args ...interface{}) {
	l.logger.Panic().Msg(fmt.Sprint(args...))
}

func (l *zeroAdapt) Fatal(ctx context.Context, args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (l *zeroAdapt) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *zeroAdapt) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

func (l *zeroAdapt) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

func (l *zeroAdapt) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

func (l *zeroAdapt) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}

func (l *zeroAdapt) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

func (l *zeroAdapt) WithField(key string, value interface{}) Logger {
	return &zeroAdapt{
		logger: l.logger.With().Interface(key, value).Logger(),
	}
}

func (l *zeroAdapt) WithFields(fields Fields) Logger {
	logFields := make(map[string]interface{})
	for k, v := range fields {
		logFields[k] = v
	}
	return &zeroAdapt{
		logger: l.logger.With().Fields(logFields).Logger(),
	}
}

func (l *zeroAdapt) level() zerolog.Level {
	var level zerolog.Level

	switch l.conf.Level {
	case DebugLevel:
		level = zerolog.DebugLevel
	case InfoLevel:
		level = zerolog.InfoLevel
	case WarnLevel:
		level = zerolog.WarnLevel
	case ErrorLevel:
		level = zerolog.ErrorLevel
	case PanicLevel:
		level = zerolog.PanicLevel
	case FatalLevel:
		level = zerolog.FatalLevel
	default:
		level = zerolog.InfoLevel
	}

	return level
}

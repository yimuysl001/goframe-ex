package rocket

import (
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/logger"
)

var ctx = gctx.GetInitCtx()

type RocketMqLogger struct {
	Flag     string
	LevelLog string
}

func (l *RocketMqLogger) Debug(msg string, fields map[string]interface{}) {
	if l.LevelLog == "close" {
		return
	}
	if msg == "" && len(fields) == 0 {
		return
	}
	if l.LevelLog == "debug" || l.LevelLog == "all" {
		logger.Logger().Debug(ctx, msg)
	}
}

func (l *RocketMqLogger) Level(level string) {
	logger.Logger().Info(ctx, level)
}

func (l *RocketMqLogger) OutputPath(path string) (err error) {
	logger.Logger().Info(ctx, path)
	return nil
}

func (l *RocketMqLogger) Info(msg string, fields map[string]interface{}) {
	if l.LevelLog == "close" {
		return
	}
	if msg == "" && len(fields) == 0 {
		return
	}

	if l.LevelLog == "info" || l.LevelLog == "all" {
		logger.Logger().Info(ctx, msg)
	}
}

func (l *RocketMqLogger) Warning(msg string, fields map[string]interface{}) {
	if l.LevelLog == "close" {
		return
	}
	if msg == "" && len(fields) == 0 {
		return
	}

	if l.LevelLog == "warn" || l.LevelLog == "all" {
		logger.Logger().Warning(ctx, msg)
	}
}

func (l *RocketMqLogger) Error(msg string, fields map[string]interface{}) {
	if l.LevelLog == "close" {
		return
	}
	if msg == "" && len(fields) == 0 {
		return
	}
	if l.LevelLog == "error" || l.LevelLog == "all" {
		logger.Logger().Error(ctx, msg)
	}
}

func (l *RocketMqLogger) Fatal(msg string, fields map[string]interface{}) {
	if l.LevelLog == "close" {
		return
	}
	if msg == "" && len(fields) == 0 {
		return
	}

	if l.LevelLog == "fatal" || l.LevelLog == "all" {
		logger.Logger().Fatal(ctx, msg)
	}
}

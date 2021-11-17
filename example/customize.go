package main

import (
	"fmt"

	"github.com/strpc/zaptelegram"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger, _ := zap.NewProduction()

	zaptelegram.BaseAPIURL = "https://localhost:8000/botapi"
	telegramHook, _ := zaptelegram.NewTelegramHook(
		"0123456789:XXXXXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYY",
		[]int{123456789, 123456789},

		zaptelegram.WithLevel(zapcore.DebugLevel),
		//zaptelegram.WithStrongLevel(zapcore.ErrorLevel),
		zaptelegram.WithTimeout(3),
		zaptelegram.WithDisabledNotification(),
		zaptelegram.WithoutAsyncOpt(),
		zaptelegram.WithFormatter(func(e zapcore.Entry) string {
			return fmt.Sprintf("service: auth service\n%s\n%s\n%s", e.Time, e.Level, e.Message)
		}),
	)

	logger = logger.WithOptions(zap.Hooks(telegramHook.GetHook()))
	logger.Warn("first event")
	logger.Error("second event")
}

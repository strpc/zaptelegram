package main

import (
	"github.com/strpc/zaptelegram"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	telegramHook, _ := zaptelegram.NewTelegramHook(
		"xxxxxxxxxx:YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY",
		[]int{123456789},
	)
	logger = logger.WithOptions(zap.Hooks(telegramHook.GetHook()))
	logger.Info("foo")
	logger.Error("bar")
}

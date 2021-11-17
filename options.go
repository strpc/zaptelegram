package zaptelegram

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Option func(*TelegramHook) error

func WithLevel(l zapcore.Level) Option {
	return func(h *TelegramHook) error {
		levels := getLevelThreshold(l)
		h.levels = levels
		return nil
	}
}

func WithStrongLevel(l zapcore.Level) Option {
	return func(h *TelegramHook) error {
		h.levels = []zapcore.Level{l}
		return nil
	}
}

func WithDisabledNotification() Option {
	return func(h *TelegramHook) error {
		h.telegramClient.disabledNotification = true
		return nil
	}
}

func WithTimeout(t time.Duration) Option {
	return func(h *TelegramHook) error {
		h.telegramClient.httpClient.Timeout = t
		return nil
	}
}

func WithFormatter(f func(e zapcore.Entry) string) Option {
	return func(h *TelegramHook) error {
		h.telegramClient.formatter = f
		return nil
	}
}

func WithoutAsyncOpt() Option {
	return func(h *TelegramHook) error {
		h.async = false
		return nil
	}
}

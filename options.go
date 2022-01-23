package zaptelegram

import (
	"context"
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
		if h.queue {
			return AsyncOptError
		}
		h.async = false
		return nil
	}
}

func WithQueue(ctx context.Context, interval time.Duration, queueSize int) Option {
	return func(h *TelegramHook) error {
		h.async = false
		h.queue = true
		h.intervalQueue = interval
		h.entriesChan = make(chan zapcore.Entry, queueSize)
		go func() {
			_ = h.consumeEntriesQueue(ctx)
		}()
		return nil
	}
}

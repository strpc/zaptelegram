package zaptelegram

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	defaultLevel    = zapcore.WarnLevel
	defaultAsyncOpt = true
	defaultQueueOpt = false
)

var AllLevels = [6]zapcore.Level{
	zapcore.DebugLevel,
	zapcore.InfoLevel,
	zapcore.WarnLevel,
	zapcore.ErrorLevel,
	zapcore.FatalLevel,
	zapcore.PanicLevel,
}

var (
	TokenError    = errors.New("token not defined")
	ChatIDsError  = errors.New("chat ids not defined")
	AsyncOptError = errors.New("async option not worked with queue option")
)

type TelegramHook struct {
	telegramClient *telegramClient
	levels         []zapcore.Level
	async          bool
	queue          bool
	intervalQueue  time.Duration
	entriesChan    chan zapcore.Entry
}

func NewTelegramHook(token string, chatIDs []int, opts ...Option) (*TelegramHook, error) {
	if token == "" {
		return &TelegramHook{}, TokenError
	} else if len(chatIDs) == 0 {
		return &TelegramHook{}, ChatIDsError
	}
	c := newTelegramClient(token, chatIDs)
	h := &TelegramHook{
		telegramClient: c,
		levels:         []zapcore.Level{defaultLevel},
		async:          defaultAsyncOpt,
		queue:          defaultQueueOpt,
	}
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

func (h TelegramHook) GetHook() func(zapcore.Entry) error {
	return func(e zapcore.Entry) error {
		if !h.isActualLevel(e.Level) {
			return nil
		} else if h.async {
			go func() {
				_ = h.telegramClient.sendMessage(e)
			}()
			return nil
		} else if h.queue {
			h.entriesChan <- e
			return nil
		} else if err := h.telegramClient.sendMessage(e); err != nil {
			return err
		}
		return nil
	}
}

func (h TelegramHook) isActualLevel(l zapcore.Level) bool {
	for _, level := range h.levels {
		if level == l {
			return true
		}
	}
	return false
}

func (h TelegramHook) consumeEntriesQueue(ctx context.Context) error {
	ticker := time.NewTicker(h.intervalQueue)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.handleNewEntries()
		case <-ctx.Done():
			h.handleNewEntries()
			return ctx.Err()
		}
	}
}

func (h TelegramHook) handleNewEntries() {
	for len(h.entriesChan) > 0 {
		_ = h.telegramClient.sendMessage(<-h.entriesChan)
	}
}

func getLevelThreshold(l zapcore.Level) []zapcore.Level {
	for i := range AllLevels {
		if AllLevels[i] == l {
			return AllLevels[i:]
		}
	}
	return []zapcore.Level{}
}

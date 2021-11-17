package zaptelegram

import (
	"errors"

	"go.uber.org/zap/zapcore"
)

const (
	defaultLevel    = zapcore.WarnLevel
	defaultAsyncOpt = true
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
	TokenError   = errors.New("token not defined")
	ChatIDsError = errors.New("chat ids not defined")
)

type TelegramHook struct {
	telegramClient *telegramClient
	levels         []zapcore.Level
	async          bool
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
	}
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

func (h *TelegramHook) GetHook() func(zapcore.Entry) error {
	return func(e zapcore.Entry) error {
		if !h.isActualLevel(e.Level) {
			return nil
		}
		if h.async {
			go func() {
				_ = h.telegramClient.sendMessage(e)
			}()
			return nil
		}
		if err := h.telegramClient.sendMessage(e); err != nil {
			return err
		}
		return nil
	}
}

func (h *TelegramHook) isActualLevel(l zapcore.Level) bool {
	for _, level := range h.levels {
		if level == l {
			return true
		}
	}
	return false
}

func getLevelThreshold(l zapcore.Level) []zapcore.Level {
	for i := range AllLevels {
		if AllLevels[i] == l {
			return AllLevels[i:]
		}
	}
	return []zapcore.Level{}
}

package zaptelegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"
)

var (
	BaseAPIURL        = "https://api.telegram.org/bot"
	defaultLoggerName = "zaptelegram"
	RequestError      = errors.New("request error")
)

const (
	defaultDisableNotification = false
	defaultTimeout             = 10 * time.Second
)

type telegramClient struct {
	httpClient           *http.Client
	token                string
	chatIDs              []int
	disabledNotification bool
	formatter            func(e zapcore.Entry) string
}

func newTelegramClient(token string, chatIDs []int) *telegramClient {
	c := &http.Client{Timeout: defaultTimeout}
	return &telegramClient{
		httpClient:           c,
		token:                token,
		chatIDs:              chatIDs,
		disabledNotification: defaultDisableNotification,
	}
}

// Logger: zaptelegram
// 11:25:59 01.01.2007
// info
// Hello bar
func (c *telegramClient) formatMessage(e zapcore.Entry) string {
	if c.formatter != nil {
		return c.formatter(e)
	}
	loggerName := defaultLoggerName
	if e.LoggerName != "" {
		loggerName = e.LoggerName
	}
	return fmt.Sprintf("Logger: %s\n%s\n%s\n%s", loggerName, e.Time, e.Level, e.Message)
}

func (c *telegramClient) sendMessage(e zapcore.Entry) error {
	url := fmt.Sprintf("%s%s/sendMessage", BaseAPIURL, c.token)
	body := struct {
		ChatID              int    `json:"chat_id"`
		Text                string `json:"text"`
		DisableNotification bool   `json:"disable_notification"`
	}{
		Text:                c.formatMessage(e),
		DisableNotification: c.disabledNotification,
	}
	for _, chatID := range c.chatIDs {
		body.ChatID = chatID
		msg, err := json.Marshal(&body)
		if err != nil {
			return err
		}
		if err := c.post(url, msg); err != nil {
			return err
		}
	}
	return nil
}

func (c *telegramClient) post(url string, body []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req)
}

func (c *telegramClient) doRequest(req *http.Request) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode > 299 {
		return RequestError
	}
	if err := res.Body.Close(); err != nil {
		return err
	}
	return nil
}

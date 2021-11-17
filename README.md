# zaptelegram

Hook for sending events to telegram for zap logger.


Install:  
`go get -u github.com/strpc/zaptelegram`


Basic usage:
```go
package main

import (
	"github.com/strpc/zaptelegram"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	telegramHook, _ := zaptelegram.NewTelegramHook(
		"0123456789:XXXXXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYY",  // telegram token. https://t.me/BotFather
		[]int{123456789},                                  // Slice of chat_id for send events
	)
	logger = logger.WithOptions(zap.Hooks(telegramHook.GetHook()))
	
	logger.Warn("first event")  // by default hook handled level Warn and higher
	logger.Error("second event")
}
```
---

Customize usage:

```go
package main

import (
	"fmt"
	"github.com/strpc/zaptelegram"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger, _ := zap.NewProduction()

	zaptelegram.BaseAPIURL = "https://localhost:8000/botapi" // change telegram bot api url
	telegramHook, _ := zaptelegram.NewTelegramHook(
		"0123456789:XXXXXXXXXXXXXXXXXXXXYYYYYYYYYYYYYYY",
		[]int{123456789, 123456789},
		
		zaptelegram.WithLevel(zapcore.DebugLevel),  // send all events
		//zaptelegram.WithStrongLevel(zapcore.ErrorLevel),  // send only errors-events
		zaptelegram.WithTimeout(3),              // set timeout for send event to telegram. by default - 10 sec 
		zaptelegram.WithDisabledNotification(),  // disable notification for send message
		zaptelegram.WithoutAsyncOpt(),           // disable async send message. by default - enabled
		zaptelegram.WithFormatter(func(e zapcore.Entry) string {
			return fmt.Sprintf("service: auth service\n%s - %s - %s", e.Time, e.Level, e.Message)
		}),  // set custom format message.
	)

	logger = logger.WithOptions(zap.Hooks(telegramHook.GetHook()))
	logger.Warn("first event")
	logger.Error("second event")
}
```
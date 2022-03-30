package main

import (
	"context"
	"github.com/pangjianan/qq_robot/conf"
	"github.com/pangjianan/qq_robot/handler"
	"github.com/pangjianan/qq_robot/redis"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"syscall"
	"time"
)

// 消息处理器，持有 openapi 对象
var processor handler.Processor

func main() {
	conf.ConfigInit()
	redis.Init(conf.GlobalConfig)
	ctx := context.Background()
	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig("./conf/config.yaml"); err != nil {
		log.Fatalln(err)
	}

	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)

	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	processor = handler.Processor{Api: api}

	websocket.RegisterResumeSignal(syscall.SIGUSR1)
	// 根据不同的回调，生成 intents
	intent := websocket.RegisterHandlers(
		processor.ATMessageEventHandler(),
	)

	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

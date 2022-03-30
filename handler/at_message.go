package handler

import (
	"context"
	"fmt"
	"github.com/pangjianan/qq_robot/redis"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"strings"
	"time"
)

type Processor struct {
	Api openapi.OpenAPI
}

const singKey = "user:sing:%s:%d:%d" //用户打卡信息 %s=发送方uid %d=年 %d=月 bitmap

// ATMessageEventHandler 实现处理 at 消息的回调
func (p Processor) ATMessageEventHandler() websocket.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return p.processMessage(input, data)
	}
}

func (p Processor) processMessage(input string, data *dto.WSATMessageData) error {
	if input != "打卡" {
		return nil
	}
	ctx := context.Background()
	now := time.Now()
	key := fmt.Sprintf(singKey, data.Author.ID, now.Year(), now.Month())
	offset := int64(now.Day())
	_, err := redis.GlobalRedis.SetBit(ctx, key, offset, 1).Result()
	if err != nil {
		log.Printf("processMessage | SetBit fail uid=%s|key=%s|offset=%d|err=%s", data.Author.ID, key, offset, err)
		return err
	}

	count, err := redis.GlobalRedis.BitCount(ctx, key, nil).Result()
	if err != nil {
		log.Printf("processMessage | BitCount fail key=%s|err=%s", key, err)
		return err
	}
	msg := dto.MessageToCreate{
		MessageReference: &dto.MessageReference{
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
		Embed: &dto.Embed{
			Title:       "本月打卡统计",
			Description: "",
			Prompt:      "消息通知",
			Fields: []*dto.EmbedField{
				{fmt.Sprintf("本月累计打卡%d次", count), ""},
				{fmt.Sprintf("漏打卡%d次", int64(now.Day())-count), ""},
			},
		},
	}
	_, err = p.Api.PostMessage(ctx, data.ChannelID, &msg)
	if err != nil {
		log.Printf("processMessage | PostMessage fail ChannelID=%s|msg=%+v|err=%s", data.ChannelID, msg, err)
		return err
	}
	return nil
}

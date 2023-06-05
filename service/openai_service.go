package service

import (
	"context"
	"dou_yin/model/VO"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

var GptClient *openai.Client

func GetGPTMessage(msg *VO.MessageVO) (*VO.MessageVO, error) {
	resp, err := GptClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: msg.Message,
				},
			},
		},
	)
	if err != nil {
		fmt.Println("[GetGPTMessage], CreateChatCompletion err is ", err.Error())
		return nil, err
	}
	result := VO.MessageVO{
		MsgID: utils.ShiftToStringFromInt64(snowflake.GenID()),
		MsgType: msg.MsgType,
		Message: resp.Choices[0].Message.Content,
		CreateTime: utils.GetNowTime(),
		SenderID: utils.ShiftToStringFromInt64(999999),
		DataType: 0,
		IsAnonymous: false,
	}
	if result.MsgType == 0 { 	// 私聊
		result.ReceiverID = msg.SenderID
	}else{ 						// 群聊
		result.ReceiverID = msg.ReceiverID
	}
	return &result, nil
}

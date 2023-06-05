package service

import (
	"context"
	"dou_yin/dao/mysql/group_chat_dao"
	"dou_yin/model/VO"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
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
					Role:    openai.ChatMessageRoleUser,
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
		MsgID:       utils.ShiftToStringFromInt64(snowflake.GenID()),
		MsgType:     msg.MsgType,
		Message:     resp.Choices[0].Message.Content,
		CreateTime:  utils.GetNowTime(),
		SenderID:    utils.ShiftToStringFromInt64(999999),
		DataType:    0,
		IsAnonymous: false,
	}
	if result.MsgType == 0 { // 私聊
		result.ReceiverID = msg.SenderID
	} else { // 群聊
		result.ReceiverID = msg.ReceiverID
		err = CreatGroupMsg(result)
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}

// 接收新消息并保存
func CreatGroupMsg(msg VO.MessageVO) error {
	msgDO, err := DO.MGetMsgDOfromVO(&msg)
	if err != nil {
		return err
	}
	msgPO, err := DO.MGetGroupMsgPOfromDO(msgDO)
	if err != nil {
		return err
	}
	err = group_chat_dao.WriteGroupMsg(*msgPO)
	if err != nil {
		return err
	}

	return nil
}
package service

import (
	"dou_yin/dao/mysql/friend_dao"
	"dou_yin/dao/mysql/privateChat_dao"
	"dou_yin/model/PO"
	"dou_yin/model/VO"
	"dou_yin/model/VO/param"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"encoding/json"
)

func QueryPrivateChatMsg(param param.QueryPrivateChatMsgParam) (messageList DO.MessageList, err error) {
	friendship, err := friend_dao.QueryFriendshipBy2ID(utils.ShiftToNum64(param.UserID), utils.ShiftToNum64(param.FriendID))
	if err != nil {
		return messageList, err
	}

	var extra PO.FriendExtra
	if friendship.Extra != nil {
		err = json.Unmarshal([]byte(*friendship.Extra), &extra)
		if err != nil {
			return messageList, err
		}
	}

	var readTime string
	if friendship.FirstID == utils.ShiftToNum64(param.UserID) {
		readTime = extra.FirstReadTime
	} else {
		readTime = extra.SecondReadTime
	}

	MsgPOs, err := privateChat_dao.Query(friendship.FriendshipID, param.Num, param.PageNum, readTime)
	for _, msg := range MsgPOs {
		msgDO := VO.MessageVO{
			MsgID:      utils.ShiftToStringFromInt64(msg.MsgID),
			MsgType:    0,
			Message:    msg.Message,
			SenderID:   utils.ShiftToStringFromInt64(msg.SenderID),
			ReceiverID: utils.ShiftToStringFromInt64(msg.ReceiverID),
			CreateTime: msg.CreateTime,
			DataType:   msg.Type,
		}

		messageList.Messages = append(messageList.Messages, msgDO)
	}

	return messageList, nil
}

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

	MsgPOs, err := privateChat_dao.QueryByFriendshipID(friendship.FriendshipID, param.Num, param.PageNum, readTime)
	for _, msg := range MsgPOs {
		if msg.SenderID == utils.ShiftToNum64(param.UserID) {
			if msg.Deleted_list&1 == 1 {
				continue
			}
		} else {
			if msg.Deleted_list&2 == 1 {
				continue
			}
		}

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

func DeletePrivateChatMsg(param param.DeletePrivateChatMsgParam) (err error) {
	msg, err := privateChat_dao.QueryByMsgID(utils.ShiftToNum64(param.MsgID))
	if err != nil {
		return err
	}

	if msg.SenderID == utils.ShiftToNum64(param.UserID) {
		msg.Deleted_list |= 1
	} else {
		msg.Deleted_list |= 2
	}

	err = privateChat_dao.UpdateDeletedList(msg.MsgID, msg.Deleted_list)
	if err != nil {
		return err
	}

	return nil
}

func SavePrivateChatMsg(msg VO.MessageVO) (err error) {
	friendship, err := friend_dao.QueryFriendshipBy2ID(utils.ShiftToNum64(msg.SenderID), utils.ShiftToNum64(msg.ReceiverID))
	if err != nil {
		return err
	}

	po := PO.PrivateMsgPO{
		MsgID:        utils.ShiftToNum64(msg.MsgID),
		FriendshipID: friendship.FriendshipID,
		SenderID:     utils.ShiftToNum64(msg.SenderID),
		ReceiverID:   utils.ShiftToNum64(msg.ReceiverID),
		Message:      msg.Message,
		Type:         msg.DataType,
		CreateTime:   msg.CreateTime,
		Deleted_list: 0,
		Extra:        nil,
	}

	err = privateChat_dao.Insert(po)
	if err != nil {
		return err
	}

	return nil
}

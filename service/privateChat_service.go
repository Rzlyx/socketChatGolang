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
	"fmt"
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

	MsgPOs, err := privateChat_dao.QueryReadMsgByFriendshipID(friendship.FriendshipID, param.Num, param.PageNum, readTime)
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

func QueryPrivateChatMsgByDate(param param.QueryPrivateChatMsgByDateParam) (messageList DO.MessageList, err error) {
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

	MsgPOs, err := privateChat_dao.QueryReadMsgByFriendshipIDandDate(friendship.FriendshipID, readTime, param.StartTime, param.EndTime)
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

func QueryPrivateChatMsgByReadTime(param param.QueryPrivateChatMsgByReadTimeParam) (messageList DO.MessageList, err error) {
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

	MsgPOs, err := privateChat_dao.QueryReadMsgByReadTime(friendship.FriendshipID, param.ReadTime, param.Num)
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

func QueryAllUnreadPrivateChatMsg(UserID int64) (err error) {
	friendships, err := friend_dao.QueryFriendshipList(UserID)
	if err != nil {
		return err
	}

	var msgVOs []VO.MessageVO
	for _, friendship := range friendships {
		var extra PO.FriendExtra
		if friendship.Extra != nil {
			err = json.Unmarshal([]byte(*friendship.Extra), &extra)
			if err != nil {
				return err
			}
		}
		var readTime string
		if friendship.FirstID == UserID {
			readTime = extra.FirstReadTime
		} else {
			readTime = extra.SecondReadTime
		}

		fmt.Println("开始查未读消息： ", friendship.FirstID, extra.FirstReadTime, friendship.SecondID, extra.SecondReadTime, UserID, readTime)
		msgs, err := privateChat_dao.QueryUnreadMsgByFriendshipID(friendship.FriendshipID, readTime)
		fmt.Println("查到未读私聊消息：", msgs)
		if err != nil {
			return err
		}
		for _, msg := range msgs {
			msgVO := VO.MessageVO{
				MsgID:       utils.ShiftToStringFromInt64(msg.MsgID),
				MsgType:     0,
				Message:     msg.Message,
				CreateTime:  msg.CreateTime,
				SenderID:    utils.ShiftToStringFromInt64(msg.SenderID),
				ReceiverID:  utils.ShiftToStringFromInt64(UserID),
				DataType:    msg.Type,
				ErrString:   "",
				IsAnonymous: false,
			}
			fmt.Println("/*/*/*create_time ", msg)
			msgVOs = append(msgVOs, msgVO)
		}
	}

	for _, msgVO := range msgVOs {
		fmt.Println("发送未读私聊消息：", msgVO)
		MsgChan <- msgVO
	}

	return nil
}

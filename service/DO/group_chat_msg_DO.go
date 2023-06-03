package DO

import (
	"dou_yin/dao/mysql/group_chat_dao"
	"dou_yin/model/VO"
	"dou_yin/pkg/utils"
	"encoding/json"
	"fmt"
)

type GroupMsgDO struct {
	MsgID       int64          `json:"message_id" db:"message_id"`
	GroupID     int64          `json:"group_id" db:"group_id"`
	SenderID    int64          `json:"sender_id" db:"sender_id"`
	Message     string         `json:"message" db:"message"`
	Type        int            `json:"type" db:"type"`
	IsAnonymous bool           `json:"is_anonymous" db:"is_anonymous"`
	CreateTime  string         `json:"create_time" db:"create_time"`
	DeletedList *[]int64       `json:"deleted_list" db:"deleted_list"`
	Extra       *GroupMsgExtra `json:"extra" db:"extra"`
}

type GroupMsgExtra struct {
	Path string `json:"path"`
}

func MGetGroupMsgDOfromPO(msg *group_chat_dao.GroupMsgPO) (*GroupMsgDO, error) {
	var deletedList []int64
	if msg.DeletedList != nil {
		err := json.Unmarshal([]byte(*msg.DeletedList), &deletedList)
		if err != nil {
			fmt.Println("[MGetGroupMsgDOfromPO], Unmarshal deletedList err is ", err.Error())
			return nil, err
		}
	}
	var extra GroupMsgExtra
	if msg.Extra != nil {
		err := json.Unmarshal([]byte(*msg.Extra), &extra)
		if err != nil {
			fmt.Println("[MGetGroupMsgDOfromPO], Unmarshal extra err is ", err.Error())
			return nil, err
		}
	}

	return &GroupMsgDO{
		MsgID:       msg.MsgID,
		GroupID:     msg.GroupID,
		SenderID:    msg.SenderID,
		Message:     msg.Message,
		Type:        msg.Type,
		IsAnonymous: msg.IsAnonymous,
		CreateTime:  msg.CreateTime,
		DeletedList: &deletedList,
		Extra:       &extra,
	}, nil
}

func MGetGroupMsgPOfromDO(msg *GroupMsgDO) (*group_chat_dao.GroupMsgPO, error) {
	result := group_chat_dao.GroupMsgPO{
		MsgID:       msg.MsgID,
		GroupID:     msg.GroupID,
		SenderID:    msg.SenderID,
		Message:     msg.Message,
		Type:        msg.Type,
		IsAnonymous: msg.IsAnonymous,
		CreateTime:  msg.CreateTime,
	}

	result.DeletedList = nil
	if msg.DeletedList != nil {
		data, err := json.Marshal(*msg.DeletedList)
		if err != nil {
			fmt.Println("[MGetGroupMsgPOfromDO], Marshal deletedList err is ", err.Error())
			return nil, err
		}
		if len(*msg.DeletedList) > 0 {
			deletedList := string(data)
			result.DeletedList = &deletedList
		}
	}

	result.Extra = nil
	if msg.Extra != nil {
		data, err := json.Marshal(*msg.Extra)
		if err != nil {
			fmt.Println("[MGetGroupMsgPOfromDO], Marshal extra err is ", err.Error())
			return nil, err
		}
		extra := string(data)
		result.Extra = &extra
	}

	return &result, nil
}

func MGetMsgVOfromDO(msg *GroupMsgDO, Type int) *VO.MessageVO {
	return &VO.MessageVO{
		MsgID: utils.ShiftToStringFromInt64(msg.MsgID),
		MsgType: Type,
		Message: msg.Message,
		CreateTime: msg.CreateTime,
		SenderID: utils.ShiftToStringFromInt64(msg.SenderID),
		ReceiverID: utils.ShiftToStringFromInt64(msg.GroupID),
		DataType: msg.Type,
	}
}

func MGetMsgDOfromVO(msg *VO.MessageVO) (*GroupMsgDO, error) {
	result := GroupMsgDO{
		MsgID: utils.ShiftToNum64(msg.MsgID),
		GroupID: utils.ShiftToNum64(msg.ReceiverID),
		SenderID: utils.ShiftToNum64(msg.SenderID),
		Message: msg.Message,
		Type: msg.MsgType,
		IsAnonymous: msg.IsAnonymous,
		CreateTime: msg.CreateTime,
	}
	var deletedList []int64
	result.DeletedList = &deletedList
	result.Extra = &GroupMsgExtra{}
	
	return &result, nil
}
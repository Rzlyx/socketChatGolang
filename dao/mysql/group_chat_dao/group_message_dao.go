package group_chat_dao

import (
	"dou_yin/dao/mysql"
	"fmt"
)

// 获取历史消息
func MGetGroupMsgOldList(GroupID int64, ReadTime string, pageNum, num int) (*[]GroupMsgPO, error) {
	var list []GroupMsgPO
	startIndex := num * pageNum
	strSql := "select * from group_message where group_id = ? and create_time < ? limit ?, ?"
	err := mysql.DB.Select(&list, strSql, GroupID, ReadTime, startIndex, num)
	if err != nil {
		fmt.Println("[MGetGroupMsgOldList], select old err ", err.Error())
		return nil, err
	}
	return &list, nil
}

// 获取未读消息
func MGetGroupNewList(GroupID int64, ReadTime string, pageNum, num int) (*[]GroupMsgPO, error) {
	var list []GroupMsgPO
	startIndex := num * pageNum
	strSql := "select * from group_message where group_id = ? and create_time > ? limit ?, ?"
	err := mysql.DB.Select(&list, strSql, GroupID, ReadTime, startIndex, num)
	if err != nil {
		fmt.Println("[MGetGroupMsgOldList], select old err ", err.Error())
		return nil, err
	}
	return &list, nil
}

// 写入新消息
func WriteGroupMsg(msg GroupMsgPO) error {
	strSql := "insert group_message (message_id, group_id, sender_id, message, type, is_anonymous, create_time, deleted_list, extra) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := mysql.DB.Exec(strSql, 
		msg.MsgID,
		msg.GroupID,
		msg.SenderID,
		msg.Message,
		msg.Type,
		msg.IsAnonymous,
		msg.CreateTime,
		msg.DeletedList,
		msg.Extra,
	)
	if err != nil {
		fmt.Println("[WriteGroupMsg], insert err is ", err.Error())
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		fmt.Println("[WriteGroupMsg], insert err is ", err.Error())
		return err
	}
	return nil
}

// 更新消息
func UpdateMsgbyPO(msg GroupMsgPO) (error) {
	strSql := "UPDATE group_message SET message_id = ?, group_id = ?, sender_id = ?, message = ?, type = ?, is_anonymous = ?, create_time = ?, deleted_list = ?, extra = ? WHERE message_id = ?"
	_, err := mysql.DB.Exec(strSql,
		msg.MsgID,
        msg.GroupID,
        msg.SenderID,
        msg.Message,
        msg.Type,
        msg.IsAnonymous,
        msg.CreateTime,
        msg.DeletedList,
        msg.Extra,
        msg.MsgID,
	)
	if err != nil {
		return err
	}
	return nil
}
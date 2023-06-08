package privateChat_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"fmt"
)

func Insert(privateMsgPO PO.PrivateMsgPO) (err error) {
	sqlStr := "insert into private_message(message_id, friendship_id, sender, receiver, message, type, create_time) values(?, ?, ?, ?, ?, ?, ?)"
	_, err = mysql.DB.Exec(sqlStr, privateMsgPO.MsgID, privateMsgPO.FriendshipID, privateMsgPO.SenderID, privateMsgPO.ReceiverID, privateMsgPO.Message, privateMsgPO.Type, privateMsgPO.CreateTime)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func QueryReadMsgByFriendshipID(friendshipID int64, num int, pageNum int, readTime string) (privateMsgPOs []PO.PrivateMsgPO, err error) {
	startIndex := num * pageNum
	sqlStr := "select * from private_message where friendship_id = ? and create_time < ? limit ?, ?"
	err = mysql.DB.Select(&privateMsgPOs, sqlStr, friendshipID, readTime, startIndex, num)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return privateMsgPOs, nil
}

func QueryReadMsgByFriendshipIDandDate(friendshipID int64, readTime string, startTime string, endTime string) (privateMsgPOs []PO.PrivateMsgPO, err error) {
	sqlStr := "select * from private_message where friendship_id = ? and create_time < ? and create_time < ? and create_time > ?"
	err = mysql.DB.Select(&privateMsgPOs, sqlStr, friendshipID, readTime, endTime, startTime)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return privateMsgPOs, nil
}

func QueryReadMsgByReadTime(friendshipID int64, lastTime string, num int) (privateMsgPOs []PO.PrivateMsgPO, err error) {
	sqlStr := "select * from private_message where friendship_id = ? and create_time < ? limit ?"
	err = mysql.DB.Select(&privateMsgPOs, sqlStr, friendshipID, lastTime, num)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return privateMsgPOs, nil
}

func QueryUnreadMsgByFriendshipID(friendshipID int64, readTime string) (privateMsgPOs []PO.PrivateMsgPO, err error) {
	sqlStr := "select * from private_message where friendship_id = ? and create_time > ?"
	err = mysql.DB.Select(&privateMsgPOs, sqlStr, friendshipID, readTime)
	fmt.Println("查询私聊消息(历史未读):", readTime, privateMsgPOs)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return privateMsgPOs, nil
}

func QueryByMsgID(MsgID int64) (po PO.PrivateMsgPO, err error) {
	sqlStr := "select * from private_message where message_id = ?"
	err = mysql.DB.Get(&po, sqlStr, MsgID)
	if err != nil {
		logger.Log.Error(err.Error())
		return po, err
	}

	return po, nil
}

func UpdateDeletedList(msgID int64, deleted_list int) (err error) {
	sqlStr := "update private_message set deleted_list = ? where message_id = ?"
	_, err = mysql.DB.Exec(sqlStr, deleted_list, msgID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return err
}

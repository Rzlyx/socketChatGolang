package privateChat_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
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

// todo: deleted list
func Query(friendshipID int64, num int, pageNum int, readTime string) (privateMsgPOs []PO.PrivateMsgPO, err error) {
	startIndex := num * pageNum
	sqlStr := "select * from private_message where friendship_id = ? and create_time < ? limit ?, ?"
	err = mysql.DB.Select(&privateMsgPOs, sqlStr, friendshipID, readTime, startIndex, num)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return privateMsgPOs, nil
}

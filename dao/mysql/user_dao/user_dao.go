package user_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
)

// todo: concurrent

func QueryUserInfo(userID int64) (userPO PO.UserPO, err error) {
	sqlStr := "select * from user where user_id = ? "
	err = mysql.DB.Get(&userPO, sqlStr, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return userPO, err
	}

	return userPO, nil
}

func UpdatePrivateChatWhite(userID int64, whiteList string) (err error) {
	sqlStr := "update user set private_chat_white = ? where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdatePrivateChatBlack(userID int64, blackList string) (err error) {
	sqlStr := "update user set private_chat_black = ? where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleWhite(userID int64, whiteList string) (err error) {
	sqlStr := "update user set friend_circle_white = ? where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleBlack(userID int64, blackList string) (err error) {
	sqlStr := "update user set friend_circle_black = ? where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdatePrivateChatBlackWhite(userID int64, whiteList string, blackList string) (err error) {
	sqlStr := "update user set private_chat_white = ? and private_chat_black where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, whiteList, blackList, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleBlackWhite(userID int64, whiteList string, blackList string) (err error) {
	sqlStr := "update user set friend_circle_white = ? and friend_circle_black where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, whiteList, blackList, userID)
	if err != nil {
		return err
	}

	return nil
}

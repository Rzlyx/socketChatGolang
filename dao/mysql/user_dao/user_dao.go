package user_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"fmt"
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
	if whiteList == "" {
		sqlStr := "update user set private_chat_white = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	} else {
		sqlStr := "update user set private_chat_white = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	}

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdatePrivateChatBlack(userID int64, blackList string) (err error) {
	if blackList == "" {
		sqlStr := "update user set private_chat_black = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	} else {
		sqlStr := "update user set private_chat_black = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleWhite(userID int64, whiteList string) (err error) {
	if whiteList == "" {
		sqlStr := "update user set friend_circle_white = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	} else {
		sqlStr := "update user set friend_circle_white = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleBlack(userID int64, blackList string) (err error) {
	if blackList == "" {
		sqlStr := "update user set friend_circle_black = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	} else {
		sqlStr := "update user set friend_circle_black = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdatePrivateChatBlackWhite(userID int64, whiteList string, blackList string) (err error) {
	if whiteList == "" && blackList == "" {
		sqlStr := "update user set private_chat_white = NULL, private_chat_black  = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	} else if whiteList == "" && blackList != "" {
		sqlStr := "update user set private_chat_white = NULL, private_chat_black  = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	} else if whiteList != "" && blackList == "" {
		sqlStr := "update user set private_chat_white = ?, private_chat_black  = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	} else {
		sqlStr := "update user set private_chat_white = ?, private_chat_black  = ? where user_id = ?"
		fmt.Println(whiteList, blackList)
		_, err = mysql.DB.Exec(sqlStr, whiteList, blackList, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleBlackWhite(userID int64, whiteList string, blackList string) (err error) {
	if whiteList == "" && blackList == "" {
		sqlStr := "update user set friend_circle_white = NULL , friend_circle_black  = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	} else if whiteList == "" && blackList != "" {
		sqlStr := "update user set friend_circle_white = NULL , friend_circle_black  = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	} else if whiteList != "" && blackList == "" {
		sqlStr := "update user set friend_circle_white = ? , friend_circle_black  = NULL where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	} else {
		sqlStr := "update user set friend_circle_white = ? , friend_circle_black  = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, blackList, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

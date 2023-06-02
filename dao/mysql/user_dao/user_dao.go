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

func CreateUserInfoByPO(user *PO.UserPO) error {
	strSql := "INSERT INTO user (user_id, user_name, password, sex, phone_number, e_mail, signature, birthday, status, private_chat_white, private_chat_black, friend_circle_white, friend_circle_black, friend_circle_visiable, group_chat_white, group_chat_black, group_chat_gray, create_time, is_deleted, extra ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := mysql.DB.Exec(strSql,
		user.UserID,
		user.UserName,
		user.Password,
		user.Sex,
		user.PhoneNumber,
		user.Email,
		user.Signature,
		user.Birthday,
		user.Status,
		user.PrivateChatWhite,
		user.PrivateChatBlack,
		user.FriendCircleWhite,
		user.FriendCircleBlack,
		user.FriendCircleVisiable,
		user.GroupChatWhite,
		user.GroupChatBlack,
		user.GroupChatGray,
		user.CreateTime,
		user.IsDeleted,
		user.Extra)
	if err != nil {
		fmt.Println("[CreateUserInfoByPO], insert err is ", err.Error())
		return err
	}
	return nil
}

func UpdateUserInfoByPO(user *PO.UserPO) error {
	strSql := "UPDATE user SET user_name = ?, password = ?, sex = ?, phone_number = ?, e_mail = ?, signature = ?, birthday = ?, status = ?, private_chat_white = ?, private_chat_black = ?, friend_circle_white = ?, friend_circle_black = ?, friend_circle_visiable = ?, group_chat_white = ?, group_chat_black = ?, group_chat_gray = ?, create_time = ?, is_deleted = ?, extra = ? WHERE user_id = ?"
	_, err := mysql.DB.Exec(strSql,
		user.UserName,
		user.Password,
		user.Sex,
		user.PhoneNumber,
		user.Email,
		user.Signature,
		user.Birthday,
		user.Status,
		user.PrivateChatWhite,
		user.PrivateChatBlack,
		user.FriendCircleWhite,
		user.FriendCircleBlack,
		user.FriendCircleVisiable,
		user.GroupChatWhite,
		user.GroupChatBlack,
		user.GroupChatGray,
		user.CreateTime,
		user.IsDeleted,
		user.Extra,
		user.UserID)
	if err != nil {
		fmt.Println("[UpdateUserInfoByPO], update err is ", err.Error())
		return err
	}
	return nil
}

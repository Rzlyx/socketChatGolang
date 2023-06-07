package user_dao

import (
	"database/sql"
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"fmt"
)

// todo: concurrent
func Register(p *PO.UserPO) error {
	sqlStr := "insert into user (user_id, user_name, password, sex, phone_number, email, signature, birthdaay) VALUES (?,?,?,?,?,?,?,?)"
	_, err := mysql.DB.Exec(sqlStr, p.UserID, p.UserName, p.Password, p.Sex, p.PhoneNumber, p.Email, p.Signature, p.Birthday)
	if err != nil {
		fmt.Println("[Register], insert err is ", err)
	}
	return err
}

func Login(username string) (p1 *PO.UserPO, err error) {
	p1 = new(PO.UserPO)
	sqlStr := "select * from user where user_name = ?"
	err = mysql.DB.Get(p1, sqlStr, username)
	return p1, err
}

func GetContactorList(Id string) (p *PO.ContactorList, err error) {
	p = new(PO.ContactorList)
	sqlStr := "select * from `" + Id + "`"
	err = mysql.DB.Select(&p.ContactorList, sqlStr)
	fmt.Println(err)
	return p, err
}

func QueryUserInfo(userID int64) (userPO PO.UserPO, err error) {
	sqlStr := "select * from user where user_id = ? "
	err = mysql.DB.Get(&userPO, sqlStr, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return userPO, err
	}

	return userPO, nil
}

func UpdatePrivateChatWhite(tx *sql.Tx, userID int64, whiteList string) (err error) {
	if whiteList != "" {
		sqlStr := "update user set private_chat_white = ? where user_id = ?"
		_, err = tx.Exec(sqlStr, whiteList, userID)
	} else {
		sqlStr := "update user set private_chat_white = null where user_id = ?"
		_, err = tx.Exec(sqlStr, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdatePrivateChatBlack(tx *sql.Tx, userID int64, blackList string) (err error) {
	if blackList != "" {
		sqlStr := "update user set private_chat_black = ? where user_id = ?"
		_, err = tx.Exec(sqlStr, blackList, userID)
	} else {
		sqlStr := "update user set private_chat_black = null where user_id = ?"
		_, err = tx.Exec(sqlStr, userID)
	}

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleWhite(tx *sql.Tx, userID int64, whiteList string) (err error) {
	if whiteList == "" {
		sqlStr := "update user set friend_circle_white = null where user_id = ?"
		_, err = tx.Exec(sqlStr, userID)
	} else {
		sqlStr := "update user set friend_circle_white = ? where user_id = ?"
		_, err = tx.Exec(sqlStr, whiteList, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleBlack(tx *sql.Tx, userID int64, blackList string) (err error) {
	if blackList != "" {
		sqlStr := "update user set friend_circle_black = ? where user_id = ?"
		_, err = tx.Exec(sqlStr, blackList, userID)
	} else {
		sqlStr := "update user set friend_circle_black = null where user_id = ?"
		_, err = tx.Exec(sqlStr, userID)
	}

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdatePrivateChatBlackWhite(userID int64, whiteList string, blackList string) (err error) {
	fmt.Println(whiteList, blackList)
	if whiteList != "" && blackList != "" {
		fmt.Println("*")
		sqlStr := "update user set private_chat_white = ?, private_chat_black = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, blackList, userID)
	} else if whiteList == "" && blackList != "" {
		fmt.Println("**")
		sqlStr := "update user set private_chat_white = null, private_chat_black = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	} else if whiteList != "" && blackList == "" {
		fmt.Println("***")
		sqlStr := "update user set private_chat_white = ?, private_chat_black = null where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	} else {
		fmt.Println("****")
		sqlStr := "update user set private_chat_white = null, private_chat_black = null where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFriendCircleBlackWhite(userID int64, whiteList string, blackList string) (err error) {
	if whiteList != "" && blackList != "" {
		sqlStr := "update user set friend_circle_white = ?, friend_circle_black = ? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, blackList, userID)
	} else if whiteList == "" && blackList != "" {
		sqlStr := "update user set friend_circle_white = null, friend_circle_black =? where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, blackList, userID)
	} else if whiteList != "" && blackList == "" {
		sqlStr := "update user set friend_circle_white = ?, friend_circle_black = null where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, whiteList, userID)
	} else {
		sqlStr := "update user set friend_circle_white = null, friend_circle_black = null where user_id = ?"
		_, err = mysql.DB.Exec(sqlStr, userID)
	}
	if err != nil {
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

func UpdateUserInfoByPO(tx *sql.Tx, user *PO.UserPO) error {
	fmt.Printf("func(UpdateUserInfoByPO): [param: %v]\n", user)
	strSql := "UPDATE user SET user_name = ?, password = ?, sex = ?, phone_number = ?, e_mail = ?, signature = ?, birthday = ?, status = ?, private_chat_white = ?, private_chat_black = ?, friend_circle_white = ?, friend_circle_black = ?, friend_circle_visiable = ?, group_chat_white = ?, group_chat_black = ?, group_chat_gray = ?, create_time = ?, is_deleted = ?, extra = ? WHERE user_id = ?"
	//fmt.Println("func(UpdateUserInfoByPO): ", strSql)
	_, err := tx.Exec(strSql,
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

func UpdateUserInfoByPOTx(tx *sql.Tx, user *PO.UserPO) error {
	strSql := "UPDATE user SET user_name = ?, password = ?, sex = ?, phone_number = ?, e_mail = ?, signature = ?, birthday = ?, status = ?, private_chat_white = ?, private_chat_black = ?, friend_circle_white = ?, friend_circle_black = ?, friend_circle_visiable = ?, group_chat_white = ?, group_chat_black = ?, group_chat_gray = ?, create_time = ?, is_deleted = ?, extra = ? WHERE user_id = ?"
	_, err := tx.Exec(strSql,
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

func QueryLike(str string) (users []PO.UserPO, err error) {
	sqlStr := "select * from user where user_name like '%" + str + "%'"
	err = mysql.DB.Select(&users, sqlStr)
	if err != nil {
		return users, err
	}

	return users, nil
}

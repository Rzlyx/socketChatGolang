package user_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/model/PO"
	"encoding/json"
)

func QueryUserInfo(userID int64) (userPO PO.UserPO, err error) {
	sqlStr := "select * from user where user_id = ? "
	err = mysql.DB.QueryRow(sqlStr, userID).Scan(&userPO)
	if err != nil {
		return userPO, err
	}

	return userPO, nil
}

func AddPrivateChatWhite(userID int64, friendID int64) (err error) {
	var userPO PO.UserPO
	var whileList PO.PrivateChatWhite
	sqlStr := "selete * from user where user_id = ?"
	err = mysql.DB.QueryRow(sqlStr, userID).Scan(&userPO)
	if err != nil {
		return err
	}

	if userPO.PrivateChatWhite != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whileList)
		if err != nil {
			return nil
		}
	}
	whileList.WhiteList = append(whileList.WhiteList, friendID)
	whiteJson, err := json.Marshal(whileList)
	if err != nil {
		return nil
	}
	WhiteStr := string(whiteJson[:])

	sqlStr = "update user set private_chat_white = ? where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, WhiteStr, userID)
	if err != nil {
		return err
	}

	return nil
}

func AddPrivateChatWhiteFromBlack(userID int64, friendID int64) (err error) {
	var userPO PO.UserPO
	var whileList PO.PrivateChatWhite
	var blackList PO.PrivateChatBlack
	sqlStr := "select * from user where user_id = ?"
	err = mysql.DB.QueryRow(sqlStr, userID).Scan(&userPO)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whileList)
	if err != nil {
		return err
	}

	whileList.WhiteList = append(whileList.WhiteList, friendID)
	whileJson, err := json.Marshal(whileList)
	if err != nil {
		return err
	}
	whileStr := string(whileJson[:])

	err = json.Unmarshal([]byte(*userPO.PrivateChatBlack), &blackList)
	if err != nil {
		return err
	}
	for index, id := range blackList.BlackList {
		if id == friendID {
			blackList.BlackList = append(blackList.BlackList[:index], blackList.BlackList[index+1:]...)
			break
		}
	}
	blackJson, err := json.Marshal(blackList)
	if err != nil {
		return err
	}
	blackStr := string(blackJson[:])

	sqlStr = "update user set private_chat_white = ? and private_char_black where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, whileStr, blackStr, userID)
	if err != nil {
		return err
	}

	return nil
}

func AddPrivateChatBlack(userID int64, friendID int64) (err error) {
	var userPO PO.UserPO
	var blackList PO.PrivateChatBlack
	var whiteList PO.PrivateChatWhite
	sqlStr := "select * from user where user_id = ?"
	err = mysql.DB.QueryRow(sqlStr, userID).Scan(&userPO)
	if userPO.PrivateChatBlack != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatBlack), &blackList)
		if err != nil {
			return err
		}
	}
	blackList.BlackList = append(blackList.BlackList, friendID)
	blackJson, err := json.Marshal(blackList)
	if err != nil {
		return err
	}
	blackStr := string(blackJson[:])

	err = json.Unmarshal([]byte(*userPO.PrivateChatBlack), &whiteList)
	if err != nil {
		return err
	}
	for index, id := range whiteList.WhiteList {
		if id == friendID {
			whiteList.WhiteList = append(whiteList.WhiteList[:index], whiteList.WhiteList[index+1:]...)
			break
		}
	}
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		return err
	}
	whiteStr := string(whiteJson[:])

	sqlStr = "update user set private_chat_black = ? and private_chat white where user_id = ?"
	_, err = mysql.DB.Exec(sqlStr, blackStr, whiteStr, userID)
	if err != nil {
		return err
	}

	return nil
}

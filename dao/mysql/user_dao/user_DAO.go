package user_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/model/PO"
)

func QueryUserInfo(userID int64) (userPO PO.UserPO, err error) {
	sqlStr := "select * from user where user_id = ? "
	err = mysql.DB.QueryRow(sqlStr, userID).Scan(userPO)
	if err != nil {
		return userPO, err
	}

	return userPO, nil
}

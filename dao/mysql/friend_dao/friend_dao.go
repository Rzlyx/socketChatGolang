package friend_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"dou_yin/service/DO"
	"errors"
	"fmt"
)

func QueryFriendshipList(userID int64) (friends []PO.FriendPO, err error) {
	sqlStr := "select * from friend where first_id = ? or second_id = ?"
	err = mysql.DB.Select(&friends, sqlStr, userID, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return friends, nil
}

func QueryFriendship(friendShipID int64) (friend PO.FriendPO, err error) {
	sqlStr := "select * from friend where friendship_id = ?"
	err = mysql.DB.Get(&friend, sqlStr, friendShipID)
	if err != nil {
		logger.Log.Error(err.Error())
		return friend, err
	}

	return friend, nil
}

func QueryFriendshipBy2ID(firstID int64, secondID int64) (friend PO.FriendPO, err error) {
	fmt.Println(firstID, secondID)
	sqlStr := "select * from friend where (first_id = ? and second_id = ?) or (first_id = ? and second_id = ?)"
	err = mysql.DB.Get(&friend, sqlStr, firstID, secondID, secondID, firstID)
	if err != nil {
		logger.Log.Error(err.Error())
		return friend, err
	}

	return friend, nil
}

func DeleteFriend(friendshipID int64) (err error) {
	sqlStr := "delete from friend where friendship_id = ?"
	ret, err := mysql.DB.Exec(sqlStr, friendshipID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	rows, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("not friend")
	}

	return nil
}

func Insert(friendship DO.Friendship) (err error) {
	sqlStr := "insert friend(friendship_id, first_id, second_id, f_remark_s, s_remark_f, is_f_remark_s, is_s_remark_f) values(?, ?, ?, ?, ?, ?, ?)"
	_, err = mysql.DB.Exec(sqlStr, friendship.FriendshipID, friendship.FirstID, friendship.SecondID, friendship.FirstRemarkSecond, friendship.SecondRemarkFirst, false, false)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateFirstRemarkSecond(friendshipID int64, remark string, realName string) (err error) {
	if remark == "" {
		sqlStr := "update friend set f_remark_s = ?, is_f_remark_s = ? where friendship_id = ?"
		_, err = mysql.DB.Exec(sqlStr, realName, false, friendshipID)
	} else {
		sqlStr := "update friend set f_remark_s = ?, is_f_remark_s = ? where friendship_id = ?"
		_, err = mysql.DB.Exec(sqlStr, remark, true, friendshipID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateSecondRemarkFirst(friendshipID int64, remark string, realName string) (err error) {
	if remark == "" {
		sqlStr := "update friend set s_remark_f = ? , is_s_remark_f = ? where friendship_id = ?"
		_, err = mysql.DB.Exec(sqlStr, realName, false, friendshipID)
	} else {
		sqlStr := "update friend set s_remark_f = ? , is_s_remark_f = ? where friendship_id = ?"
		_, err = mysql.DB.Exec(sqlStr, remark, true, friendshipID)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateExtra(friendshipID int64, extra string) (err error) {
	sqlStr := "update friend set extra = ? where friendship_id = ?"
	_, err = mysql.DB.Exec(sqlStr, extra, friendshipID)
	if err != nil {
		return err
	}

	return nil
}

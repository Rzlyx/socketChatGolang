package friend_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/model/PO"
	"dou_yin/service/DO"
	"errors"
)

func QueryFriendshipList(userID int64) (friends []PO.FriendPO, err error) {
	sqlStr := "select * from friend where first_id = ? or second_id = ?"
	rows, err := mysql.DB.Query(sqlStr, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend PO.FriendPO
		err := rows.Scan(&friend)
		if err != nil {
			return nil, err
		}

		friends = append(friends, friend)
	}

	return friends, nil
}

func QueryFriendship(friendShipID int64) (friend PO.FriendPO, err error) {
	sqlStr := "select * from friend where friendship_id = ?"
	err = mysql.DB.QueryRow(sqlStr, friendShipID).Scan(friend)
	if err != nil {
		return friend, err
	}

	return friend, nil
}

func QueryFriendshipBy2ID(firstID int64, secondID int64) (friend PO.FriendPO, err error) {
	sqlStr := "select * from friend where (first_id = ? and second_id = ?) or (first_id = ? and second_id = ?)"
	err = mysql.DB.QueryRow(sqlStr, firstID, secondID, secondID, firstID).Scan(friend)
	if err != nil {
		return friend, err
	}

	return friend, nil
}

func DeleteFriend(friendshipID int64) (err error) {
	sqlStr := "delete from friend where friendship_id = ?"
	ret, err := mysql.DB.Exec(sqlStr, friendshipID)
	if err != nil {
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
		return err
	}

	return nil
}

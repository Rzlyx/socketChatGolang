package friend_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/model/PO"
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

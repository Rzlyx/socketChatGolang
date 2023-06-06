package friendCircle_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/model/PO"
)

func QueryFriendCircle(userID int64) (context []PO.FriendCirclePO, err error) {
	sqlStr := "select * from friend_circle where sender = ?"
	err = mysql.DB.Select(&context, sqlStr, userID)
	if err != nil {
		return context, err
	}

	return context, nil
}

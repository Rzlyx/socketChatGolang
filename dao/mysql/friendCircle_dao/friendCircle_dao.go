package friendCircle_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/model/PO"
	"fmt"
)

func CreateFriendCircle(circle *PO.FriendCirclePO) error {
	strSql := "insert friend_circle (news_id,sender_id,news,type,black_list,white_list,create_time,likes,is_deleted,extra) values (?,?,?,?,?,?,?,?,?,?)"
	result, err := mysql.DB.Exec(strSql, 
		circle.NewsID,
		circle.SenderID,
		circle.News,
		circle.Type,
		circle.BlackList,
		circle.WhiteList,
		circle.CreateTime,
		circle.Likes,
		circle.IsDeleted,
		circle.Extra,
	)
	if err != nil {
		fmt.Println("[WriteGroupMsg], insert err is ", err.Error())
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		fmt.Println("[WriteGroupMsg], insert err is ", err.Error())
		return err
	}
	return nil
}

func UpdateFriendCircle(circle *PO.FriendCirclePO) error {
	strSql := "UPDATE friend_circle SET sender_id = ?, news = ?, type = ?, black_list = ?, white_list = ?, create_time = ?, likes = ?, is_deleted = ?, extra = ? WHERE news_id = ?"
	_, err := mysql.DB.Exec(strSql, 
		circle.SenderID,
		circle.News,
		circle.Type,
		circle.BlackList,
		circle.WhiteList,
		circle.CreateTime,
		circle.Likes,
		circle.IsDeleted,
		circle.Extra,
		circle.NewsID,
	)
	if err != nil {
		fmt.Println("[UpdateFriendCircle], update err is ", err.Error())
		return err
	}
	return nil
}

func MGetFriendCircle(NewsID int64) (*PO.FriendCirclePO, error) {
	var circle PO.FriendCirclePO
	strSql := "SELECT * from friend_circle where news_id = ?"
	err := mysql.DB.Get(&circle, strSql, NewsID)
	if err != nil {
		fmt.Println("[MGetFriendCircle], select err is ", err.Error())
		return nil, err
	}

	return &circle, nil
}


func QueryFriendCircle(userID int64) (context []PO.FriendCirclePO, err error) {
	sqlStr := "select * from friend_circle where sender = ?"
	err = mysql.DB.Select(&context, sqlStr, userID)
	if err != nil {
		return context, err
	}

	return context, nil
}

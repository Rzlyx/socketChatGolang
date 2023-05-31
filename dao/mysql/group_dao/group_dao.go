package group_dao

import (
	"dou_yin/dao/mysql"
	"fmt"
)

func MGetGroupListByUserID(UserID int64) ([]GroupDO, error) {
	var list []GroupDO
	strSql := "select * from group where user_id = ?"
	rows, err := mysql.DB.Query(strSql, UserID)
	if err != nil {
		fmt.Println("[MGetGroupListByUserID] query mysql err is ", err.Error())
		return list, err
	}

	for rows.Next() {
		var group GroupDO
		err := rows.Scan(&group)
		if err != nil {
			fmt.Println("[MGetGroupListByUserID] scan group err is ", err.Error())
			return list, err
		}
		list = append(list, group)
	}
	return list, nil
}

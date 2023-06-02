package group_dao

import (
	"database/sql"
	"dou_yin/dao/mysql"
	"fmt"
)

func MGetGroupInfoByGroupID(GroupID int64) (*GroupInfoPO, error) {
	var group GroupInfoPO
	strSql := "select * from group_info where group_id = ?"
	err := mysql.DB.Get(&group, strSql, GroupID)
	if err != nil {
		fmt.Println("[MGetGroupInfoByGroupID] query or get err is ", err.Error())
		return nil, err
	}
	return &group, nil
}

func CreateGroupInfo(info *GroupInfoPO) (error) {
	strSql := "INSERT group_info (group_id, owner_id, group_name, description, user_ids, admin_ids, slience_list, create_time, is_deleted, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := mysql.DB.Exec(strSql, 
		info.OwnerID,
		info.GroupName,
		info.Description,
		info.UserIds,
		info.AdminIds,
		info.SilenceList,
		info.CreateTime,
		info.IsDeleted,
		info.Extra)
	if err != nil {
		fmt.Println("[CreateGroupInfo] insert group_info err is ", err.Error())
		return err
	}
	return  nil
}

func UpdateGroupInfo(info *GroupInfoPO) error {
	strSql := "UPDATE group_info SET owner_id = ?, group_name = ?, description = ?, user_ids = ?, admin_ids = ?, slience_list = ?, create_time = ?, is_deleted = ?, extra = ? WHERE group_id = ?"
	_, err := mysql.DB.Exec(strSql, 
		info.OwnerID,
		info.GroupName,
		info.Description,
		info.UserIds,
		info.AdminIds,
		info.SilenceList,
		info.CreateTime,
		info.IsDeleted,
		info.Extra,
		info.GroupID)
	if err != nil {
		fmt.Println("[UpdateGroupInfo] update group_info err is ", err.Error())
		return err
	}
	return nil
}

func DeleteGroupInfoByGroupID(GroupID int64) error {
	strSql := "delete from group_info where group_id = ?"
	_, err := mysql.DB.Exec(strSql, GroupID)
	if err != nil {
		fmt.Println("[DeleteGroupInfoByGroupID] delete group_info err is ", err.Error())
		return err
	}
	return nil
}

func MGetGroupListByUserID(UserID int64) ([]GroupPO, error) {
	var list []GroupPO
	strSql := "select * from group_num where user_id = ?"
	err := mysql.DB.Select(&list, strSql, UserID)
	if err != nil {
		fmt.Println("[MGetGroupListByUserID] select mysql err is ", err.Error())
		return list, err
	}

	return list, nil
}

func MGetGroupByUserIDandGroupID(UserId, GroupID int64) (*GroupPO, error) {
	var group GroupPO
	strSql := "select * from group_num where user_id = ? and group_id = ?"
	err := mysql.DB.Get(&group, strSql, UserId, GroupID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("[MGetGroupByUserIDandGroupID] query or get err is ", err.Error())
		return nil, err
	}
	return &group, nil
}

func IsGroupUser(UserID, GroupID int64) (bool, error) {
	_, err := MGetGroupByUserIDandGroupID(UserID, GroupID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteGroupByUserIDandGroupID(UserId, GroupID int64) (bool, error) {
	strSql := "delete from group_num where user_id = ? and group_id = ?"
	_, err := mysql.DB.Exec(strSql, UserId, GroupID)
	if err != nil {
		fmt.Println("[DeleteGroupByUserIDandGroupID] delete group_num err is ", err.Error())
		return false, err
	}
	return true, nil
}

func CreateGroupByGroupPO(group GroupPO) (bool, error) {
	strSql := "INSERT INTO group_num (group_id, group_name, user_id, type, create_time, extra) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := mysql.DB.Exec(strSql, 
		group.GroupID,
		group.GroupName,
		group.UserID,
		group.Type,
		group.CreateTime,
		group.Extra)
	if err != nil {
		fmt.Println("[CreateGroupByGroupPO] insert err is ", err.Error())
		return false, err
	}
	return true, nil
}

func UpdateGroupByGroupPO(group GroupPO) (bool, error) {
	strSql := "UPDATE group_num SET group_id = ?, group_name = ?, user_id = ?, type = ?, create_time = ?, extra = ? WHERE user_id = ? and group_id = ?"
	_, err := mysql.DB.Exec(strSql, 
		group.GroupID,
		group.GroupName,
		group.UserID,
		group.Type,
		group.CreateTime,
		group.Extra,
		group.UserID,
		group.GroupID)
	if err != nil {
		fmt.Println("[UpdateGroupByGroupPO], update err is ", err.Error())
		return false, err
	}
	return true, nil
} 
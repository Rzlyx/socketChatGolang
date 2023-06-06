package apply_dao

import (
	"database/sql"
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"dou_yin/service/DO"
	"fmt"
)

func Insert(application DO.AddFriendApplication) (err error) {
	sqlStr := "insert into apply(apply_id, applicant, target_id, type, reason, create_time) values (?, ?, ?, ?, ?, ?)"
	_, err = mysql.DB.Exec(sqlStr, application.ApplyID, application.ApplicantID, application.FriendID, application.Type, application.Reason, application.CreateTime)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func QueryApplication(userID int64, friendID int64) (applyPO PO.ApplyPO, err error) {
	sqlStr := "select * from apply where applicant = ? and target_id = ?"
	err = mysql.DB.Get(&applyPO, sqlStr, userID, friendID)
	if err != nil {
		logger.Log.Error(err.Error())
		return applyPO, err
	}

	return applyPO, nil
}

func QueryFriendApply(userID int64) (applys []PO.ApplyPO, err error) {
	sqlStr := "select * from apply where target_id = ? and type = 1"
	rows, err := mysql.DB.Query(sqlStr, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var apply PO.ApplyPO
		err := rows.Scan(&apply.ApplyID, &apply.Applicant, &apply.Type, &apply.Status, &apply.Reason, &apply.CreateTime, &apply.Extra, &apply.TargetID)
		if err != nil {
			logger.Log.Error(err.Error())
			return nil, err
		}

		applys = append(applys, apply)
	}

	return applys, nil
}

func UpdateStatus(applyID int64, status int) (err error) {
	sqlStr := "update apply set status = ? where apply_id = ?"
	_, err = mysql.DB.Exec(sqlStr, status, applyID)
	if err != nil {
		return err
	}

	return nil
}

func Delete(applyID int64) (err error) {
	sqlStr := "delete from apply where apply_id = ?"
	_, err = mysql.DB.Exec(sqlStr, applyID)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func CreateApplication(apply *PO.ApplyPO) (error) {
	sqlStr := "INSERT INTO apply (apply_id, applicant, target_id, type, status, reason, create_time, extra) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := mysql.DB.Exec(sqlStr, 
		apply.ApplyID,
		apply.Applicant,
		apply.TargetID,
		apply.Type,
		apply.Status,
		apply.Reason,
		apply.CreateTime,
		apply.Extra)
	if err != nil {
		fmt.Println("[CreateApplication], insert err is ", err.Error())
		return err
	}
	return nil
}

func MGetApplicationListByGroupID(GroupID int64) (*[]PO.ApplyPO, error) {
	var list []PO.ApplyPO
	sqlStr := "select * from apply where target_id = ? and type = 0"
	err := mysql.DB.Select(&list, sqlStr, GroupID)
	if err != nil {
		fmt.Println("[MGetApplicationListByGroupID], query select err is ", err.Error())
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return &list, nil
}

// 获取被邀请记录
func MGetApplicationListByUserID(UserID int64) (*[]PO.ApplyPO, error) {
	var list []PO.ApplyPO
	sqlStr := "select * from apply where applicant = ? and type = 2"
	err := mysql.DB.Select(&list, sqlStr, UserID)
	if err != nil {
		fmt.Println("[MGetApplicationListByUserID], query select err is ", err.Error())
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return &list, nil
}

func MGetApplicationByGroupIDandUserID(GroupID, UserID int64) (*PO.ApplyPO, error) {
	var application PO.ApplyPO
	sqlStr := "select * from apply where target_id = ? and applicant = ? and type = 0"
	err := mysql.DB.Get(&application, sqlStr, GroupID, UserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		fmt.Println("[MGetApplicationByGroupID], select err is ", err.Error())
		return nil, err
	}
	return &application, nil
}

func DeleteApplicationByApplyID(applyID int64) (error){
	sqlStr := "delete from apply where apply_id = ?"
	_, err := mysql.DB.Exec(sqlStr, applyID)
	if err != nil {
		fmt.Println("[DeleteApplicationByApplyID], delete err id ", err.Error())
		return err
	}

	return nil
}

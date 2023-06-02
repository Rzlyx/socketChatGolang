package apply_dao

import (
	"dou_yin/dao/mysql"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"dou_yin/service/DO"
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

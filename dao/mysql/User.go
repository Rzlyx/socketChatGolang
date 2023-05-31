package mysql

import (
	"dou_yin/model/PO"
	"dou_yin/pkg/utils"
	"fmt"
)

func Register(p *PO.User) error {
	tableName := utils.ShiftToStringFromInt64(p.UserID)
	sqlStr := "insert into user (user_id, username, password, email) VALUES (?,?,?,?)"
	sqlStr1 := "CREATE TABLE IF NOT EXISTS `" + tableName + "` (" +
		"user_id BIGINT PRIMARY KEY UNIQUE NOT NULL, " +
		"username VARCHAR(20) NOT NULL, " +
		"remark VARCHAR(20)," +
		"time VARCHAR(20)," +
		"message VARCHAR(20)" +
		")"

	// 开始事务
	tx, err := DB.Beginx()
	if err != nil {
		// 处理事务开始错误
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() // 在出现 panic 时回滚事务
			panic(p)      // 继续向上抛出 panic
		} else if err != nil {
			tx.Rollback() // 在出现其他错误时回滚事务
		} else {
			tx.Commit() // 提交事务
		}
	}()

	// 执行第一个 SQL 语句
	_, err = tx.Exec(sqlStr, p.UserID, p.UserName, p.Password, p.EMail)
	if err != nil {
		// 处理第一个 SQL 语句执行错误
		return err
	}
	// 执行第二个 SQL 语句
	_, err = tx.Exec(sqlStr1)
	if err != nil {
		// 处理第二个 SQL 语句执行错误
		return err
	}
	// 提交事务
	err = tx.Commit()
	if err != nil {
		// 处理事务提交错误
		return err
	}
	return err
}

func Login(username string) (p1 *PO.UserPO, err error) {
	p1 = new(PO.UserPO)
	sqlStr := "select * from user where username = ?"
	err = DB.Get(p1, sqlStr, username)
	return p1, err
}

func GetContactorList(Id string) (p *PO.ContactorList, err error) {
	p = new(PO.ContactorList)
	sqlStr := "select * from `" + Id + "`"
	err = DB.Select(&p.ContactorList, sqlStr)
	fmt.Println(err)
	return p, err
}

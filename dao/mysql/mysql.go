package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var DB *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"))
	// 也可以使用MustConnect连接不成功就panic
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	DB.SetMaxOpenConns(viper.GetInt("mysql,max_open_conns"))
	DB.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}

func Close() {
	_ = DB.Close()
}

// TxFunc 定义了在事务中执行的函数类型
type TxFunc func(*sql.Tx) error

// Tx 函数用于在事务中执行给定的函数
func Tx(db *sqlx.DB, fn TxFunc) error {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 执行给定的函数
	err = fn(tx)
	if err != nil {
		// 发生错误，回滚事务
		_ = tx.Rollback()
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		// 提交事务时发生错误，回滚事务
		_ = tx.Rollback()
		return err
	}

	return nil
}


// TODO:Tx
// err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
	
// 	return nil
// })
// if err != nil {
// 	return err
// }
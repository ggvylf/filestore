package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// 声明sql连接
var db *sql.DB

// 初始化连接
func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fileserver?charset=utf8mb4")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect mysql,err=", err.Error())
		os.Exit(1)
	}
}

// 调用的时候返回已经声明过的连接 而不是新生成一个 实现连接复用
func DBConn() *sql.DB {
	return db
}

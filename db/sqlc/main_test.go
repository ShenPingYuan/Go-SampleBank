// Go测试入口，使用的mysql数据库
package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// 定义全局变量
var testQueries *Queries
var testDB *sql.DB
var testContext context.Context

const (
	dbDriver = "mysql"
	dbSource = "root:1230@tcp(localhost:3306)/simple_bank?parseTime=true"
)

// 初始化数据库连接
func TestMain(m *testing.M) {
	var err error
	//连接数据库
	testDB, err = sql.Open(dbDriver, dbSource)

	//创建上下文
	testContext = context.Background()

	//判断是否连接成功
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	//初始化Queries
	testQueries = New(testDB)
	//执行测试
	os.Exit(m.Run())
}

// Go测试入口，使用的mysql数据库
package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ShenPingYuan/go-webdemo/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// 定义全局变量
var testQueries *Queries
var testDB *sql.DB
var testContext context.Context

// 初始化数据库连接
func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	//连接数据库
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	//判断是否连接成功
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	//创建上下文
	testContext = context.Background()

	//准备迁移
	driver, err := mysql.WithInstance(testDB, &mysql.Config{})
	if err != nil {
		log.Fatal("cannot create migration driver:", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://../migration",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("cannot create migration instance:", err)
	}

	//执行迁移
	//migration.Down()
	migration.Up()
	//migration.Version()

	//初始化Queries
	testQueries = New(testDB)

	//插入一行user
	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: util.RandomPassword(),
		Email:          util.RandomEmail(),
		FullName:       util.RandomUsername(),
	}
	_, err = testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		log.Fatal("cannot create user:", err)
	}
	//执行测试
	os.Exit(m.Run())
}

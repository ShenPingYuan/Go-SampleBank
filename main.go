package main

import (
	"database/sql"

	"github.com/ShenPingYuan/go-webdemo/api"
	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// 服务地址
	serverAddress = ":8080"
	// 数据库地址
	dbDriver = "mysql"
	dbSource = "root:1230@tcp(localhost:3306)/simple_bank?parseTime=true&multiStatements=true"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		panic(err)
	}
}

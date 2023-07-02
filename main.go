package main

import (
	"database/sql"

	"github.com/ShenPingYuan/go-webdemo/api"
	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/ShenPingYuan/go-webdemo/util"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		panic(err)
	}
}

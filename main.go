package main

import (
	"database/sql"
	"log"

	"github.com/ShenPingYuan/go-webdemo/api"
	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/ShenPingYuan/go-webdemo/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	server := api.NewServer(config, store)

	//准备迁移
	driver, err := mysql.WithInstance(conn, &mysql.Config{})
	if err != nil {
		log.Fatal("cannot create migration driver:", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://./db/migration/",
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

	err = server.Start(config.ServerAddress)
	if err != nil {
		panic(err)
	}
}

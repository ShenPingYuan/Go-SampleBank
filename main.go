package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/ShenPingYuan/go-webdemo/api"
	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/ShenPingYuan/go-webdemo/grpc_api"
	pb "github.com/ShenPingYuan/go-webdemo/protobuffer"
	"github.com/ShenPingYuan/go-webdemo/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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

	runGrpcServer(config, store, conn)
}

func runGrpcServer(config util.Config, store db.Store, conn *sql.DB) {

	server := grpc_api.NewServer(config, store)

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start grpc server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}

func runHttpServer(config util.Config, store db.Store, conn *sql.DB) {
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

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		panic(err)
	}
}

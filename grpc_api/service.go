package grpc_api

import (
	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	pb "github.com/ShenPingYuan/go-webdemo/protobuffer"
	"github.com/ShenPingYuan/go-webdemo/token"
	"github.com/ShenPingYuan/go-webdemo/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) *Server {

	tokenMaker, err := token.NewJWTMaker(config.SymmetricKey)
	if err != nil {
		panic(err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server
}

func (server *Server) Start(address string) error {
	return nil
}

package api

import (
	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/ShenPingYuan/go-webdemo/token"
	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
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

	router := gin.Default()

	binding.Validator.Engine().(*validator.Validate).RegisterValidation("currency", validCurrency)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.getPagedAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfer", server.createTransfer)

	router.POST("/users", server.createUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

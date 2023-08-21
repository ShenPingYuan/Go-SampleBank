package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.getPagedAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfer", server.createTransfer)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.login)

	server.router = router
}

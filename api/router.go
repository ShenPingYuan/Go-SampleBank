package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()

	authRouters := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouters.POST("/accounts", server.createAccount)
	authRouters.GET("/accounts", server.getPagedAccounts)
	authRouters.GET("/accounts/:id", server.getAccount)
	authRouters.DELETE("/accounts/:id", server.deleteAccount)

	authRouters.POST("/transfer", server.createTransfer)

	authRouters.GET("/users/userinfo", server.getUserInfo)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.login)

	router.POST("/token/renew", server.renewAccessToken)

	server.router = router
}

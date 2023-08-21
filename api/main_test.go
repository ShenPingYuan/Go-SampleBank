package api

import (
	"os"
	"testing"
	"time"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		DBDriver:            "mysql",
		DBSource:            "root:1230@tcp(localhost:3306)/simple_bank?parseTime=true&multiStatements=true",
		ServerAddress:       ":8080",
		SymmetricKey:        util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server := NewServer(config, store)

	return server
}

// main
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

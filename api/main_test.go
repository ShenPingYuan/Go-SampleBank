package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// main
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

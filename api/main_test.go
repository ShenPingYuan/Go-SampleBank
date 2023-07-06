package api

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// main
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

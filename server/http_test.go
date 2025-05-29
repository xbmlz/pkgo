package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPServer(t *testing.T) {
	r := gin.Default()
	r.Handle("GET", "/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	server := NewHTTPServer(r, ":8080")

	go func() {
		time.Sleep(2 * time.Second)

		resp, err := http.Get("http://localhost:8080")
		assert.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		time.Sleep(1 * time.Second)
		server.Shutdown()
		t.Log("shutdown completed")
	}()

	server.Run()
}

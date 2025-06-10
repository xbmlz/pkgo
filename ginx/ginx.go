package ginx

import (
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/xbmlz/pkgo/log"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(ginzap.Ginzap(log.GetLogger(), time.DateTime, true))
	r.Use(ginzap.CustomRecoveryWithZap(log.GetLogger(), true, func(c *gin.Context, err any) {
		ResponseCustom(c, http.StatusInternalServerError, "Internal Server Error", nil)
	}))
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.NoRoute(func(ctx *gin.Context) {
		ResponseCustom(ctx, http.StatusNotFound, "Not Found", nil)
	})
	return r
}

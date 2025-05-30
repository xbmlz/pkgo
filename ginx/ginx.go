package ginx

import (
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/xbmlz/pkgo/logx"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(logx.GetLogger(), time.DateTime, true))
	r.Use(ginzap.CustomRecoveryWithZap(logx.GetLogger(), true, func(c *gin.Context, err any) {
		ResponseCustom(c, http.StatusInternalServerError, "Internal Server Error", nil)
	}))
	r.Use(static.ServeRoot("/static", "./static"))
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.NoRoute(func(ctx *gin.Context) {
		ResponseCustom(ctx, http.StatusNotFound, "Not Found", nil)
	})
	return r
}

package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetQuery[T any](c *gin.Context, key string) T {
	return convert[T](c.Query(key))
}

func GetParam[T any](c *gin.Context, key string) T {
	return convert[T](c.Param(key))
}

func GetForm[T any](c *gin.Context, key string) T {
	return convert[T](c.PostForm(key))
}

func convert[T any](value string) T {
	var target T

	switch any(target).(type) {
	case string:
		return any(value).(T)
	case int:
		return any(cast.ToInt(value)).(T)
	case int64:
		return any(cast.ToInt64(value)).(T)
	case float32:
		return any(cast.ToFloat32(value)).(T)
	case float64:
		return any(cast.ToFloat64(value)).(T)
	case bool:
		return any(cast.ToBool(value)).(T)
	default:
		return any(value).(T)
	}
}

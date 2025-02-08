package middleware

import (
	"gotoraft/pkg/errors"
	"gotoraft/pkg/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 中间件处理panic
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				logger.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)

				// 包装错误
				var apiErr *errors.Error
				switch e := err.(type) {
				case *errors.Error:
					apiErr = e
				case error:
					apiErr = errors.WrapError(e, errors.ErrInternalServer, "服务器内部错误")
				default:
					apiErr = errors.NewError(errors.ErrInternalServer, "服务器内部错误")
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    apiErr.Code,
					"message": apiErr.Message,
				})
			}
		}()
		c.Next()
	}
}

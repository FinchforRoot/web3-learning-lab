package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 日志记录中间件
func LoggerMiddleware() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return gin.LoggerWithFormatter(
		func(param gin.LogFormatterParams) string {
			logger.WithFields(
				logrus.Fields{
					"status_code": param.StatusCode,
					"latency":     param.Latency,
					"client_ip":   param.ClientIP,
					"method":      param.Method,
					"path":        param.Path,
					"user_agent":  param.Request.UserAgent(),
					"error":       param.ErrorMessage,
					"timestamp":   param.TimeStamp.Format(time.RFC3339),
				}).Info("HTTP Request")
			return ""
		})
}

// 全局异常处理
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error":  err,
					"path":   ctx.Request.URL.Path,
					"method": ctx.Request.Method,
				}).Error("Panic recovered")
				ctx.JSON(500, gin.H{
					"code":    500,
					"message": "Internal server error",
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

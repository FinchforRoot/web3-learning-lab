package middleware

import (
	"my-blog-project/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 请求头获取token
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			utils.Unauthorized(ctx, "未获取到Authorization")
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.Unauthorized(ctx, "未获取到Bearer")
			ctx.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Unauthorized(ctx, "Invalid token")
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}

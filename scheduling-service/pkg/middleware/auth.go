package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/web"
	"net/http"
	"os"
)

func Authentication() gin.HandlerFunc {
	requiredToken := os.Getenv("TOKEN")

	return func(ctx *gin.Context) {
		tokenReceived := ctx.GetHeader("SECRET_TOKEN")

		if tokenReceived == "" {
			web.BadResponse(ctx, http.StatusUnauthorized, "error", "Token not found")
			ctx.Abort()
			return
		}

		if tokenReceived != requiredToken {
			web.BadResponse(ctx, http.StatusUnauthorized, "error", "Invalid token provided")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

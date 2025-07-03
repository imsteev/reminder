package handler

import (
	"net/http"
	"reminder-app/controller/protocol"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
)

func httpOptionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func clerkAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, protocol.ErrorResponse{Error: "authorization header required"})
			c.Abort()
			return
		}

		// Wrap the request/response for clerk validation
		wrappedHandler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, protocol.ErrorResponse{Error: "invalid or expired token"})
				c.Abort()
				return
			}

			c.Set("userID", claims.Subject)
		}))

		wrappedHandler.ServeHTTP(c.Writer, c.Request)

		// Only continue if auth was successful (userID was set)
		if _, exists := c.Get("userID"); exists {
			c.Next()
		}
	}
}

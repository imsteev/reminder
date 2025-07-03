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

		// AI generated.... it works but I don't fully understand it. Come back to this to see if we can simplify it.
		wrappedHandler := clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, protocol.ErrorResponse{Error: "invalid or expired token"})
				c.Abort()
				return
			}
			c.Set("clerkID", claims.Subject)
		}))

		wrappedHandler.ServeHTTP(c.Writer, c.Request)

		if _, exists := c.Get("clerkID"); exists {
			c.Next()
		}
	}
}

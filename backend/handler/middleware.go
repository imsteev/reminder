package handler

import (
	"net/http"
	"reminder-app/controller/protocol"
	"reminder-app/lib/actor"
	"reminder-app/models"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func injectActorMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		clerkID, exists := c.Get("clerkID")
		if !exists {
			c.JSON(http.StatusUnauthorized, protocol.ErrorResponse{Error: "unauthorized"})
			c.Abort()
			return
		}

		user := &models.User{}
		if err := db.First(user, "clerk_id = ?", clerkID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, protocol.ErrorResponse{Error: "unauthorized"})
			c.Abort()
			return
		}

		c.Set("user", actor.New(user.ID, user.ClerkID))
		c.Next()
	}
}

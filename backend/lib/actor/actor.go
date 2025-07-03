package actor

import "github.com/gin-gonic/gin"

type Actor struct {
	UserID  uint
	ClerkID string
}

func New(id uint, clerkID string) *Actor {
	return &Actor{UserID: id, ClerkID: clerkID}
}

func (a *Actor) GetUserIDInt64() int64 {
	return int64(a.UserID)
}

func (a *Actor) GetClerkID() string {
	return a.ClerkID
}

func FromGin(c *gin.Context) *Actor {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}
	return user.(*Actor)
}

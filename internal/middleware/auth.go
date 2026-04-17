package middleware

import (
	"net/http"

	"github.com/abhaybhargav/go-expense-boilerplate/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		if uid == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("user_id", uid)
		if role := session.Get("role"); role != nil {
			c.Set("role", role)
		}
		c.Next()
	}
}

func RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := c.Get("role")
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		current := models.Role(raw.(string))
		for _, r := range roles {
			if r == current {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

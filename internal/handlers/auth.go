package handlers

import (
	"context"
	"net/http"

	"github.com/abhaybhargav/go-expense-boilerplate/internal/database"
	"github.com/abhaybhargav/go-expense-boilerplate/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"Title": "Login"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var (
		id    int64
		hash  string
		role  string
		name  string
		active bool
	)
	err := h.DB.Pool.QueryRow(context.Background(),
		`SELECT id, password_hash, role, full_name, active FROM users WHERE email = $1`, email).
		Scan(&id, &hash, &role, &name, &active)
	if err != nil || !active {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", id)
	session.Set("role", role)
	session.Set("full_name", name)
	_ = session.Save()

	c.Redirect(http.StatusFound, redirectFor(models.Role(role)))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func redirectFor(r models.Role) string {
	switch r {
	case models.RoleSuperAdmin:
		return "/admin"
	case models.RoleProjectManager:
		return "/approvals"
	default:
		return "/dashboard"
	}
}

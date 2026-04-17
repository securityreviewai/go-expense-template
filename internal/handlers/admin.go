package handlers

import (
	"context"
	"net/http"

	"github.com/abhaybhargav/go-expense-boilerplate/internal/database"
	"github.com/abhaybhargav/go-expense-boilerplate/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	DB *database.DB
}

func NewAdminHandler(db *database.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) Index(c *gin.Context) {
	rows, err := h.DB.Pool.Query(context.Background(),
		`SELECT id, email, full_name, role, active, created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin.html", gin.H{"Error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.Active, &u.CreatedAt); err == nil {
			users = append(users, u)
		}
	}
	c.HTML(http.StatusOK, "admin.html", gin.H{"Users": users, "Title": "Admin"})
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	email := c.PostForm("email")
	name := c.PostForm("full_name")
	role := c.PostForm("role")
	password := c.PostForm("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = h.DB.Pool.Exec(context.Background(),
		`INSERT INTO users (email, password_hash, full_name, role, active) VALUES ($1,$2,$3,$4,true)`,
		email, string(hash), name, role)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/admin")
}

package router

import (
	"html/template"
	"net/http"

	"github.com/abhaybhargav/go-expense-boilerplate/internal/config"
	"github.com/abhaybhargav/go-expense-boilerplate/internal/database"
	"github.com/abhaybhargav/go-expense-boilerplate/internal/handlers"
	"github.com/abhaybhargav/go-expense-boilerplate/internal/middleware"
	"github.com/abhaybhargav/go-expense-boilerplate/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config, db *database.DB) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	store := cookie.NewStore([]byte(cfg.SessionKey))
	store.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 60 * 60 * 24 * 7, Secure: cfg.Environment == "production"})
	r.Use(sessions.Sessions("expense_session", store))

	r.SetFuncMap(template.FuncMap{
		"dollars": func(cents int64) string {
			return formatCents(cents)
		},
	})
	r.LoadHTMLGlob("web/templates/**/*")
	r.Static("/static", "web/static")

	auth := handlers.NewAuthHandler(db)
	expense := handlers.NewExpenseHandler(db)
	admin := handlers.NewAdminHandler(db)

	r.GET("/healthz", handlers.Health(db))
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })
	r.GET("/login", auth.ShowLogin)
	r.POST("/login", auth.Login)
	r.POST("/logout", auth.Logout)

	authed := r.Group("/", middleware.RequireAuth())
	{
		authed.GET("/dashboard", expense.Dashboard)
		authed.POST("/reports", expense.CreateReport)
		authed.POST("/reports/:id/submit", expense.SubmitReport)
	}

	managers := r.Group("/", middleware.RequireAuth(), middleware.RequireRole(models.RoleProjectManager, models.RoleSuperAdmin))
	{
		managers.GET("/approvals", expense.Approvals)
		managers.POST("/approvals/:id/decide", expense.Decide)
	}

	adminG := r.Group("/admin", middleware.RequireAuth(), middleware.RequireRole(models.RoleSuperAdmin))
	{
		adminG.GET("", admin.Index)
		adminG.POST("/users", admin.CreateUser)
	}

	return r
}

func formatCents(c int64) string {
	neg := c < 0
	if neg {
		c = -c
	}
	dollars := c / 100
	cents := c % 100
	sign := ""
	if neg {
		sign = "-"
	}
	return sign + formatInt(dollars) + "." + twoDigit(cents)
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func twoDigit(n int64) string {
	if n < 10 {
		return "0" + formatInt(n)
	}
	return formatInt(n)
}

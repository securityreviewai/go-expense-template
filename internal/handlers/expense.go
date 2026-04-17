package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/abhaybhargav/go-expense-boilerplate/internal/database"
	"github.com/abhaybhargav/go-expense-boilerplate/internal/models"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	DB *database.DB
}

func NewExpenseHandler(db *database.DB) *ExpenseHandler {
	return &ExpenseHandler{DB: db}
}

func (h *ExpenseHandler) Dashboard(c *gin.Context) {
	uid := c.GetInt64("user_id")
	rows, err := h.DB.Pool.Query(context.Background(),
		`SELECT id, title, status, created_at FROM expense_reports WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`, uid)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "dashboard.html", gin.H{"Error": err.Error()})
		return
	}
	defer rows.Close()

	var reports []models.ExpenseReport
	for rows.Next() {
		var r models.ExpenseReport
		if err := rows.Scan(&r.ID, &r.Title, &r.Status, &r.CreatedAt); err == nil {
			reports = append(reports, r)
		}
	}
	c.HTML(http.StatusOK, "dashboard.html", gin.H{"Reports": reports, "Title": "My Reports"})
}

func (h *ExpenseHandler) CreateReport(c *gin.Context) {
	uid := c.GetInt64("user_id")
	title := c.PostForm("title")
	if title == "" {
		c.String(http.StatusBadRequest, "title is required")
		return
	}
	_, err := h.DB.Pool.Exec(context.Background(),
		`INSERT INTO expense_reports (user_id, title, status) VALUES ($1, $2, 'draft')`, uid, title)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/dashboard")
}

func (h *ExpenseHandler) SubmitReport(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	uid := c.GetInt64("user_id")
	_, err := h.DB.Pool.Exec(context.Background(),
		`UPDATE expense_reports SET status='pending', submitted_at=NOW() WHERE id=$1 AND user_id=$2 AND status='draft'`, id, uid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/dashboard")
}

func (h *ExpenseHandler) Approvals(c *gin.Context) {
	rows, err := h.DB.Pool.Query(context.Background(),
		`SELECT r.id, r.title, r.status, r.created_at, u.full_name
		 FROM expense_reports r JOIN users u ON u.id = r.user_id
		 WHERE r.status='pending' ORDER BY r.submitted_at ASC`)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "approvals.html", gin.H{"Error": err.Error()})
		return
	}
	defer rows.Close()

	type pending struct {
		models.ExpenseReport
		Submitter string
	}
	var list []pending
	for rows.Next() {
		var p pending
		if err := rows.Scan(&p.ID, &p.Title, &p.Status, &p.CreatedAt, &p.Submitter); err == nil {
			list = append(list, p)
		}
	}
	c.HTML(http.StatusOK, "approvals.html", gin.H{"Pending": list, "Title": "Pending Approvals"})
}

func (h *ExpenseHandler) Decide(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	approverID := c.GetInt64("user_id")
	action := c.PostForm("action")

	var status models.ExpenseStatus
	switch action {
	case "approve":
		status = models.StatusApproved
	case "reject":
		status = models.StatusRejected
	default:
		c.String(http.StatusBadRequest, "invalid action")
		return
	}
	_, err := h.DB.Pool.Exec(context.Background(),
		`UPDATE expense_reports SET status=$1, approved_by_id=$2, approved_at=NOW(), rejected_reason=$3
		 WHERE id=$4 AND status='pending'`,
		status, approverID, c.PostForm("reason"), id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/approvals")
}

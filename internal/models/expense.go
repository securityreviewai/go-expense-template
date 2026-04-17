package models

import "time"

type ExpenseStatus string

const (
	StatusDraft    ExpenseStatus = "draft"
	StatusPending  ExpenseStatus = "pending"
	StatusApproved ExpenseStatus = "approved"
	StatusRejected ExpenseStatus = "rejected"
)

type Expense struct {
	ID          int64         `json:"id"`
	ReportID    int64         `json:"report_id"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	AmountCents int64         `json:"amount_cents"`
	Currency    string        `json:"currency"`
	IncurredOn  time.Time     `json:"incurred_on"`
	ReceiptURL  string        `json:"receipt_url,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type ExpenseReport struct {
	ID             int64         `json:"id"`
	UserID         int64         `json:"user_id"`
	Title          string        `json:"title"`
	Status         ExpenseStatus `json:"status"`
	SubmittedAt    *time.Time    `json:"submitted_at,omitempty"`
	ApprovedByID   *int64        `json:"approved_by_id,omitempty"`
	ApprovedAt     *time.Time    `json:"approved_at,omitempty"`
	RejectedReason string        `json:"rejected_reason,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

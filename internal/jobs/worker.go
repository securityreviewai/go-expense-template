package jobs

import (
	"context"
	"log"
	"time"

	"github.com/abhaybhargav/go-expense-boilerplate/internal/database"
)

type Worker struct {
	DB       *database.DB
	Interval time.Duration
}

func NewWorker(db *database.DB) *Worker {
	return &Worker{DB: db, Interval: 30 * time.Second}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	log.Println("background worker started")
	for {
		select {
		case <-ctx.Done():
			log.Println("background worker stopping")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	if err := w.notifyStalePending(ctx); err != nil {
		log.Printf("worker notifyStalePending: %v", err)
	}
}

func (w *Worker) notifyStalePending(ctx context.Context) error {
	_, err := w.DB.Pool.Exec(ctx,
		`INSERT INTO notifications (user_id, kind, message, created_at)
		 SELECT r.user_id, 'stale_pending', 'Your expense report is awaiting approval', NOW()
		 FROM expense_reports r
		 WHERE r.status = 'pending'
		   AND r.submitted_at < NOW() - INTERVAL '3 days'
		   AND NOT EXISTS (
		     SELECT 1 FROM notifications n
		     WHERE n.user_id = r.user_id AND n.kind = 'stale_pending'
		       AND n.created_at > NOW() - INTERVAL '1 day'
		   )`)
	return err
}

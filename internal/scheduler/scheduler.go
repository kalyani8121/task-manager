package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kalyani8121/task-manager/internal/email"
)

// Scheduler runs background jobs automatically
type Scheduler struct {
	db     *sqlx.DB
	mailer *email.EmailSender
}

// NewScheduler creates a new scheduler
func NewScheduler(db *sqlx.DB, mailer *email.EmailSender) *Scheduler {
	return &Scheduler{
		db:     db,
		mailer: mailer,
	}
}

// Start begins the scheduler in background
// Runs as a goroutine — doesn't block the API!
func (s *Scheduler) Start() {
	go s.run()
	log.Println("Email scheduler started")
}

// run is the main scheduler loop
func (s *Scheduler) run() {
	for {
		now := time.Now()

		// Calculate time until next 9AM
		next9AM := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			9, 0, 0, 0,
			now.Location(),
		)

		// If 9AM already passed today → schedule for tomorrow
		if now.After(next9AM) {
			next9AM = next9AM.Add(24 * time.Hour)
		}

		// Wait until 9AM
		waitTime := next9AM.Sub(now)
		log.Printf("Next email scheduled at: %s (in %s)",
			next9AM.Format("2006-01-02 15:04:05"),
			waitTime.Round(time.Minute),
		)

		time.Sleep(waitTime)

		// Send emails to all users
		s.sendDailyEmails()
	}
}

// UserWithEmail holds user info for sending emails
type UserWithEmail struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

// PriorityTaskDB holds task from database
type PriorityTaskDB struct {
	Title        string  `db:"title"`
	Urgency      string  `db:"urgency"`
	DueDate      *string `db:"due_date"`
	DaysUntilDue *int    `db:"days_until_due"`
}

// sendDailyEmails sends priority email to every user
func (s *Scheduler) sendDailyEmails() {
	log.Println("Sending daily priority emails...")

	// Get all users
	var users []UserWithEmail
	err := s.db.Select(&users,
		`SELECT id, name, email FROM users`)
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		return
	}

	// Send email to each user
	for _, user := range users {
		err := s.sendEmailToUser(user)
		if err != nil {
			log.Printf("Failed to send email to %s: %v",
				user.Email, err)
		} else {
			log.Printf("Email sent to %s", user.Email)
		}
	}
}

// sendEmailToUser sends priority email to one user
func (s *Scheduler) sendEmailToUser(user UserWithEmail) error {
	// Get user's pending/in_progress tasks
	var tasks []PriorityTaskDB
	query := `
		SELECT
			title,
			CASE
				WHEN due_date IS NULL THEN NULL
				ELSE TO_CHAR(due_date, 'YYYY-MM-DD')
			END as due_date,
			CASE
				WHEN due_date IS NULL THEN NULL
				ELSE EXTRACT(DAY FROM due_date - NOW())::int
			END as days_until_due,
			CASE
				WHEN due_date IS NULL THEN 'NO DEADLINE'
				WHEN due_date < NOW() THEN 'OVERDUE'
				WHEN EXTRACT(DAY FROM due_date - NOW()) = 0
					THEN 'DUE TODAY'
				WHEN EXTRACT(DAY FROM due_date - NOW()) = 1
					THEN 'DUE TOMORROW'
				WHEN EXTRACT(DAY FROM due_date - NOW()) <= 3
					THEN 'URGENT'
				WHEN EXTRACT(DAY FROM due_date - NOW()) <= 7
					THEN 'UPCOMING'
				ELSE 'NORMAL'
			END as urgency
		FROM tasks
		WHERE user_id = $1
		AND status NOT IN ('completed', 'cancelled')
		ORDER BY
			CASE
				WHEN due_date IS NULL THEN 2
				ELSE 1
			END,
			due_date ASC
	`

	err := s.db.Select(&tasks, query, user.ID)
	if err != nil {
		return fmt.Errorf("get tasks: %w", err)
	}

	// Skip if no pending tasks
	if len(tasks) == 0 {
		log.Printf("No pending tasks for %s — skipping",
			user.Email)
		return nil
	}

	// Convert to email format
	emailTasks := make([]email.PriorityTaskEmail, len(tasks))
	for i, t := range tasks {
		dueDate := ""
		if t.DueDate != nil {
			dueDate = *t.DueDate
		}
		emailTasks[i] = email.PriorityTaskEmail{
			Rank:    i + 1,
			Title:   t.Title,
			Urgency: t.Urgency,
			DueDate: dueDate,
		}
	}

	// Send the email!
	return s.mailer.SendPriorityQueueEmail(
		user.Email,
		user.Name,
		emailTasks,
	)
}

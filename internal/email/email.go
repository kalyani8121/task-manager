package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailSender handles sending emails
type EmailSender struct {
	host    string
	port    string
	email   string
	password string
}

// NewEmailSender creates a new email sender
func NewEmailSender(host,port,email,password string) *EmailSender {
	return &EmailSender{
		host: host,
		port: port,
		email: email,
		password: password,
	}
}

// PriorityTaskEmail holds task info for email
type PriorityTaskEmail struct {
	Rank     int
	Title   string
	DueDate string
	Urgency string
	
}

// SendPriorityQueueEmail sends daily priority email to user
func (s *EmailSender) SendPriorityQueueEmail (
	toEmail   string,
	userName  string,
	tasks     []PriorityTaskEmail,
) error{
	// Build email subject
	subject := " Your Daily Task Priority Queue "

	// Count overdue tasks
	overdueCount := 0
	for _, t := range tasks{
		if strings.Contains(t.Urgency, "OVERDUE") {
			overdueCount++
		}
	}

	if overdueCount > 0 {
		subject = fmt.Sprintf(
			"URGENT: %d Overdue Tasks Need Attention!",overdueCount,
		)
	}

	// Build email body
	body := buildEmailBody(userName, tasks)

	// Send email
	return s.sendEmail(toEmail, subject, body)

}

// build EmailBody creates the email content
func buildEmailBody (userName string, tasks []PriorityTaskEmail) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Hi %s!\n\n", userName))
	sb.WriteString("Here is your daily priority queue:\n")
	sb.WriteString("==========================\n\n")

	if len(tasks) == 0 {
		sb.WriteString("Great job! No pending tasks!\n")
		return sb.String()
	}

	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf(
			"RANK %d: %s\n",
			task.Rank,
			task.Title,
		))
		sb.WriteString(fmt.Sprintf(
			"Status: %s\n",
			task.Urgency,
		))
		if task.DueDate != ""{
			sb.WriteString(fmt.Sprintf(
				"Due: %s\n",
				task.DueDate,
			))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("====================\n")
	sb.WriteString("Complete RANK 1 first!\n\n")
	sb.WriteString("- Task Manager System\n")

	return sb.String()
}

// Send Email sends the actual email via Gmail SMTP

func (s *EmailSender) sendEmail(
	to string,
	subject string,
	body string,
) error {
	// Gmail authentication
	auth := smtp.PlainAuth(
		"",
		s.email,
		s.password,
		s.host,
	)

	// Build email message
	message := fmt.Sprintf(
		"From: Task Manager <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		s.email,
		to,
		subject,
		body,
	)

	//Send via Gmail SMTP
	addr := fmt.Sprintf("%s:%s",s.host,s.port)
	err := smtp.SendMail(
		addr,
		auth,
		s.email,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
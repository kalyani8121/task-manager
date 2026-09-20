package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type EmailSender struct {
	apiKey    string
	fromEmail string
}

func NewEmailSender(apiKey, fromEmail string) *EmailSender {
	return &EmailSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}
}

type PriorityTaskEmail struct {
	Rank    int
	Title   string
	Urgency string
	DueDate string
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

func (s *EmailSender) SendPriorityQueueEmail(
	toEmail string,
	userName string,
	tasks []PriorityTaskEmail,
) error {
	subject := "Your Daily Task Priority Queue"

	overdueCount := 0
	for _, t := range tasks {
		if strings.Contains(t.Urgency, "OVERDUE") {
			overdueCount++
		}
	}

	if overdueCount > 0 {
		subject = fmt.Sprintf(
			"URGENT: %d Overdue Tasks Need Attention!",
			overdueCount,
		)
	}

	body := buildEmailBody(userName, tasks)
	return s.sendEmail(toEmail, subject, body)
}

func buildEmailBody(userName string, tasks []PriorityTaskEmail) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Hi %s!\n\n", userName))
	sb.WriteString("Here is your daily priority queue:\n")
	sb.WriteString("================================\n\n")

	if len(tasks) == 0 {
		sb.WriteString("Great job! No pending tasks!\n")
		return sb.String()
	}

	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf("RANK %d: %s\n", task.Rank, task.Title))
		sb.WriteString(fmt.Sprintf("  Status: %s\n", task.Urgency))
		if task.DueDate != "" {
			sb.WriteString(fmt.Sprintf("  Due: %s\n", task.DueDate))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("================================\n")
	sb.WriteString("Complete RANK 1 first!\n\n")
	sb.WriteString("- Task Manager System\n")

	return sb.String()
}

func (s *EmailSender) sendEmail(
	to string,
	subject string,
	body string,
) error {
	reqBody := resendRequest{
		From:    fmt.Sprintf("Task Manager <%s>", s.fromEmail),
		To:      []string{to},
		Subject: subject,
		Text:    body,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.resend.com/emails",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("resend API error: status %d",
			resp.StatusCode)
	}

	return nil
}
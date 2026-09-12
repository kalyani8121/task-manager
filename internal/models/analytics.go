package models

// Summary shows overall task statistics
type AnalyticsSummary struct {
	TotalTasks       int `json:"total_tasks"`    
	Completed        int `json:"completed"`
	Pending          int `json:"pending"`
	InProgress       int `json:"in_progress"`
	Cancelled        int  `json:"cancelled"`
	Overdue          int  `json:"overdue"`
	CompletedPercent float64 `json:"completed_percent"`
}

// StatusCount shows count for one status
type StatusCount struct{
	Status string `json:"status" db:"status"`
	Count  int    `json:"count" db:"count"`
}

// OverdueTask shows a task that is past due date
type OverdueTask struct {
	ID          string `json:"id"            db:"id"`
	Title       string `json:"title"         db:"title"`
	DueDate     string `json:"due_date"      db:"due_date"`
	DaysOverdue int    `json:"days_overdue"  db:"days_overdue"`
	Status      string `json:"status"        db:"status"`
}

// OverdueResponse is what we send back
type OverdueResponse struct {
	OverdueCount int           `json:"overdue_count"`
	Tasks        []OverdueTask `json:"tasks"`
}
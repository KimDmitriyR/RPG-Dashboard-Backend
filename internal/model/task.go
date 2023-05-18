package model

type Task struct {
	ID             int    `json:"id"`
	Email_curator  string `json:"email_curator"`
	Email_employee string `json:"email_employee"`
	Description    string `json:"description"`
}

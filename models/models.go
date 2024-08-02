package models

type EmailReq struct {
	Email string `json:"email"`
}

type SendEmailReq struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

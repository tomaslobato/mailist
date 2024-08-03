package models

type EmailReq struct {
	Email string `json:"email"`
}

type SendEmailReq struct {
	AdminCode string `json:"adminCode"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

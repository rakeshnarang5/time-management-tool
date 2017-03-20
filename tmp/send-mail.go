package main 

import (
	"fmt"
	"net/smtp"
)

type EmailUser struct {
    Username    string
    Password    string
    EmailServer string
    Port        int
}

emailUser := &EmailUser{'yourGmailUsername', 'password', 'smtp.gmail.com', 587}

auth := smtp.PlainAuth("",
    emailUser.Username,
    emailUser.Password,
    emailUser.EmailServer
)




package qmailer

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

var (
	host              string
	port              int
	username          string
	password          string
	encryption        mail.Encryption
	connectionTimeout time.Duration
	sendTimeout       time.Duration
)

func getSmtpClient() (*mail.SMTPClient, error) {
	client := mail.NewSMTPClient()

	//SMTP Client
	client.Host = host
	client.Port = port
	client.Username = username
	client.Password = password
	client.Encryption = encryption
	client.ConnectTimeout = connectionTimeout
	client.SendTimeout = sendTimeout

	//For authentication, you can use AuthPlain, AuthLogin or AuthCRAMMD5
	client.Authentication = mail.AuthPlain

	//KeepAlive true because the connection need to be open for multiple emails
	//For avoid inactivity timeout, every 30 second you can send a NO OPERATION command to smtp client
	//use smtpClient.Client.Noop() after 30 second of inactivity in this example
	client.KeepAlive = true

	//Connect to client
	smtpClient, err := client.Connect()
	if err != nil {
		return nil, err
	}

	return smtpClient, nil
}

package smtp

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

type Config struct {
	MailHost string
	MailPort int
	MailUsername string
	MailPassword string
	MailEncryption mail.Encryption
	MailConnectionTimeout time.Duration
	MailSendTimeout time.Duration
	SmtpMaxConnections uint32
	SmtpAttempts uint32
	SmtpTimeout time.Duration
}

func Setup(cfg *Config) {
	host = cfg.MailHost
	port = cfg.MailPort
	username = cfg.MailUsername
	password = cfg.MailPassword
	encryption = cfg.MailEncryption
	connectionTimeout = cfg.MailConnectionTimeout
	sendTimeout = cfg.MailSendTimeout

	once.Do(func() {
		poolInstance = &smtpPool{
			0,
			0,
			0,
			cfg.SmtpMaxConnections,
			cfg.SmtpAttempts,
			cfg.SmtpTimeout,
		}
	})
}

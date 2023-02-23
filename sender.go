package smtp

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"sync"
)

//TODO add company signature

type EmailData struct {
	From, Subject, Sender, ReplyTo string
	Addresses, Attachments         []string
	TemplateData                   map[string]string
}

var in chan *mail.Email

func init() {
	in = make(chan *mail.Email)
}

func SendEmail(templateName string, data *EmailData) (chan *result, error) {
	var (
		mu   sync.Mutex
		body bytes.Buffer
	)
	out := make(chan *result)

	paths := []string{
		"mailer\\templates\\base.html", "mailer\\templates\\" + templateName, "mailer\\" +
			"templates\\styles.html"}

	templateData, err := template.ParseFiles(paths...)
	//template, err := ParseTemplateDir("./mailer/templates")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse template")
	}

	templateData = templateData.Lookup(templateName)
	err = templateData.Execute(&body, &data.TemplateData)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not execute template")
	}

	e := mail.NewMSG()

	e.SetFrom(data.From).
		SetSubject(data.Subject)

	//Get from each mail
	e.GetFrom()
	e.SetBody(mail.TextHTML, body.String())

	if data.Sender != "" {
		e.SetSender(data.Sender)
	}
	if data.ReplyTo != "" {
		e.SetReplyTo(data.ReplyTo)
	}
	//if data.ReturnPath != "" {
	//	e.SetReturnPath(data.ReturnPath)
	//}
	//if len(data.Cc) > 0 {
	//	e.AddCc(data.Cc...)
	//}
	//if len(data.Sender) > 0 {
	//	e.AddBcc(data.Bcc...)
	//}

	for _, attachment := range data.Attachments {
		e.Attach(&mail.File{FilePath: attachment})
	}

	//Send with high priority
	e.SetPriority(mail.PriorityHigh)

	// always check error after send
	if e.Error != nil {
		fmt.Println("EmailSendError send in sendEmail func")
		return nil, e.Error
	}

	for _, address := range data.Addresses {
		e.AddTo(address)

		mu.Lock()
		if poolInstance.Active < poolInstance.GoMax && poolInstance.Pending == 0 {
			go poolInstance.Connect(in, out)
		}
		mu.Unlock()

		in <- e
	}

	//TODO close(in)

	return out, nil
}

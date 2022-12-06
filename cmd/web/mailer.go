package main

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	Domain     string
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
	From       string
	FromName   string
	Wait       *sync.WaitGroup
	JobCh      chan MailerJob
	ErrorCh    chan error
	DoneCh     chan bool
}

type MailerJob struct {
	From          string
	FromName      string
	To            string
	Subject       string
	Attachments   []string // full pathnames
	AttachmentMap map[string]string
	Data          any            // body
	DataMap       map[string]any // convenient way to get data into the template we're using
	Template      string
}

// Init in App
func (app *App) createMailer() {
	errorCh := make(chan error)
	jobCh := make(chan MailerJob, 100)
	doneCh := make(chan bool)

	app.Mailer = &Mailer{
		Domain:     "localhost",
		Host:       "localhost",
		Port:       1025,
		Encryption: "none",
		From:       "info@mycompany.com",
		FromName:   "Info",
		Wait:       app.Wait,
		ErrorCh:    errorCh,
		JobCh:      jobCh,
		DoneCh:     doneCh,
	}
}

func buildHTMLMessage(msg MailerJob) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template)
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}
	htmlStr := tpl.String()
	htmlStr, err = inlineCSS(htmlStr)
	if err != nil {
		return "", err
	}

	return htmlStr, nil
}

func buildPlainMessage(msg MailerJob) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.plain.gohtml", msg.Template)
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}
	plainStr := tpl.String()
	return plainStr, nil
}

func inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}
	return html, nil
}

func asEncryptionType(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

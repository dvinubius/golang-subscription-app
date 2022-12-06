package main

import (
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// listen in App
func (app *App) runMailer() {
	for {
		select {
		case job := <-app.Mailer.JobCh:
			go app.Mailer.execJob(job, app.Mailer.ErrorCh)
		case err := <-app.Mailer.ErrorCh:
			app.ErrorLog.Println(err)
		case <-app.Mailer.DoneCh:
			return
		}
	}
}

func (app *App) sendEmail(job MailerJob) {
	app.Wait.Add(1)
	app.Mailer.JobCh <- job
}

func (m *Mailer) execJob(mJob MailerJob, errorCh chan error) {
	defer m.Wait.Done()
	if mJob.Template == "" {
		mJob.Template = "mail"
	}
	if mJob.From == "" {
		mJob.From = m.From
	}
	if mJob.FromName == "" {
		mJob.FromName = m.FromName
	}
	if mJob.AttachmentMap == nil {
		mJob.AttachmentMap = make(map[string]string)
	}

	if len(mJob.DataMap) == 0 {
		mJob.DataMap = make(map[string]any)
	}

	mJob.DataMap["message"] = mJob.Data

	formattedMessage, err := buildHTMLMessage(mJob)
	if err != nil {
		errorCh <- err
	}
	plainMessage, err := buildPlainMessage(mJob)
	if err != nil {
		errorCh <- err
	}

	smtpClient := mail.NewSMTPClient()
	smtpClient.Host = m.Host
	smtpClient.Port = m.Port
	smtpClient.Username = m.Username
	smtpClient.Password = m.Password
	smtpClient.Encryption = asEncryptionType(m.Encryption)
	smtpClient.KeepAlive = false
	smtpClient.ConnectTimeout = 10 * time.Second
	smtpClient.SendTimeout = 10 * time.Second

	clientConn, err := smtpClient.Connect()
	if err != nil {
		errorCh <- err
	}

	email := mail.NewMSG()
	email.SetFrom(mJob.From).AddTo(mJob.To).SetSubject(mJob.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.SetBody(mail.TextHTML, formattedMessage)

	if len(mJob.Attachments) > 0 {
		for _, x := range mJob.Attachments {
			email.AddAttachment(x)
		}
	}

	if len(mJob.AttachmentMap) > 0 {
		for key, value := range mJob.AttachmentMap {
			email.AddAttachment(value, key)
		}
	}

	err = email.Send(clientConn)
	if err != nil {
		errorCh <- err
	}
}

func (app *App) shutdownMailer() {
	app.Mailer.DoneCh <- true
	close(app.Mailer.JobCh)
	close(app.Mailer.DoneCh)
}

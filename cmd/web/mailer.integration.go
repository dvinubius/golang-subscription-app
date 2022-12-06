package main

import "github.com/dvinubius/golang-subscription-app/cmd/web/mailer"

// listen in App
func (app *App) runMailer() {
	for {
		select {
		case job := <-app.Mailer.JobCh:
			go app.Mailer.ExecJob(job, app.Mailer.ErrorCh)
		case err := <-app.Mailer.ErrorCh:
			app.ErrorLog.Println(err)
		case <-app.Mailer.DoneCh:
			return
		}
	}
}

func (app *App) sendEmail(job mailer.MailerJob) {
	app.Wait.Add(1)
	app.Mailer.JobCh <- job
}

func (app *App) shutdownMailer() {
	app.Mailer.DoneCh <- true
	close(app.Mailer.JobCh)
	close(app.Mailer.DoneCh)
}

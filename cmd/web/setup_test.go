package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dvinubius/golang-subscription-app/cmd/web/mailer"
	"github.com/dvinubius/golang-subscription-app/data"
)

var testApp App

func TestMain(m *testing.M) {
	gob.Register(data.User{})

	tempPath = "./../../tmp"
	pathToManual = "./../../pdf"

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	testApp = App{
		Session:     session,
		DB:          nil,
		InfoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:    log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Wait:        &sync.WaitGroup{},
		Models:      data.TestNew(nil),
		ErrorCh:     make(chan error),
		ErrorDoneCh: make(chan bool),
	}

	// dummy mailer
	testApp.Mailer = &mailer.Mailer{
		Wait:    testApp.Wait,
		ErrorCh: make(chan error),
		JobCh:   make(chan mailer.MailerJob, 100),
		DoneCh:  make(chan bool),
	}

	go func() {
		for {
			select {
			case <-testApp.Mailer.JobCh:
				testApp.Wait.Done()
			case <-testApp.Mailer.ErrorCh:
			case <-testApp.Mailer.DoneCh:

			}
		}
	}()

	go func() {
		for {
			select {
			case <-testApp.ErrorCh:
			case <-testApp.ErrorDoneCh:
				return
			}
		}
	}()

	os.Exit(m.Run())
}

func getCtxWithSession(r *http.Request) context.Context {
	// even if there is no session token in the header, we do this in order for the ctx to be
	// context.Background.WithValue(type scs.contextKey, val <not Stringer>)
	// as opposed to merely
	// context.Background
	ctx, err := testApp.Session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx

}

package main

import (
	"database/sql"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/dvinubius/golang-subscription-app/cmd/web/mailer"
	"github.com/dvinubius/golang-subscription-app/data"
)

type App struct {
	Session     *scs.SessionManager
	DB          *sql.DB
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	Wait        *sync.WaitGroup
	Models      data.Models
	Mailer      *mailer.Mailer
	ErrorCh     chan error
	ErrorDoneCh chan bool
}

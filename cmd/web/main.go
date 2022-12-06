package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dvinubius/golang-subscription-app/data"
)

const webPort = "80"

func main() {
	// db connect
	db := initDB()
	// sessions
	session := initSession()
	// loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// channels

	// wg
	wg := sync.WaitGroup{}

	// app cfg
	app := App{
		Session:     session,
		DB:          db,
		Wait:        &wg,
		InfoLog:     infoLog,
		ErrorLog:    errorLog,
		Models:      data.New(db),
		ErrorCh:     make(chan error),
		ErrorDoneCh: make(chan bool),
	}

	// mail
	app.createMailer()
	go app.runMailer()

	// listen for signals
	go app.listenForShutdown()
	// listen for errors
	go app.handleErrors()
	// listen for requests
	app.serve()
}

func (app *App) handleErrors() {
	for {
		select {
		case err := <-app.ErrorCh:
			app.ErrorLog.Println(err)
		case <-app.ErrorDoneCh:
			return
		}
	}
}

func (app *App) serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting Web Server on port " + webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *App) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *App) shutdown() {
	app.InfoLog.Println("running cleanup")
	app.Wait.Wait()
	app.InfoLog.Println("closing channels & shutting down app...")
	app.shutdownMailer()
	app.ErrorDoneCh <- true
	close(app.ErrorCh)
	close(app.ErrorDoneCh)
}

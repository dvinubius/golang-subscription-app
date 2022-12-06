package main

import (
	"database/sql"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/dvinubius/golang-subscription-app/data"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB

func initDB() (db *sql.DB) {
	db = attemptConn()
	if db == nil {
		log.Panic("No Database!")
	}

	return
}

func attemptConn() (db *sql.DB) {
	dsn := os.Getenv("DSN")

	for counts := 0; counts < 10; counts++ {
		time.Sleep(time.Second)
		var err error
		db, err = sql.Open("pgx", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err != nil {
			log.Println("Postgres not ready yet...")
			log.Println("Retrying in 1 sec")
			continue
		} else {
			log.Println("Connected!")
			return
		}
	}
	return
}

// SESSIONS

func initSession() (session *scs.SessionManager) {
	gob.Register(data.User{})
	session = scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	return
}

func initRedis() (pool *redis.Pool) {
	pool = &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return
}

# Subscription Web App in Go

> Fullstack Web App in Go showcasing basic flows like 
> - register user (w/ confirmation email)
> - login user
> - logout user 
> - manage a subscription
>   - view available plans
>   - purchase a plan (w/ invoice email & purchased content email)

## Demonstrates go concurrency
Go routines are used for
- send confirmation emails
- send invoice emails
- create custom purchased content (pdf file) & send via email

Showcases **graceful shutdown** given pending actions via channels

## Inspiration
This app was created along the lines of the Udemy Course called "Working With Concurrency in Go".

A few architectural and stylistic improvements were made:
- data models are just data models, not to be abused as receivers of db-query functionality
- regrouping of functionality
- more appropriate names

# Ops
- Docker setup w/ **Postgres**, **Mailhog**, **Redis** 
- `makefile` for easy setup (build, start, stop, restart)

# Architecture
  - `App` struct is the app config and receiver of various funcs
    - route handlers
    - rendering of go templates for frontend
    - integrations of mailer

  - **routing** via chi mux, custom **Middleware** for Auth & SessionLoad 
  - custom mailer package
    - email verification links created with **signer package** `go-alone` - crypto sign message & attach as suffix
  - Redis used for sessions
  - Middleware for authenticated-status and loading sessions
  
  - **error handling** via channels on 2 levels
      - app level: `app.ErrorCh` and `app.ErrorDoneCh`
      - mailer level: `mailer.ErrroCh` and `mailer.ErrorDoneCh`
  
  - DB
      - setup with explicit sequence definitions, pg.catalog, …
      - Models
          - **no gorm**, plain sql queries & custom defined models
          - db operations on tables as functions on modelConnectors : special purpose receivers (`DBUsers` and `DBPlans`)
# Tests

A few patterns for testing
- routes
- renderer
- auth
- handlers

## Data Models
As opposed to using the repository Pattern, handler tests use custom test modelConnectors (db independent) which replicate the interface of the production models. Quite a lot of duplicated code, expensive to maintain, not optimal. Definitely preferring repository model

# Todo

## Functionality 
Add payment solution

## Architecture
Repository pattern

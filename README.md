# Subscription Web App in Go

Fullstack Web App in Go showcasing basic flows like 
- register (confirmation email)
- login
- logout
- manage subscription
  - view available plans
  - purchase a plan (invoice email & purchased content email)

## Demonstrates go concurrency. 
Go routines are used for
- send confirmation emails
- send invoice emails
- create custom purchased content (pdf file) & send via email

Showcasing **graceful shutdown** given pending actions via channels

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
      - setup - explicit sequence definitions, pg.catalog, …
      - Models
          - **no gorm**, plain sql queries & custom defined models
          - db operation on table as method on model struct
          - **empty struct pattern** for entry-independent queries (`app.Models.User.getAll()` is a function call on receiver `data.User{}`)

# Tests

A few patterns for testing
- routes
- renderer
- auth
- handlers

## Data Models
As opposed to using the repository Pattern, handler tests are use custom testmodels (db independent) which reflect the production models. Quite a lot of duplicated code, not optimal.

# Todo

## Functionality 
Add payment solution

## Architecture
Repository pattern

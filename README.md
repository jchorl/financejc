# FinanceJC
This is a project I've been working on to track my finances

## Getting Started
1. Install `docker`
2. Run `make`

The app should be served on https://localhost:4443

## Importing Data
1. Place one or more QIF files in a folder called import in the root directory
2. Go to https://localhost:4443
3. Log in with Google if necessary
4. Click the import button

## Makefile Targets
`make` creates the Docker network, builds the DB, then the UI, then the Go server, and then serves

`make ui` builds the UI

`make ui-watch` builds the UI and watches for changes

`make db` builds the database

`make connect-db` creates a postgres container running psql and connects to the database

`make serve` serves the site

`make build` builds the Go server

`make clean` kills the db, webserver, postgres container connected to the db, and network
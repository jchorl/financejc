# FinanceJC
This is a project I've been working on to track my finances

## Getting Started
1. Install `docker`
2. Run `make`

The app should be served on localhost:8080

## Importing Data
1. Place one or more QIF files in a folder called import in the root directory
2. Go to localhost:8080
3. Log in with Google if necessary
4. Click the import button

## Makefile Targets
`make` builds the UI, then builds and starts the server, watching for changes

`make ui` builds the UI

`make ui-watch` builds the UI and watches for changes

`make serve` serves the site and watches for changes

`make build` builds the main server image

If developing server code, run `make` and leave it running

For frontend, run `make build`, and when that completes, run `make ui-watch` and `make serve` in separate terminals to continuously rebuild the UI and update the server on change
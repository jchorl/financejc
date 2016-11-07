# FinanceJC
This is a project I've been working on to track my finances.
Features:
- Create accounts in many different currencies
- Create categorized transactions under accounts
- Suggestions for transactions while you type based on previous transactions
- Scheduled transactions at a fixed interval or on a set day of the week/month/year
- REST APIs to allow other frontends and better visualization of finances
- Import data from QIF files
- Login with Google

## Getting Started
1. Install `docker`
2. Run `make dev`

The app should be served on https://localhost

## Architecture
### Webserver
The webserver is written in Go. It uses popular frameworks such as github.com/labstack/echo and github.com/Sirupsen/logrus. All SQL is done using database/sql and queries are written in raw SQL.

### Frontend
The frontend is written in ES6 with React/Redux/Immutable etc. It is transpiled using Webpack, combined with CSS and bundled into a bundle.js file. Webpack adds a link to the bundle into index.html and also minifies index.html.

### Database
The app uses a Postgres database. The init script can be found in /db.

### Elasticsearch
Elasticsearch indexes all transactions to provide suggestions as a user types. When typing the name for a transaction, elasticsearch will query previous transactions and provide suggestions. Clicking a suggestion will fill all fields, except for the date. When typing the category, elasticsearch will query categories and provide suggestions. Clicking a suggestion will only fill the category field.

### Nginx
Nginx is used for serving static files and load balancing. TLS certs are also managed in the nginx container using LetsEncrypt.

## Importing Data
1. Log in if necessary
2. Hover over your email address in the top right and click import
3. Select a QIF file

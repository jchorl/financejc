CREATE TABLE users (
    id serial PRIMARY KEY,
    googleId varchar(40) UNIQUE,
    email varchar(40) UNIQUE
);

CREATE TABLE accounts (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    currency varchar(3) NOT NULL,
    userId integer NOT NULL references users(id) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE transactions (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    occurred date NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    relatedTransactionId integer references transactions(id) DEFERRABLE INITIALLY DEFERRED,
    accountId integer NOT NULL references accounts(id) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE recurringTransactions (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    nextOccurs date NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    accountId integer NOT NULL references accounts(id) DEFERRABLE INITIALLY DEFERRED,

    scheduleType varchar(20) NOT NULL,
    secondsBetween integer,
    dayOf integer,
    secondsBeforeToPost integer NOT NULL
);

CREATE TABLE transactionTemplates (
    id serial PRIMARY KEY,
    templateName varchar(40) NOT NULL,
    name varchar(100) NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    accountId integer NOT NULL references accounts(id) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX ON users(googleId);
CREATE INDEX ON accounts(userId);
CREATE INDEX ON transactions(accountId, occurred DESC, id);
CREATE INDEX ON recurringTransactions(accountId);
CREATE INDEX ON recurringTransactions((nextOccurs - interval '1 second' * secondsBeforeToPost));
CREATE INDEX ON transactionTemplates(accountId);

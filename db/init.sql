CREATE TABLE users (
    id serial PRIMARY KEY,
    google_id varchar(40) UNIQUE,
    email varchar(40) UNIQUE
);

CREATE TABLE accounts (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    currency varchar(3) NOT NULL,
    user_id integer NOT NULL references users(id) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE transactions (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    occurred date NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    related_transaction_id integer references transactions(id) DEFERRABLE INITIALLY DEFERRED,
    account_id integer NOT NULL references accounts(id) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE recurring_transactions (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    next_occurs date NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    account_id integer NOT NULL references accounts(id) DEFERRABLE INITIALLY DEFERRED,

    schedule_type varchar(20) NOT NULL,
    seconds_between integer,
    day_of integer,
    seconds_before_to_post integer NOT NULL
);

CREATE TABLE templates (
    id serial PRIMARY KEY,
    template_name varchar(40) NOT NULL,
    name varchar(100) NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    account_id integer NOT NULL references accounts(id) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX ON users(google_id);
CREATE INDEX ON accounts(user_id);
CREATE INDEX ON transactions(account_id, occurred DESC, id);
CREATE INDEX ON recurring_transactions(account_id);
CREATE INDEX ON recurring_transactions((next_occurs - interval '1 second' * seconds_before_to_post));
CREATE INDEX ON templates(account_id);

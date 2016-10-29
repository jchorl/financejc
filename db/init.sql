CREATE TABLE users (
    id serial PRIMARY KEY,
    googleId varchar(40)
);

CREATE TABLE accounts (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    currency varchar(3) NOT NULL,
    userId integer NOT NULL references users(id)
);

CREATE TABLE transactions (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    occurred date NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    relatedTransaction integer references transactions(id),
    account integer NOT NULL references accounts(id)
);

CREATE TABLE recurring_ransactions (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    next_occurs date NOT NULL,
    category varchar(100),
    amount integer NOT NULL,
    note varchar(256),
    account integer NOT NULL references accounts(id),

    schedule_type varchar(20) NOT NULL,
    seconds_between integer,
    day_of integer,
    seconds_before_to_post integer NOT NULL
);

CREATE INDEX ON accounts(userId);
CREATE INDEX ON transactions(account);

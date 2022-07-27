-- Represents financial institutions like bank names or investment firms
CREATE TABLE organizations
(
    id         integer primary key autoincrement,
    name       varchar   not null,
    is_deleted tinyint            default 0,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp -- needs to be manually updated on updates
);

-- represents the accounts
CREATE TABLE accounts
(
    id              integer primary key autoincrement,
    name            varchar   not null,
    type            varchar   not null,                          -- could be assets like savings, cash, investment or liabilities like credit card, loan
    subtype         varchar,                                     -- could be sub types of investments like MFs, stocks or insurance
    org_id          integer references organizations (id),
    current_balance NUMERIC,
    is_deleted      tinyint            default 0,
    created_ts      timestamp not null default current_timestamp,
    updated_ts      timestamp not null default current_timestamp -- needs to be manually updated on updates
);

-- represents the ledger of transactions
CREATE TABLE transactions
(
    id               integer primary key autoincrement,
    transaction_time timestamp not null,
    type             timestamp not null,                          -- transfer or normal
    account_id       integer references accounts (id),
    description      varchar,
    memo             varchar,                                     -- personal notes about a transaction
    amount           NUMERIC,                                     -- negative value indicates expense, positive value indicates income
    currency         varchar,                                     -- to handle multi currency computations
    payee            integer references payees (id),
    category_id      varchar references categories (id),
    tags             varchar,                                     -- no foreign key because it's a json column
    is_deleted       tinyint            default 0,
    created_ts       timestamp not null default current_timestamp,
    updated_ts       timestamp not null default current_timestamp -- needs to be manually updated on updates
);

-- payees for another dimension of tracking
CREATE TABLE payees
(
    id         integer primary key autoincrement,
    name       varchar   not null,
    is_deleted tinyint            default 0,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp -- needs to be manually updated on updates
);

-- hierarchical representation of categories
CREATE TABLE categories
(
    id         integer primary key autoincrement,
    name       varchar   not null,
    parent     integer            default name,
    is_deleted tinyint            default 0,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp -- needs to be manually updated on updates
);

-- represents a tag on each transaction. Does not support soft deletes.
CREATE TABLE tags
(
    name       varchar primary key,
    created_ts timestamp not null default current_timestamp,
    updated_ts timestamp not null default current_timestamp -- needs to be manually updated on updates
);
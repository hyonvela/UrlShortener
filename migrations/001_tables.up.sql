CREATE TABLE IF NOT EXISTS short_url (
    id bigint PRIMARY KEY,
    short_url VARCHAR(10) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS long_url (
    id bigint PRIMARY KEY,
    long_url TEXT UNIQUE NOT NULL
);
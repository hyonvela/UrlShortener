CREATE TABLE IF NOT EXISTS url_mappings (
    id bigint PRIMARY KEY,
    short_url VARCHAR(10) UNIQUE NOT NULL,
    long_url TEXT UNIQUE NOT NULL
);

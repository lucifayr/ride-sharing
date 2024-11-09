DROP TABLE users;


CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('google', 'microsoft'))
);

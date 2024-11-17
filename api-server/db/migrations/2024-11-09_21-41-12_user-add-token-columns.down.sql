PRAGMA foreign_keys = off;


BEGIN TRANSACTION;


ALTER TABLE users
RENAME TO _users_old;


CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('google', 'microsoft'))
);


INSERT INTO
    users (id, name, email, provider)
SELECT
    id,
    name,
    email,
    provider
FROM
    _users_old;


DROP TABLE _users_old;


COMMIT;


PRAGMA foreign_keys = ON;

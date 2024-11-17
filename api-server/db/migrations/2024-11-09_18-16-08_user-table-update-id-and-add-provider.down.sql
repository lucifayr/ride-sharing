PRAGMA foreign_keys = off;


BEGIN TRANSACTION;


ALTER TABLE users
RENAME TO _users_old;


CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL
);


INSERT INTO
    users (id, name, email)
SELECT
    CAST(id AS INTEGER),
    name,
    email
FROM
    _users_old;


DROP TABLE _users_old;


COMMIT;


PRAGMA foreign_keys = ON;

-- See sqlc docs for more information:
-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html#schema-and-queries
--
-- name: UsersCreate :one
INSERT INTO
    users (id, name, email, provider)
VALUES
    (?, ?, ?, ?) RETURNING *;


-- name: UsersGetById :one
SELECT
    *
FROM
    users
WHERE
    id = ?;


-- name: UsersUpdateNameAndEmail :one
UPDATE users
SET
    name = ?,
    email = ?
WHERE
    id = ? RETURNING *;


-- name: UsersSetTokens :exec
UPDATE users
SET
    access_token = ?,
    refresh_token = ?
WHERE
    id = ?;


-- name: UsersSetBlocked :exec
UPDATE users
SET
    is_blocked = ?
WHERE
    id = ?;

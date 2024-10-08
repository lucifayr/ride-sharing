-- See sqlc docs for more information:
-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html#schema-and-queries
--
-- name: UsersGetById :one
SELECT
    *
FROM
    users
WHERE
    id = ?;


-- name: UsersUpdateName :one
UPDATE users
SET
    name = ?
WHERE
    id = ? RETURNING *;

-- See sqlc docs for more information:
-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html#schema-and-queries
--
-- name: GroupsCreate :one
INSERT INTO
    ride_groups (name, description, created_by)
VALUES
    (?, ?, ?) RETURNING *;


-- name: GroupsGetMany :many
SELECT
    id,
    name,
    description,
    created_by
FROM
    ride_groups
ORDER BY
    name
LIMIT
    50
OFFSET
    ?;

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


-- name: GroupsGetById :one
SELECT
    id,
    name,
    description,
    created_by
FROM
    ride_groups
WHERE
    id = ?;


-- name: GroupsUpdateName :exec
UPDATE ride_groups
SET
    name = ?
WHERE
    id = ?;


-- name: GroupsUpdateDescription :exec
UPDATE ride_groups
SET
    description = ?
WHERE
    id = ?;


-- name: GroupsMembersGet :many
SELECT
    group_id,
    user_id,
    u.email,
    join_status
FROM
    ride_group_members
    INNER JOIN users u ON u.id = user_id
WHERE
    group_id = ?
ORDER BY
    (
        SELECT
            gso.ordering
        FROM
            ride_group_members_join_status_ordering gso
        WHERE
            gso.status = join_status
    );


-- name: GroupsMembersJoin :exec
INSERT INTO
    ride_group_members (group_id, user_id)
VALUES
    (?, ?);


-- name: GroupsMembersSetStatus :exec
UPDATE ride_group_members
SET
    join_status = ?
WHERE
    group_id = ?
    AND user_id = ?;

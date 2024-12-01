-- See sqlc docs for more information:
-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html#schema-and-queries
--
-- name: RidesCreate :one
INSERT INTO
    rides (
        location_from,
        location_to,
        tacking_place_at,
        created_by,
        driver,
        transport_limit
    )
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING *;


-- name: RidesGetMany :many
SELECT
    rides.id,
    location_from,
    location_to,
    tacking_place_at,
    created_by,
    created_at,
    transport_limit,
    driver,
    users.email AS driver_email
FROM
    rides
    INNER JOIN users ON rides.driver = users.id
ORDER BY
    created_at DESC
LIMIT
    50
OFFSET
    ?;


-- name: RidesGetById :one
SELECT
    rides.id,
    location_from,
    location_to,
    tacking_place_at,
    created_by,
    created_at,
    transport_limit,
    driver,
    users.email AS driver_email
FROM
    rides
    INNER JOIN users ON rides.driver = users.id
WHERE
    rides.id = ?;

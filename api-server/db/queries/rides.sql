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
    ud.email AS driver_email,
    uc.email AS created_by_email
FROM
    rides
    INNER JOIN users ud ON rides.driver = ud.id
    INNER JOIN users uc ON rides.created_by = uc.id
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
    ud.email AS driver_email,
    uc.email AS created_by_email
FROM
    rides
    INNER JOIN users ud ON rides.driver = ud.id
    INNER JOIN users uc ON rides.created_by = uc.id
WHERE
    rides.id = ?;

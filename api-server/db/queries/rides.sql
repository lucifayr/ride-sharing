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
    *
FROM
    rides
ORDER BY
    created_at DESC
LIMIT
    50
OFFSET
    ?;

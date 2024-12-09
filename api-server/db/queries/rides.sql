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
    (?, ?, ?, ?, ?, ?) RETURNING id;


-- name: RidesCreateEvent :exec
INSERT INTO
    ride_events (
        ride_id,
        location_from,
        location_to,
        driver,
        tacking_Place_at,
        transport_limit
    )
VALUES
    (?, ?, ?, ?, ?, ?);


-- name: RidesCreateSchedule :one
INSERT INTO
    ride_schedules (ride_id, INTERVAL, unit)
VALUES
    (?, ?, ?) RETURNING id;


-- name: RidesCreateScheduleWeekday :exec
INSERT INTO
    ride_schedule_weekdays (ride_schedule_id, weekday)
VALUES
    (?, ?);


-- name: RidesMarkPastEventsDone :exec
UPDATE ride_events
SET
    status = 'done'
WHERE
    status = 'upcoming'
    AND tacking_place_at <= ?;


-- name: RidesGetLatest :one
SELECT
    r.id AS ride_id,
    re.id AS ride_event_id,
    re.location_from,
    re.location_to,
    re.tacking_place_at,
    r.created_by,
    re.transport_limit,
    re.driver,
    re.status,
    ud.email AS driver_email,
    uc.email AS created_by_email,
    rs.id AS ride_schedule_id,
    rs.unit AS ride_schedule_unit,
    rs.interval AS ride_schedule_interval,
    r.location_from AS base_location_from,
    r.location_to AS base_location_to,
    r.transport_limit AS base_transport_limit,
    r.driver AS base_driver
FROM
    ride_events re
    INNER JOIN rides r ON re.id = r.id
    LEFT OUTER JOIN ride_schedules rs ON rs.ride_id = r.id
    INNER JOIN users ud ON r.driver = ud.id
    INNER JOIN users uc ON r.created_by = uc.id
WHERE
    re.ride_id = ?
    AND re.tacking_place_at = (
        SELECT
            MAX(tacking_place_at)
        FROM
            ride_events
        WHERE
            id = re.id
    );


-- name: RidesGetSchedule :one
SELECT
    id,
    ride_id,
    INTERVAL,
    unit
FROM
    ride_schedules
WHERE
    ride_id = ?;


-- name: RidesGetScheduleWeekdays :many
SELECT
    weekday
FROM
    ride_schedule_weekdays
WHERE
    ride_schedule_id = ?;


-- name: RidesGetMany :many
SELECT
    r.id AS ride_id,
    re.id AS ride_event_id,
    re.location_from,
    re.location_to,
    re.tacking_place_at,
    r.created_by,
    re.transport_limit,
    re.driver,
    re.status,
    ud.email AS driver_email,
    uc.email AS created_by_email,
    rs.id AS ride_schedule_id,
    rs.unit AS ride_schedule_unit,
    rs.interval AS ride_schedule_interval
FROM
    ride_events re
    INNER JOIN rides r ON re.ride_id = r.id
    LEFT OUTER JOIN ride_schedules rs ON rs.ride_id = r.id
    INNER JOIN users ud ON r.driver = ud.id
    INNER JOIN users uc ON r.created_by = uc.id
ORDER BY
    (
        SELECT
            ordering
        FROM
            ride_event_status_ordering
        WHERE
            status = re.status
    ),
    tacking_place_at DESC
LIMIT
    50
OFFSET
    ?;

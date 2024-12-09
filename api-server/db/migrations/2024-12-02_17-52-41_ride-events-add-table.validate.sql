SELECT
    id,
    ride_id,
    location_from,
    location_to,
    driver,
    tacking_place_at
FROM
    ride_events
LIMIT
    1;


SELECT
    status,
    ordering
FROM
    ride_event_status_ordering
LIMIT
    1;

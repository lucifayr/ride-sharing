SELECT
    id,
    location_from,
    location_to,
    tacking_place_at,
    created_by,
    driver
FROM
    rides;


SELECT
    user_id,
    ride_id
FROM
    ride_participants;
SELECT
    id,
    ride_id,
    schedule_interval,
    unit
FROM
    ride_schedules
LIMIT
    1;


SELECT
    ride_schedule_id,
    weekday
FROM
    ride_schedule_weekdays
LIMIT
    1;

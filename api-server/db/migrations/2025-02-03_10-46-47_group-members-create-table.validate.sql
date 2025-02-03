SELECT
    group_id,
    user_id,
    join_status
FROM
    ride_group_members
LIMIT
    1;


SELECT
    status,
    ordering
FROM
    ride_group_members_join_status_ordering
LIMIT
    1;

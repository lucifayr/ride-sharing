-- :require ./no-init-add-three-users.sql
INSERT INTO
    rides (
        id,
        location_from,
        location_to,
        tacking_place_at,
        created_by,
        driver,
        transport_limit
    )
VALUES
    (
        '123',
        'Graz',
        "Wien",
        '2044-11-26T15:18:26Z',
        'NnCaPHQLC9',
        'm6SYNABgAw',
        4
    );

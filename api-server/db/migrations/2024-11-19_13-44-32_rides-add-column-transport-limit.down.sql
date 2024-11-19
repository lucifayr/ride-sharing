PRAGMA foreign_keys = off;


BEGIN TRANSACTION;


ALTER TABLE rides
RENAME TO _rides_old;


CREATE TABLE rides (
    id INTEGER PRIMARY KEY,
    location_from TEXT NOT NULL,
    location_to TEXT NOT NULL,
    tacking_place_at TEXT NOT NULL CHECK (
        tacking_place_at = strftime('%Y-%m-%dT%H:%M:%SZ', tacking_place_at)
    ),
    created_by TEXT NOT NULL,
    driver TEXT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (driver) REFERENCES users (id)
);


INSERT INTO
    rides (
        id,
        location_from,
        location_to,
        tacking_place_at,
        created_by,
        driver
    )
SELECT
    id,
    location_from,
    location_to,
    tacking_place_at,
    created_by,
    driver
FROM
    _rides_old;


DROP TABLE _rides_old;


COMMIT;


PRAGMA foreign_keys = ON;

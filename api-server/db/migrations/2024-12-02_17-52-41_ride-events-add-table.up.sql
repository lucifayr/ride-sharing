CREATE TABLE ride_events (
    id TEXT PRIMARY KEY DEFAULT (lower(hex (randomblob (8)))),
    ride_id TEXT NOT NULL,
    location_from TEXT NOT NULL,
    location_to TEXT NOT NULL,
    driver TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('upcoming', 'done', 'canceled')) DEFAULT ('upcoming'),
    tacking_place_at TEXT NOT NULL CHECK (
        tacking_place_at = strftime ('%Y-%m-%dT%H:%M:%SZ', tacking_place_at)
    ),
    transport_limit INTEGER NOT NULL CHECK (transport_limit > 0),
    FOREIGN KEY (ride_id) REFERENCES rides (id),
    FOREIGN KEY (driver) REFERENCES users (id)
);


CREATE TABLE ride_event_status_ordering (
    status TEXT PRIMARY KEY,
    ordering INTEGER NOT NULL
);


INSERT INTO
    ride_event_status_ordering (status, ordering)
VALUES
    ('upcoming', 32),
    ('done', 64),
    ('canceled', 128);


CREATE TRIGGER ride_event_create_first AFTER INSERT ON rides BEGIN
INSERT INTO
    ride_events (
        ride_id,
        location_from,
        location_to,
        driver,
        status,
        tacking_place_at,
        transport_limit
    )
VALUES
    (
        NEW.id,
        NEW.location_from,
        NEW.location_to,
        NEW.driver,
        'upcoming',
        NEW.tacking_place_at,
        NEW.transport_limit
    );


END;

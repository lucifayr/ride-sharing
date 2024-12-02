CREATE TABLE ride_events (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(8)))),
    ride_id TEXT NOT NULL,
    location_from TEXT NOT NULL,
    location_to TEXT NOT NULL,
    driver TEXT NOT NULL,
    tacking_place_at TEXT NOT NULL CHECK (
        tacking_place_at = strftime('%Y-%m-%dT%H:%M:%SZ', tacking_place_at)
    )
);

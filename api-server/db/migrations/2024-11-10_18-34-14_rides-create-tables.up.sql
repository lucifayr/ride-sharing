CREATE TABLE rides (
    id INTEGER PRIMARY KEY,
    location_from TEXT NOT NULL,
    location_to TEXT NOT NULL,
    tacking_place_at TEXT NOT NULL CHECK (
        tacking_place_at == strftime('%Y-%m-%dT%H:%M:%S', tacking_place_at)
    ),
    created_by TEXT NOT NULL,
    driver TEXT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (driver) REFERENCES users (id)
);


CREATE TABLE ride_participants (
    user_id TEXT NOT NULL,
    ride_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, ride_id)
);

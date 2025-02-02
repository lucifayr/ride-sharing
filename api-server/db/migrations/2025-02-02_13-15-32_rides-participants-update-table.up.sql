DROP TABLE ride_participants;


CREATE TABLE ride_participants (
    user_id TEXT NOT NULL,
    ride_event_id TEXT NOT NULL,
    PRIMARY KEY (user_id, ride_event_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (ride_event_id) REFERENCES ride_events (id)
);

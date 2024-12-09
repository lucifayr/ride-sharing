CREATE TABLE ride_schedules (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(8)))),
    ride_id TEXT NOT NULL,
    schedule_interval INTEGER NOT NULL,
    unit TEXT NOT NULL CHECK (
        unit IN ('days', 'weeks', 'months', 'years', 'weekdays')
    ),
    FOREIGN KEY (ride_id) REFERENCES rides (id)
);


CREATE TABLE ride_schedule_weekdays (
    ride_schedule_id TEXT NOT NULL,
    weekday TEXT NOT NULL CHECK (
        weekday IN (
            'monday',
            'tuesday',
            'wednesday',
            'thursday',
            'friday',
            'saturday',
            'sunday'
        )
    ),
    PRIMARY KEY (ride_schedule_id, weekday),
    FOREIGN KEY (ride_schedule_id) REFERENCES ride_schedules (id)
);

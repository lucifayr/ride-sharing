CREATE TABLE ride_groups (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(8)))),
    name TEXT NOT NULL,
    created_by TEXT NOT NULL,
    description TEXT,
    FOREIGN KEY (created_by) REFERENCES users (id)
);

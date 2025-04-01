CREATE TABLE group_messages (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(8)))),
    group_id TEXT NOT NULL,
    content TEXT NOT NULL,
    sent_by TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    replies_to TEXT,
    FOREIGN KEY (group_id) REFERENCES GROUPS (id),
    FOREIGN KEY (sent_by) REFERENCES users (id),
    FOREIGN KEY (replies_to) REFERENCES group_messages (id)
);

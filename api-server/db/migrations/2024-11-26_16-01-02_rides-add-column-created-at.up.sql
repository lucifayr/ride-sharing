ALTER TABLE rides
ADD created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'));

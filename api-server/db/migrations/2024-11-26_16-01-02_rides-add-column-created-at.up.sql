ALTER TABLE rides
ADD created_at INTEGER DEFAULT (strftime('%s', 'now'));

ALTER TABLE rides
ADD transport_limit INTEGER NOT NULL CHECK (transport_limit > 0);

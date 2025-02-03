CREATE TABLE ride_group_members (
    group_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    join_status TEXT NOT NULL CHECK (join_status IN ('pending', 'member', 'banned')) DEFAULT ('pending'),
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES ride_groups (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);


CREATE TABLE ride_group_members_join_status_ordering (
    status TEXT PRIMARY KEY,
    ordering INTEGER NOT NULL
);


INSERT INTO
    ride_group_members_join_status_ordering (status, ordering)
VALUES
    ('pending', 32),
    ('member', 64),
    ('banned', 128);


INSERT INTO
    ride_group_members (group_id, user_id, join_status)
SELECT
    id,
    created_by,
    'member'
FROM
    ride_groups;


CREATE TRIGGER ride_group_owner_join_as_member INSERT ON ride_groups BEGIN
INSERT INTO
    ride_group_members (group_id, user_id, join_status)
VALUES
    (new.id, new.created_by, 'member');


END;

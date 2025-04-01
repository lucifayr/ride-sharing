-- See sqlc docs for more information:
-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html#schema-and-queries
--
-- name: GroupMessagesGetMany :many
SELECT
    gm.id,
    gm.group_id,
    gm.content,
    gm.sent_by,
    u.email AS sent_by_email,
    gm.created_at,
    gm.replies_to
FROM
    group_messages gm
    INNER JOIN users u ON u.id = sent_by
WHERE
    group_id = ?
ORDER BY
    created_at;


-- name: GroupMessagesCreate :one
INSERT INTO
    group_messages (content, group_id, sent_by, replies_to)
VALUES
    (?, ?, ?, ?) RETURNING *;

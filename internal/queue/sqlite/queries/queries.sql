-- name: GetMessage :one

DELETE FROM
    messages
WHERE id = (
    SELECT id
    FROM messages
    WHERE retry_date <= CURRENT_TIMESTAMP
    ORDER BY insertion_date ASC
    LIMIT 1
    )
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE
    id == sqlc.arg(id);

-- name: DeleteUndeliverableMessage :exec
DELETE FROM messages
WHERE
    retries > sqlc.arg(retries);

-- name: GetCount :one
SELECT
    COUNT(id)
FROM
    messages;

-- name: Insert :exec
INSERT INTO
    messages (
        subject,
        message,
        service_id,
        retry_date,
        retries
    )
VALUES (sqlc.arg(subject), sqlc.arg(message), sqlc.arg(service_id), sqlc.arg(retry_date), sqlc.arg(retries));

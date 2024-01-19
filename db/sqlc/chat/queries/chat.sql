-- name: CreateChannel :one
INSERT INTO channel (
  email_id, channel_type 
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateMembership :one
INSERT INTO membership (
  subscriber_id, subscribed_to_id 
) VALUES (
  $1, $2
)
RETURNING *;

-- name: CreateMessage :one
INSERT INTO messages (
  content, sender_id, receiver_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetMembership :many
SELECT * FROM membership
WHERE subscriber_id = $1;

-- name: GetMessages :many
SELECT * FROM messages
WHERE receiver_id = $1;


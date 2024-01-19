// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: chat.sql

package chatdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createChannel = `-- name: CreateChannel :one
INSERT INTO channel (
  email_id, channel_type 
) VALUES (
  $1, $2
)
RETURNING id, email_id, channel_type
`

type CreateChannelParams struct {
	EmailID     string
	ChannelType Channeltype
}

func (q *Queries) CreateChannel(ctx context.Context, arg CreateChannelParams) (Channel, error) {
	row := q.db.QueryRow(ctx, createChannel, arg.EmailID, arg.ChannelType)
	var i Channel
	err := row.Scan(&i.ID, &i.EmailID, &i.ChannelType)
	return i, err
}

const createMembership = `-- name: CreateMembership :one
INSERT INTO membership (
  subscriber_id, subscribed_to_id 
) VALUES (
  $1, $2
)
RETURNING id, subscriber_id, subscribed_to_id
`

type CreateMembershipParams struct {
	SubscriberID   pgtype.UUID
	SubscribedToID pgtype.UUID
}

func (q *Queries) CreateMembership(ctx context.Context, arg CreateMembershipParams) (Membership, error) {
	row := q.db.QueryRow(ctx, createMembership, arg.SubscriberID, arg.SubscribedToID)
	var i Membership
	err := row.Scan(&i.ID, &i.SubscriberID, &i.SubscribedToID)
	return i, err
}

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (
  content, sender_id, receiver_id
) VALUES (
  $1, $2, $3
)
RETURNING id, content, sender_id, receiver_id, created_at
`

type CreateMessageParams struct {
	Content    pgtype.Text
	SenderID   pgtype.UUID
	ReceiverID pgtype.UUID
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRow(ctx, createMessage, arg.Content, arg.SenderID, arg.ReceiverID)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.SenderID,
		&i.ReceiverID,
		&i.CreatedAt,
	)
	return i, err
}

const getMembership = `-- name: GetMembership :many
SELECT id, subscriber_id, subscribed_to_id FROM membership
WHERE subscriber_id = $1
`

func (q *Queries) GetMembership(ctx context.Context, subscriberID pgtype.UUID) ([]Membership, error) {
	rows, err := q.db.Query(ctx, getMembership, subscriberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Membership
	for rows.Next() {
		var i Membership
		if err := rows.Scan(&i.ID, &i.SubscriberID, &i.SubscribedToID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMessages = `-- name: GetMessages :many
SELECT id, content, sender_id, receiver_id, created_at FROM messages
WHERE receiver_id = $1
`

func (q *Queries) GetMessages(ctx context.Context, receiverID pgtype.UUID) ([]Message, error) {
	rows, err := q.db.Query(ctx, getMessages, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.SenderID,
			&i.ReceiverID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
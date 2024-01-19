package message

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Invictus9999/slack_chat/db/sqlc/chatdb"
	util "github.com/Invictus9999/slack_chat/services/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MessageServer struct {
	dbpool *pgxpool.Pool
	rdb    *redis.Client
}

func (srv *MessageServer) Send(w http.ResponseWriter, r *http.Request) {
	data := &SendRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, util.ErrInvalidRequest(err))
		return
	}

	ctx := context.Background()
	conn, err := srv.dbpool.Acquire(ctx)

	if err != nil {
		render.Render(w, r, util.ErrInvalidRequest(err))
		return
	}

	q := chatdb.New(conn)

	senderId := util.GetPGTypeUUIDFromString(data.SenderId)
	receiverId := util.GetPGTypeUUIDFromString(data.ReceiverId)

	if !senderId.Valid || !receiverId.Valid {
		render.Render(w, r, util.ErrInvalidRequest(errors.New("invalid id format")))
		return
	}

	result, err := q.CreateMessage(ctx, chatdb.CreateMessageParams{
		Content:    pgtype.Text{String: data.Content, Valid: true},
		SenderID:   senderId,
		ReceiverID: receiverId,
	})

	if err != nil {
		render.Render(w, r, util.ErrInvalidRequest(err))
		return
	}

	response := &SendResponse{
		MessageId: util.GetUUIDFromPGTypeUUID(result.ID).String(),
	}

	message := PublishMessage{
		SenderId: data.SenderId,
		Content:  data.Content,
	}

	messageStr, err := json.Marshal(message)

	if err != nil {
		render.Render(w, r, util.ErrInvalidRequest(err))
		return
	}

	srv.rdb.Publish(ctx, data.ReceiverId, messageStr)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, response)
}

func NewMessageRouter(dbpool *pgxpool.Pool) http.Handler {
	server := &MessageServer{
		dbpool: dbpool,
		rdb:    NewRedisClient(),
	}

	r := chi.NewRouter()
	r.Post("/send", server.Send)
	return r
}

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

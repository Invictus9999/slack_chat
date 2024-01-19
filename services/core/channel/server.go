package channel

import (
	"context"
	"net/http"

	"github.com/Invictus9999/slack_chat/db/sqlc/chatdb"
	util "github.com/Invictus9999/slack_chat/services/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChannelServer struct {
	dbpool *pgxpool.Pool
}

func (srv *ChannelServer) CreateChannel(w http.ResponseWriter, r *http.Request) {
	data := &CreateChannelRequest{}

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

	result, err := q.CreateChannel(ctx, chatdb.CreateChannelParams{
		EmailID:     data.Name,
		ChannelType: chatdb.Channeltype(data.Type),
	})

	// TODO: Make a user subscribe to the its own channel

	if err != nil {
		render.Render(w, r, util.ErrInvalidRequest(err))
		return
	}

	response := &CreateChannelResponse{
		Id: util.GetUUIDFromPGTypeUUID(result.ID).String(),
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, response)
}

func NewChannelRouter(dbpool *pgxpool.Pool) http.Handler {
	server := &ChannelServer{
		dbpool: dbpool,
	}

	r := chi.NewRouter()
	r.Post("/create", server.CreateChannel)
	return r
}

package membership

import (
	"context"
	"errors"
	"net/http"

	"github.com/Invictus9999/slack_chat/db/sqlc/chatdb"
	util "github.com/Invictus9999/slack_chat/services/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MembershipServer struct {
	dbpool *pgxpool.Pool
}

func (srv *MembershipServer) Subscribe(w http.ResponseWriter, r *http.Request) {
	data := &SubscribeRequest{}

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

	subcriberId := util.GetPGTypeUUIDFromString(data.SubscriberId)
	subscribedToId := util.GetPGTypeUUIDFromString(data.SubscribedToId)

	if !subcriberId.Valid || !subscribedToId.Valid {
		render.Render(w, r, util.ErrInvalidRequest(errors.New("invalid id format")))
		return
	}

	result, err := q.CreateMembership(ctx, chatdb.CreateMembershipParams{
		SubscriberID:   subcriberId,
		SubscribedToID: subscribedToId,
	})

	if err != nil {
		render.Render(w, r, util.ErrInvalidRequest(err))
		return
	}

	response := &SubscribeResponse{
		SubscriptionId: util.GetUUIDFromPGTypeUUID(result.ID).String(),
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, response)
}

func NewMembershipRouter(dbpool *pgxpool.Pool) http.Handler {
	server := &MembershipServer{
		dbpool: dbpool,
	}

	r := chi.NewRouter()
	r.Post("/subscribe", server.Subscribe)
	return r
}

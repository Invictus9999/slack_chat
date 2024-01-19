package main

import (
	"context"
	"net/http"
	"os"

	"github.com/Invictus9999/slack_chat/services/core/channel"
	"github.com/Invictus9999/slack_chat/services/core/membership"
	"github.com/Invictus9999/slack_chat/services/core/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

const POSTGRESQL_URL = "postgres://postgres:password@localhost:5432/slackchat?sslmode=disable"

func main() {
	dbpool, err := pgxpool.New(context.Background(), POSTGRESQL_URL)
	if err != nil {
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Mount("/channel", channel.NewChannelRouter(dbpool))
	r.Mount("/membership", membership.NewMembershipRouter(dbpool))
	r.Mount("/message", message.NewMessageRouter(dbpool))

	http.ListenAndServe(":3000", r)
}

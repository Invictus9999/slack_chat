package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Invictus9999/slack_chat/db/sqlc/chatdb"
	util "github.com/Invictus9999/slack_chat/services/common"
	"github.com/jackc/pgx/v5"
)

const POSTGRESQL_URL = "postgres://postgres:password@localhost:5432/slackchat?sslmode=disable"

func fetchMembership(userId string) []string {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, POSTGRESQL_URL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	q := chatdb.New(conn)
	userIdPgtype := util.GetPGTypeUUIDFromString(userId)

	if !userIdPgtype.Valid {
		return nil
	}

	result, err := q.GetMembership(ctx, userIdPgtype)

	if err != nil {
		return nil
	}

	membership := make([]string, 0)

	for _, res := range result {
		membership = append(membership, util.GetUUIDFromPGTypeUUID(res.SubscribedToID).String())
	}

	fmt.Println(membership)
	return membership
}

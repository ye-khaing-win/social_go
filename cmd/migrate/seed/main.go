package main

import (
	"github.com/ye-khaing-win/social_go/internal/db"
	"github.com/ye-khaing-win/social_go/internal/env"
	"github.com/ye-khaing-win/social_go/internal/store"
	"log"
)

func main() {
	addr := env.GetStr("DB_ADDR", "postgres://admin:password@localhost/social_db?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}

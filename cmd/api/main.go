package main

import (
	"fmt"
	"github.com/joho/godotenv"
	postgres "github.com/ye-khaing-win/social_go/internal/db"
	"github.com/ye-khaing-win/social_go/internal/env"
	"github.com/ye-khaing-win/social_go/internal/store"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	addr := env.GetInt("ADDR", 8080)
	cfg := config{
		addr:    fmt.Sprintf(":%d", addr),
		version: "1.0.0",
		env:     env.GetStr("ENV", "dev"),
		db: dbConfig{
			addr:         env.GetStr("DB_ADDR", "postgres://admin:password@localhost/social_db?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetStr("DB_MAX_IDLE_TIME", "15m"),
		},
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := postgres.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	s := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  s,
		logger: logger,
	}

	logger.Fatal(app.run(app.mount()))
}

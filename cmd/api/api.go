package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ye-khaing-win/social_go/internal/store"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr    string
	env     string
	version string
	db      dbConfig
	mail    mailConfig
}

type mailConfig struct {
	exp time.Duration
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			//r.With(app.postsContextMiddleware).Get("/{postID}", app.getPostHandler)
			//r.With(app.postsContextMiddleware).Patch("/{postID}", app.updatePostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Patch("/", app.updatePostHandler)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)

			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Post("/follow", app.followUserHandler)
				r.Post("/unfollow", app.unfollowUserHandler)
			})

			r.Route("/feed", func(r chi.Router) {
				r.Use(app.Pagination)
				r.Get("/", app.getUserFeedHandler)
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerHandler)
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)
	return srv.ListenAndServe()
}

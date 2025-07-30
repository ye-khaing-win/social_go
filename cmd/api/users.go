package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/ye-khaing-win/social_go/internal/models"
	"github.com/ye-khaing-win/social_go/internal/store"
	"net/http"
	"strconv"
)

type userContext struct{}

type FollowUserPayload struct {
	UserID int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followingUser := getUserFromContext(r)

	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(
		r.Context(),
		followingUser.ID,
		payload.UserID,
	); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowingUser := getUserFromContext(r)

	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.Followers.Unfollow(
		r.Context(),
		unfollowingUser.ID,
		payload.UserID,
	); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, id)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerErrorResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userContext{}, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *models.User {
	user, _ := r.Context().Value(userContext{}).(*models.User)
	return user
}

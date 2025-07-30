package main

import (
	"errors"
	"github.com/ye-khaing-win/social_go/internal/models"
	"github.com/ye-khaing-win/social_go/internal/store"
	"net/http"
)

type RegisterPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &models.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Users.CreateAndInvite(ctx, user, "1234", app.config.mail.exp); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.conflictResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}
}

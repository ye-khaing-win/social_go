package main

import (
	"net/http"
)

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = app.errorResponse(w, http.StatusInternalServerError, "something went wrong")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = app.errorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = app.errorResponse(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	_ = app.errorResponse(w, http.StatusConflict, err.Error())
}

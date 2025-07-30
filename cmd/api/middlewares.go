package main

import (
	"context"
	q "github.com/ye-khaing-win/social_go/internal/query"
	"net/http"
	"strconv"
)

func (app *application) Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 1
		limit := 10

		query := r.URL.Query()

		if p, err := strconv.Atoi(query.Get("page")); err == nil && p > 0 {
			page = p
		}

		if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
			limit = l
		}

		pg := q.Pagination{
			Limit:  limit,
			Offset: (page - 1) * limit,
		}

		ctx := context.WithValue(r.Context(), q.PgContext{}, pg)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

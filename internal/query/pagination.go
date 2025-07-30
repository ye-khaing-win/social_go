package query

import (
	"context"
)

type PgContext struct{}
type Pagination struct {
	Limit  int
	Offset int
}

func GetPgFromContext(ctx context.Context) Pagination {
	p, _ := ctx.Value(PgContext{}).(Pagination)

	return p
}

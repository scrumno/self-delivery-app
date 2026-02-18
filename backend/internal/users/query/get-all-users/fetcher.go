package get_all_users

import (
	"context"
)

type Fetcher struct {
	UserStorage Query
}

func (f *Fetcher) Fetch(ctx context.Context, q *Query) (any, error) {
	result := &Fetcher{}
	result.UserStorage = *q
	return result, nil
}

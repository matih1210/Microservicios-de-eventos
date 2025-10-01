package event

import "context"

type Event struct {
	ID       string `json:"id"`
	Canceled int64  `json:"canceled"`
}

type Reader interface {
	FindActiveByID(ctx context.Context, id string) (*Event, error)
}

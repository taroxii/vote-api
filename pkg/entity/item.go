package entity

import (
	"context"
	"time"
)

type Item struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	VoteCount   int       `json:"voteCount"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ItemUseCase interface {
	Fetch(c context.Context, cursor string, num int64) (res []Item, nextCursor string, err error)
	Update(c context.Context, d *Item) (err error)
	Insert(c context.Context, d *Item) (err error)
	Delete(c context.Context, id int64) (err error)
	Vote(c context.Context, id int64, user_id int64) (item *Item, err error)
	ClearVote(c context.Context, id int64) (item Item, err error)
}

type ItemRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Item, nextCursor string, err error)
	Insert(ctx context.Context, a *Item) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, ar *Item) (err error)
	Vote(ctx context.Context, id int64, user_id int64) (item *Item, err error)
	ClearVote(ctx context.Context, id int64) (item Item, err error)
}

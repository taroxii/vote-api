package entity

import (
	"context"
	"time"
)

type User struct {
	ID        int64
	Username  string
	VotedItem []Item
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	GetById(ctx context.Context, id int64) (res User, err error)
	GetByUsername(ctx context.Context, username string) (res User, err error)
}

type UserUsecase interface {
	SignIn(ctx context.Context, usr string) (tokens *string, userDM *User, err error)
}

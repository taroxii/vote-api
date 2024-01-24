package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/taroxii/vote-api/pkg/entity"
	"github.com/taroxii/vote-api/pkg/utils/logger"
)

type itemUseCase struct {
	transactionManager TransactionManager
	itemRepository     entity.ItemRepository
	userRepo           entity.UserRepository

	contextTimeout time.Duration
}

func NewItemUsecase(a entity.ItemRepository, userRepo entity.UserRepository, timeout time.Duration) entity.ItemUseCase {
	return &itemUseCase{
		itemRepository: a,
		userRepo:       userRepo,
		contextTimeout: timeout,
	}
}

func (a *itemUseCase) Fetch(c context.Context, cursor string, num int64) (res []entity.Item, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.itemRepository.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *itemUseCase) Vote(c context.Context, id int64, user_id int64) (item *entity.Item, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)

	defer cancel()

	item, err = a.itemRepository.Vote(ctx, id, user_id)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error %v", err))
		return
	}

	return
}

func (iu *itemUseCase) ClearVote(c context.Context, id int64) (item entity.Item, err error) {
	ctx, cancel := context.WithTimeout(c, iu.contextTimeout)
	defer cancel()

	item, err = iu.itemRepository.ClearVote(ctx, id)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Error %v", err))
		return
	}
	return
}

func (a *itemUseCase) Update(c context.Context, item *entity.Item) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.itemRepository.Update(ctx, item)
}

func (a *itemUseCase) Insert(c context.Context, item *entity.Item) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now
	err = a.itemRepository.Insert(ctx, item)
	return
}

func (a *itemUseCase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	return a.itemRepository.Delete(ctx, id)
}

package usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taroxii/vote-api/pkg/entity"
	"github.com/taroxii/vote-api/pkg/mocks"
	"github.com/taroxii/vote-api/pkg/usecase"
	"github.com/taroxii/vote-api/pkg/utils/logger"
)

func TestFetch(t *testing.T) {
	mockItemRepo := new(mocks.ItemRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockItem := entity.Item{
		Name:        "The one",
		Description: "This is Description",
		VoteCount:   0,
	}
	mockListItem := make([]entity.Item, 0)
	mockListItem = append(mockListItem, mockItem)
	type args struct {
		ctx    context.Context
		num    int64
		cursor string
		// req
	}

	type mocks struct {
		listItem     []entity.Item
		mockFunction func()
		userRepo     entity.UserRepository
	}

	type expected struct {
		cursor     string
		itemLength int
		itemReturn []entity.Item
		err        error
	}

	tests := []struct {
		name     string
		args     args
		mocks    mocks
		expected expected
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.TODO(),
				num:    10,
				cursor: "",
			},
			expected: expected{
				cursor:     "next-cursor",
				itemLength: len(mockListItem),
				itemReturn: mockListItem,
			},
			mocks: mocks{
				listItem: mockListItem,
				mockFunction: func() {
					mockItemRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(mockListItem, "next-cursor", nil).Once()
				},
			},
		},
		{
			name: "should return error",
			args: args{
				ctx:    context.TODO(),
				num:    10,
				cursor: "",
			},
			expected: expected{
				err:        sql.ErrNoRows,
				cursor:     "",
				itemLength: len([]entity.Item(nil)),
				itemReturn: []entity.Item(nil),
			},
			mocks: mocks{
				listItem: []entity.Item(nil),
				mockFunction: func() {
					mockItemRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return([]entity.Item(nil), "", sql.ErrNoRows).Once()
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mocks.mockFunction()
			u := usecase.NewItemUsecase(mockItemRepo, mockUserRepo, time.Second*5)
			list, nextCursor, err := u.Fetch(test.args.ctx, test.args.cursor, test.args.num)
			assert.Equal(t, test.expected.itemReturn, list)

			assert.Len(t, list, test.expected.itemLength)
			assert.Equal(t, test.expected.cursor, nextCursor)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, nextCursor)
			}

			mockItemRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}

}

func TestVote(t *testing.T) {
	mockItemRepo := new(mocks.ItemRepository)
	mockUserRepo := new(mocks.UserRepository)
	logger.InitializeLogger()
	mockItem := entity.Item{
		ID:          1,
		Name:        "The one",
		Description: "This is Description",
		VoteCount:   0,
	}

	type args struct {
		ctx     context.Context
		id      int64
		user_id int64
	}

	type mocks struct {
		item         *entity.Item
		mockFunction func()
	}

	type expected struct {
		itemReturn *entity.Item
		err        error
	}

	tests := []struct {
		name     string
		args     args
		mocks    mocks
		expected expected
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx:     context.TODO(),
				id:      1,
				user_id: 1,
			},
			expected: expected{
				itemReturn: &mockItem,
			},
			mocks: mocks{
				item: &mockItem,
				mockFunction: func() {
					mockItemRepo.On("Vote", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(&mockItem, nil).Once()
				},
			},
		},
		{
			name: "should return error",
			args: args{
				ctx:     context.TODO(),
				user_id: 1,
				id:      1,
			},
			expected: expected{
				err:        sql.ErrNoRows,
				itemReturn: nil,
			},
			mocks: mocks{

				mockFunction: func() {
					mockItemRepo.On("Vote", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(nil, sql.ErrNoRows).Once()
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mocks.mockFunction()
			u := usecase.NewItemUsecase(mockItemRepo, mockUserRepo, time.Second*5)
			item, err := u.Vote(test.args.ctx, test.args.id, test.args.user_id)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expected.itemReturn, item)
				assert.NoError(t, err)
			}
			mockItemRepo.AssertExpectations(t)
		})
	}

}

func TestClearVote(t *testing.T) {
	mockItemRepo := new(mocks.ItemRepository)
	mockUserRepo := new(mocks.UserRepository)
	logger.InitializeLogger()
	mockItem := entity.Item{
		ID:          1,
		Name:        "The one",
		Description: "This is Description",
		VoteCount:   0,
	}

	type args struct {
		ctx     context.Context
		id      int64
		user_id int64
	}

	type mocks struct {
		mockFunction func()
	}

	type expected struct {
		itemReturn entity.Item
		err        error
	}

	tests := []struct {
		name     string
		args     args
		mocks    mocks
		expected expected
		wantErr  bool
	}{
		{
			name: "success",
			args: args{
				ctx:     context.TODO(),
				id:      1,
				user_id: 1,
			},
			expected: expected{
				itemReturn: mockItem,
			},
			mocks: mocks{
				mockFunction: func() {
					mockItemRepo.On("ClearVote", mock.Anything, mock.AnythingOfType("int64")).Return(mockItem, nil).Once()

				},
			},
		},
		{
			name: "should return error",
			args: args{
				ctx:     context.TODO(),
				user_id: 1,
				id:      1,
			},
			expected: expected{
				err:        sql.ErrNoRows,
				itemReturn: entity.Item{},
			},
			mocks: mocks{

				mockFunction: func() {
					mockItemRepo.On("ClearVote", mock.Anything, mock.AnythingOfType("int64")).Return(entity.Item{}, sql.ErrNoRows).Once()
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mocks.mockFunction()
			u := usecase.NewItemUsecase(mockItemRepo, mockUserRepo, time.Second*5)
			item, err := u.ClearVote(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expected.itemReturn, item)
				assert.NoError(t, err)
			}
			mockItemRepo.AssertExpectations(t)
		})
	}

}

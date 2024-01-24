// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "github.com/taroxii/vote-api/pkg/entity"
)

// ItemRepository is an autogenerated mock type for the ItemRepository type
type ItemRepository struct {
	mock.Mock
}

// ClearVote provides a mock function with given fields: ctx, id
func (_m *ItemRepository) ClearVote(ctx context.Context, id int64) (entity.Item, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for ClearVote")
	}

	var r0 entity.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (entity.Item, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) entity.Item); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(entity.Item)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *ItemRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx, cursor, num
func (_m *ItemRepository) Fetch(ctx context.Context, cursor string, num int64) ([]entity.Item, string, error) {
	ret := _m.Called(ctx, cursor, num)

	if len(ret) == 0 {
		panic("no return value specified for Fetch")
	}

	var r0 []entity.Item
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) ([]entity.Item, string, error)); ok {
		return rf(ctx, cursor, num)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) []entity.Item); ok {
		r0 = rf(ctx, cursor, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) string); ok {
		r1 = rf(ctx, cursor, num)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int64) error); ok {
		r2 = rf(ctx, cursor, num)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Insert provides a mock function with given fields: ctx, a
func (_m *ItemRepository) Insert(ctx context.Context, a *entity.Item) error {
	ret := _m.Called(ctx, a)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Item) error); ok {
		r0 = rf(ctx, a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, ar
func (_m *ItemRepository) Update(ctx context.Context, ar *entity.Item) error {
	ret := _m.Called(ctx, ar)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Item) error); ok {
		r0 = rf(ctx, ar)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Vote provides a mock function with given fields: ctx, id, user_id
func (_m *ItemRepository) Vote(ctx context.Context, id int64, user_id int64) (*entity.Item, error) {
	ret := _m.Called(ctx, id, user_id)

	if len(ret) == 0 {
		panic("no return value specified for Vote")
	}

	var r0 *entity.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (*entity.Item, error)); ok {
		return rf(ctx, id, user_id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) *entity.Item); ok {
		r0 = rf(ctx, id, user_id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, id, user_id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewItemRepository creates a new instance of ItemRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewItemRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ItemRepository {
	mock := &ItemRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

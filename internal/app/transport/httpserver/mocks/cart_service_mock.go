// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/cronnoss/bookshop-home-task/internal/app/domain"

	mock "github.com/stretchr/testify/mock"
)

// CartService is an autogenerated mock type for the CartService type
type CartService struct {
	mock.Mock
}

type CartService_Expecter struct {
	mock *mock.Mock
}

func (_m *CartService) EXPECT() *CartService_Expecter {
	return &CartService_Expecter{mock: &_m.Mock}
}

// Checkout provides a mock function with given fields: ctx, userID
func (_m *CartService) Checkout(ctx context.Context, userID int) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for Checkout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CartService_Checkout_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Checkout'
type CartService_Checkout_Call struct {
	*mock.Call
}

// Checkout is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int
func (_e *CartService_Expecter) Checkout(ctx interface{}, userID interface{}) *CartService_Checkout_Call {
	return &CartService_Checkout_Call{Call: _e.mock.On("Checkout", ctx, userID)}
}

func (_c *CartService_Checkout_Call) Run(run func(ctx context.Context, userID int)) *CartService_Checkout_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *CartService_Checkout_Call) Return(_a0 error) *CartService_Checkout_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CartService_Checkout_Call) RunAndReturn(run func(context.Context, int) error) *CartService_Checkout_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCartAndStocks provides a mock function with given fields: ctx, cart
func (_m *CartService) UpdateCartAndStocks(ctx context.Context, cart domain.Cart) (domain.Cart, error) {
	ret := _m.Called(ctx, cart)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCartAndStocks")
	}

	var r0 domain.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Cart) (domain.Cart, error)); ok {
		return rf(ctx, cart)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Cart) domain.Cart); ok {
		r0 = rf(ctx, cart)
	} else {
		r0 = ret.Get(0).(domain.Cart)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Cart) error); ok {
		r1 = rf(ctx, cart)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CartService_UpdateCartAndStocks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCartAndStocks'
type CartService_UpdateCartAndStocks_Call struct {
	*mock.Call
}

// UpdateCartAndStocks is a helper method to define mock.On call
//   - ctx context.Context
//   - cart domain.Cart
func (_e *CartService_Expecter) UpdateCartAndStocks(ctx interface{}, cart interface{}) *CartService_UpdateCartAndStocks_Call {
	return &CartService_UpdateCartAndStocks_Call{Call: _e.mock.On("UpdateCartAndStocks", ctx, cart)}
}

func (_c *CartService_UpdateCartAndStocks_Call) Run(run func(ctx context.Context, cart domain.Cart)) *CartService_UpdateCartAndStocks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Cart))
	})
	return _c
}

func (_c *CartService_UpdateCartAndStocks_Call) Return(_a0 domain.Cart, _a1 error) *CartService_UpdateCartAndStocks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CartService_UpdateCartAndStocks_Call) RunAndReturn(run func(context.Context, domain.Cart) (domain.Cart, error)) *CartService_UpdateCartAndStocks_Call {
	_c.Call.Return(run)
	return _c
}

// NewCartService creates a new instance of CartService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCartService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CartService {
	mock := &CartService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

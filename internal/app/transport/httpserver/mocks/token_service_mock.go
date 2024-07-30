// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/cronnoss/bookshop-home-task/internal/app/domain"

	mock "github.com/stretchr/testify/mock"
)

// TokenService is an autogenerated mock type for the TokenService type
type TokenService struct {
	mock.Mock
}

type TokenService_Expecter struct {
	mock *mock.Mock
}

func (_m *TokenService) EXPECT() *TokenService_Expecter {
	return &TokenService_Expecter{mock: &_m.Mock}
}

// GenerateToken provides a mock function with given fields: user
func (_m *TokenService) GenerateToken(user domain.User) (string, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for GenerateToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.User) (string, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(domain.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(domain.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenService_GenerateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateToken'
type TokenService_GenerateToken_Call struct {
	*mock.Call
}

// GenerateToken is a helper method to define mock.On call
//   - user domain.User
func (_e *TokenService_Expecter) GenerateToken(user interface{}) *TokenService_GenerateToken_Call {
	return &TokenService_GenerateToken_Call{Call: _e.mock.On("GenerateToken", user)}
}

func (_c *TokenService_GenerateToken_Call) Run(run func(user domain.User)) *TokenService_GenerateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.User))
	})
	return _c
}

func (_c *TokenService_GenerateToken_Call) Return(_a0 string, _a1 error) *TokenService_GenerateToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenService_GenerateToken_Call) RunAndReturn(run func(domain.User) (string, error)) *TokenService_GenerateToken_Call {
	_c.Call.Return(run)
	return _c
}

// GetUser provides a mock function with given fields: token
func (_m *TokenService) GetUser(token string) (domain.User, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.User, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) domain.User); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenService_GetUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUser'
type TokenService_GetUser_Call struct {
	*mock.Call
}

// GetUser is a helper method to define mock.On call
//   - token string
func (_e *TokenService_Expecter) GetUser(token interface{}) *TokenService_GetUser_Call {
	return &TokenService_GetUser_Call{Call: _e.mock.On("GetUser", token)}
}

func (_c *TokenService_GetUser_Call) Run(run func(token string)) *TokenService_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TokenService_GetUser_Call) Return(_a0 domain.User, _a1 error) *TokenService_GetUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenService_GetUser_Call) RunAndReturn(run func(string) (domain.User, error)) *TokenService_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewTokenService creates a new instance of TokenService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenService(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenService {
	mock := &TokenService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

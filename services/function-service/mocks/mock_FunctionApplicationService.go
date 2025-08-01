// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
	mock "github.com/stretchr/testify/mock"
)

// MockFunctionApplicationService is an autogenerated mock type for the FunctionApplicationService type
type MockFunctionApplicationService struct {
	mock.Mock
}

type MockFunctionApplicationService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFunctionApplicationService) EXPECT() *MockFunctionApplicationService_Expecter {
	return &MockFunctionApplicationService_Expecter{mock: &_m.Mock}
}

// PersistFunction provides a mock function with given fields: ctx, function
func (_m *MockFunctionApplicationService) PersistFunction(ctx context.Context, function *domain.Function) (domain.FunctionID, error) {
	ret := _m.Called(ctx, function)

	if len(ret) == 0 {
		panic("no return value specified for PersistFunction")
	}

	var r0 domain.FunctionID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Function) (domain.FunctionID, error)); ok {
		return rf(ctx, function)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Function) domain.FunctionID); ok {
		r0 = rf(ctx, function)
	} else {
		r0 = ret.Get(0).(domain.FunctionID)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Function) error); ok {
		r1 = rf(ctx, function)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFunctionApplicationService_PersistFunction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PersistFunction'
type MockFunctionApplicationService_PersistFunction_Call struct {
	*mock.Call
}

// PersistFunction is a helper method to define mock.On call
//   - ctx context.Context
//   - function *domain.Function
func (_e *MockFunctionApplicationService_Expecter) PersistFunction(ctx interface{}, function interface{}) *MockFunctionApplicationService_PersistFunction_Call {
	return &MockFunctionApplicationService_PersistFunction_Call{Call: _e.mock.On("PersistFunction", ctx, function)}
}

func (_c *MockFunctionApplicationService_PersistFunction_Call) Run(run func(ctx context.Context, function *domain.Function)) *MockFunctionApplicationService_PersistFunction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Function))
	})
	return _c
}

func (_c *MockFunctionApplicationService_PersistFunction_Call) Return(_a0 domain.FunctionID, _a1 error) *MockFunctionApplicationService_PersistFunction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFunctionApplicationService_PersistFunction_Call) RunAndReturn(run func(context.Context, *domain.Function) (domain.FunctionID, error)) *MockFunctionApplicationService_PersistFunction_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFunctionApplicationService creates a new instance of MockFunctionApplicationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFunctionApplicationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFunctionApplicationService {
	mock := &MockFunctionApplicationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

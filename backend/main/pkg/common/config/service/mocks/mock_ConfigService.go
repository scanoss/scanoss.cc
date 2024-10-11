// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	mock "github.com/stretchr/testify/mock"
)

// MockConfigService is an autogenerated mock type for the ConfigService type
type MockConfigService struct {
	mock.Mock
}

type MockConfigService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConfigService) EXPECT() *MockConfigService_Expecter {
	return &MockConfigService_Expecter{mock: &_m.Mock}
}

// ReadConfig provides a mock function with given fields:
func (_m *MockConfigService) ReadConfig() (entities.Config, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadConfig")
	}

	var r0 entities.Config
	var r1 error
	if rf, ok := ret.Get(0).(func() (entities.Config, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() entities.Config); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(entities.Config)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConfigService_ReadConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadConfig'
type MockConfigService_ReadConfig_Call struct {
	*mock.Call
}

// ReadConfig is a helper method to define mock.On call
func (_e *MockConfigService_Expecter) ReadConfig() *MockConfigService_ReadConfig_Call {
	return &MockConfigService_ReadConfig_Call{Call: _e.mock.On("ReadConfig")}
}

func (_c *MockConfigService_ReadConfig_Call) Run(run func()) *MockConfigService_ReadConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockConfigService_ReadConfig_Call) Return(_a0 entities.Config, _a1 error) *MockConfigService_ReadConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConfigService_ReadConfig_Call) RunAndReturn(run func() (entities.Config, error)) *MockConfigService_ReadConfig_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConfigService creates a new instance of MockConfigService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConfigService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConfigService {
	mock := &MockConfigService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

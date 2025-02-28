// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/scanoss/scanoss.cc/backend/entities"
	mock "github.com/stretchr/testify/mock"
)

// MockTreeService is an autogenerated mock type for the TreeService type
type MockTreeService struct {
	mock.Mock
}

type MockTreeService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTreeService) EXPECT() *MockTreeService_Expecter {
	return &MockTreeService_Expecter{mock: &_m.Mock}
}

// GetTree provides a mock function with given fields: rootPath
func (_m *MockTreeService) GetTree(rootPath string) ([]entities.TreeNode, error) {
	ret := _m.Called(rootPath)

	if len(ret) == 0 {
		panic("no return value specified for GetTree")
	}

	var r0 []entities.TreeNode
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]entities.TreeNode, error)); ok {
		return rf(rootPath)
	}
	if rf, ok := ret.Get(0).(func(string) []entities.TreeNode); ok {
		r0 = rf(rootPath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.TreeNode)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(rootPath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTreeService_GetTree_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTree'
type MockTreeService_GetTree_Call struct {
	*mock.Call
}

// GetTree is a helper method to define mock.On call
//   - rootPath string
func (_e *MockTreeService_Expecter) GetTree(rootPath interface{}) *MockTreeService_GetTree_Call {
	return &MockTreeService_GetTree_Call{Call: _e.mock.On("GetTree", rootPath)}
}

func (_c *MockTreeService_GetTree_Call) Run(run func(rootPath string)) *MockTreeService_GetTree_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockTreeService_GetTree_Call) Return(_a0 []entities.TreeNode, _a1 error) *MockTreeService_GetTree_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTreeService_GetTree_Call) RunAndReturn(run func(string) ([]entities.TreeNode, error)) *MockTreeService_GetTree_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTreeService creates a new instance of MockTreeService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTreeService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTreeService {
	mock := &MockTreeService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

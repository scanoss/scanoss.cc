// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"
	mock "github.com/stretchr/testify/mock"
)

// MockFileController is an autogenerated mock type for the FileController type
type MockFileController struct {
	mock.Mock
}

type MockFileController_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFileController) EXPECT() *MockFileController_Expecter {
	return &MockFileController_Expecter{mock: &_m.Mock}
}

// GetLocalFile provides a mock function with given fields: path
func (_m *MockFileController) GetLocalFile(path string) (entities.FileDTO, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for GetLocalFile")
	}

	var r0 entities.FileDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (entities.FileDTO, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) entities.FileDTO); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(entities.FileDTO)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFileController_GetLocalFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLocalFile'
type MockFileController_GetLocalFile_Call struct {
	*mock.Call
}

// GetLocalFile is a helper method to define mock.On call
//   - path string
func (_e *MockFileController_Expecter) GetLocalFile(path interface{}) *MockFileController_GetLocalFile_Call {
	return &MockFileController_GetLocalFile_Call{Call: _e.mock.On("GetLocalFile", path)}
}

func (_c *MockFileController_GetLocalFile_Call) Run(run func(path string)) *MockFileController_GetLocalFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFileController_GetLocalFile_Call) Return(_a0 entities.FileDTO, _a1 error) *MockFileController_GetLocalFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFileController_GetLocalFile_Call) RunAndReturn(run func(string) (entities.FileDTO, error)) *MockFileController_GetLocalFile_Call {
	_c.Call.Return(run)
	return _c
}

// GetRemoteFile provides a mock function with given fields: path
func (_m *MockFileController) GetRemoteFile(path string) (entities.FileDTO, error) {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for GetRemoteFile")
	}

	var r0 entities.FileDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (entities.FileDTO, error)); ok {
		return rf(path)
	}
	if rf, ok := ret.Get(0).(func(string) entities.FileDTO); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(entities.FileDTO)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFileController_GetRemoteFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRemoteFile'
type MockFileController_GetRemoteFile_Call struct {
	*mock.Call
}

// GetRemoteFile is a helper method to define mock.On call
//   - path string
func (_e *MockFileController_Expecter) GetRemoteFile(path interface{}) *MockFileController_GetRemoteFile_Call {
	return &MockFileController_GetRemoteFile_Call{Call: _e.mock.On("GetRemoteFile", path)}
}

func (_c *MockFileController_GetRemoteFile_Call) Run(run func(path string)) *MockFileController_GetRemoteFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFileController_GetRemoteFile_Call) Return(_a0 entities.FileDTO, _a1 error) *MockFileController_GetRemoteFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFileController_GetRemoteFile_Call) RunAndReturn(run func(string) (entities.FileDTO, error)) *MockFileController_GetRemoteFile_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFileController creates a new instance of MockFileController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFileController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFileController {
	mock := &MockFileController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

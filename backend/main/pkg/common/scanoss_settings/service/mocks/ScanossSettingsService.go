// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ScanossSettingsService is an autogenerated mock type for the ScanossSettingsService type
type ScanossSettingsService struct {
	mock.Mock
}

type ScanossSettingsService_Expecter struct {
	mock *mock.Mock
}

func (_m *ScanossSettingsService) EXPECT() *ScanossSettingsService_Expecter {
	return &ScanossSettingsService_Expecter{mock: &_m.Mock}
}

// HasUnsavedChanges provides a mock function with given fields:
func (_m *ScanossSettingsService) HasUnsavedChanges() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HasUnsavedChanges")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScanossSettingsService_HasUnsavedChanges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasUnsavedChanges'
type ScanossSettingsService_HasUnsavedChanges_Call struct {
	*mock.Call
}

// HasUnsavedChanges is a helper method to define mock.On call
func (_e *ScanossSettingsService_Expecter) HasUnsavedChanges() *ScanossSettingsService_HasUnsavedChanges_Call {
	return &ScanossSettingsService_HasUnsavedChanges_Call{Call: _e.mock.On("HasUnsavedChanges")}
}

func (_c *ScanossSettingsService_HasUnsavedChanges_Call) Run(run func()) *ScanossSettingsService_HasUnsavedChanges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ScanossSettingsService_HasUnsavedChanges_Call) Return(_a0 bool, _a1 error) *ScanossSettingsService_HasUnsavedChanges_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ScanossSettingsService_HasUnsavedChanges_Call) RunAndReturn(run func() (bool, error)) *ScanossSettingsService_HasUnsavedChanges_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields:
func (_m *ScanossSettingsService) Save() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ScanossSettingsService_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type ScanossSettingsService_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
func (_e *ScanossSettingsService_Expecter) Save() *ScanossSettingsService_Save_Call {
	return &ScanossSettingsService_Save_Call{Call: _e.mock.On("Save")}
}

func (_c *ScanossSettingsService_Save_Call) Run(run func()) *ScanossSettingsService_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ScanossSettingsService_Save_Call) Return(_a0 error) *ScanossSettingsService_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ScanossSettingsService_Save_Call) RunAndReturn(run func() error) *ScanossSettingsService_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewScanossSettingsService creates a new instance of ScanossSettingsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScanossSettingsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScanossSettingsService {
	mock := &ScanossSettingsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
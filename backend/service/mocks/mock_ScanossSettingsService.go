// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockScanossSettingsService is an autogenerated mock type for the ScanossSettingsService type
type MockScanossSettingsService struct {
	mock.Mock
}

type MockScanossSettingsService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockScanossSettingsService) EXPECT() *MockScanossSettingsService_Expecter {
	return &MockScanossSettingsService_Expecter{mock: &_m.Mock}
}

// CommitStagedSkipPatterns provides a mock function with given fields:
func (_m *MockScanossSettingsService) CommitStagedSkipPatterns() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CommitStagedSkipPatterns")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScanossSettingsService_CommitStagedSkipPatterns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CommitStagedSkipPatterns'
type MockScanossSettingsService_CommitStagedSkipPatterns_Call struct {
	*mock.Call
}

// CommitStagedSkipPatterns is a helper method to define mock.On call
func (_e *MockScanossSettingsService_Expecter) CommitStagedSkipPatterns() *MockScanossSettingsService_CommitStagedSkipPatterns_Call {
	return &MockScanossSettingsService_CommitStagedSkipPatterns_Call{Call: _e.mock.On("CommitStagedSkipPatterns")}
}

func (_c *MockScanossSettingsService_CommitStagedSkipPatterns_Call) Run(run func()) *MockScanossSettingsService_CommitStagedSkipPatterns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockScanossSettingsService_CommitStagedSkipPatterns_Call) Return(_a0 error) *MockScanossSettingsService_CommitStagedSkipPatterns_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_CommitStagedSkipPatterns_Call) RunAndReturn(run func() error) *MockScanossSettingsService_CommitStagedSkipPatterns_Call {
	_c.Call.Return(run)
	return _c
}

// DiscardStagedSkipPatterns provides a mock function with given fields:
func (_m *MockScanossSettingsService) DiscardStagedSkipPatterns() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DiscardStagedSkipPatterns")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScanossSettingsService_DiscardStagedSkipPatterns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DiscardStagedSkipPatterns'
type MockScanossSettingsService_DiscardStagedSkipPatterns_Call struct {
	*mock.Call
}

// DiscardStagedSkipPatterns is a helper method to define mock.On call
func (_e *MockScanossSettingsService_Expecter) DiscardStagedSkipPatterns() *MockScanossSettingsService_DiscardStagedSkipPatterns_Call {
	return &MockScanossSettingsService_DiscardStagedSkipPatterns_Call{Call: _e.mock.On("DiscardStagedSkipPatterns")}
}

func (_c *MockScanossSettingsService_DiscardStagedSkipPatterns_Call) Run(run func()) *MockScanossSettingsService_DiscardStagedSkipPatterns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockScanossSettingsService_DiscardStagedSkipPatterns_Call) Return(_a0 error) *MockScanossSettingsService_DiscardStagedSkipPatterns_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_DiscardStagedSkipPatterns_Call) RunAndReturn(run func() error) *MockScanossSettingsService_DiscardStagedSkipPatterns_Call {
	_c.Call.Return(run)
	return _c
}

// HasUnsavedChanges provides a mock function with given fields:
func (_m *MockScanossSettingsService) HasUnsavedChanges() (bool, error) {
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

// MockScanossSettingsService_HasUnsavedChanges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasUnsavedChanges'
type MockScanossSettingsService_HasUnsavedChanges_Call struct {
	*mock.Call
}

// HasUnsavedChanges is a helper method to define mock.On call
func (_e *MockScanossSettingsService_Expecter) HasUnsavedChanges() *MockScanossSettingsService_HasUnsavedChanges_Call {
	return &MockScanossSettingsService_HasUnsavedChanges_Call{Call: _e.mock.On("HasUnsavedChanges")}
}

func (_c *MockScanossSettingsService_HasUnsavedChanges_Call) Run(run func()) *MockScanossSettingsService_HasUnsavedChanges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockScanossSettingsService_HasUnsavedChanges_Call) Return(_a0 bool, _a1 error) *MockScanossSettingsService_HasUnsavedChanges_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockScanossSettingsService_HasUnsavedChanges_Call) RunAndReturn(run func() (bool, error)) *MockScanossSettingsService_HasUnsavedChanges_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields:
func (_m *MockScanossSettingsService) Save() error {
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

// MockScanossSettingsService_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockScanossSettingsService_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
func (_e *MockScanossSettingsService_Expecter) Save() *MockScanossSettingsService_Save_Call {
	return &MockScanossSettingsService_Save_Call{Call: _e.mock.On("Save")}
}

func (_c *MockScanossSettingsService_Save_Call) Run(run func()) *MockScanossSettingsService_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockScanossSettingsService_Save_Call) Return(_a0 error) *MockScanossSettingsService_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_Save_Call) RunAndReturn(run func() error) *MockScanossSettingsService_Save_Call {
	_c.Call.Return(run)
	return _c
}

// ToggleScanningSkipPattern provides a mock function with given fields: pattern
func (_m *MockScanossSettingsService) ToggleScanningSkipPattern(pattern string) error {
	ret := _m.Called(pattern)

	if len(ret) == 0 {
		panic("no return value specified for ToggleScanningSkipPattern")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(pattern)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScanossSettingsService_ToggleScanningSkipPattern_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ToggleScanningSkipPattern'
type MockScanossSettingsService_ToggleScanningSkipPattern_Call struct {
	*mock.Call
}

// ToggleScanningSkipPattern is a helper method to define mock.On call
//   - pattern string
func (_e *MockScanossSettingsService_Expecter) ToggleScanningSkipPattern(pattern interface{}) *MockScanossSettingsService_ToggleScanningSkipPattern_Call {
	return &MockScanossSettingsService_ToggleScanningSkipPattern_Call{Call: _e.mock.On("ToggleScanningSkipPattern", pattern)}
}

func (_c *MockScanossSettingsService_ToggleScanningSkipPattern_Call) Run(run func(pattern string)) *MockScanossSettingsService_ToggleScanningSkipPattern_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockScanossSettingsService_ToggleScanningSkipPattern_Call) Return(_a0 error) *MockScanossSettingsService_ToggleScanningSkipPattern_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_ToggleScanningSkipPattern_Call) RunAndReturn(run func(string) error) *MockScanossSettingsService_ToggleScanningSkipPattern_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockScanossSettingsService creates a new instance of MockScanossSettingsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockScanossSettingsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockScanossSettingsService {
	mock := &MockScanossSettingsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

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

// AddStagedScanningSkipPattern provides a mock function with given fields: pattern
func (_m *MockScanossSettingsService) AddStagedScanningSkipPattern(pattern string) error {
	ret := _m.Called(pattern)

	if len(ret) == 0 {
		panic("no return value specified for AddStagedScanningSkipPattern")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(pattern)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScanossSettingsService_AddStagedScanningSkipPattern_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddStagedScanningSkipPattern'
type MockScanossSettingsService_AddStagedScanningSkipPattern_Call struct {
	*mock.Call
}

// AddStagedScanningSkipPattern is a helper method to define mock.On call
//   - pattern string
func (_e *MockScanossSettingsService_Expecter) AddStagedScanningSkipPattern(pattern interface{}) *MockScanossSettingsService_AddStagedScanningSkipPattern_Call {
	return &MockScanossSettingsService_AddStagedScanningSkipPattern_Call{Call: _e.mock.On("AddStagedScanningSkipPattern", pattern)}
}

func (_c *MockScanossSettingsService_AddStagedScanningSkipPattern_Call) Run(run func(pattern string)) *MockScanossSettingsService_AddStagedScanningSkipPattern_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockScanossSettingsService_AddStagedScanningSkipPattern_Call) Return(_a0 error) *MockScanossSettingsService_AddStagedScanningSkipPattern_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_AddStagedScanningSkipPattern_Call) RunAndReturn(run func(string) error) *MockScanossSettingsService_AddStagedScanningSkipPattern_Call {
	_c.Call.Return(run)
	return _c
}

// CommitStagedScanningSkipPatterns provides a mock function with given fields:
func (_m *MockScanossSettingsService) CommitStagedScanningSkipPatterns() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CommitStagedScanningSkipPatterns")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CommitStagedScanningSkipPatterns'
type MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call struct {
	*mock.Call
}

// CommitStagedScanningSkipPatterns is a helper method to define mock.On call
func (_e *MockScanossSettingsService_Expecter) CommitStagedScanningSkipPatterns() *MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call {
	return &MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call{Call: _e.mock.On("CommitStagedScanningSkipPatterns")}
}

func (_c *MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call) Run(run func()) *MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call) Return(_a0 error) *MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call) RunAndReturn(run func() error) *MockScanossSettingsService_CommitStagedScanningSkipPatterns_Call {
	_c.Call.Return(run)
	return _c
}

// DiscardStagedScanningSkipPatterns provides a mock function with given fields:
func (_m *MockScanossSettingsService) DiscardStagedScanningSkipPatterns() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DiscardStagedScanningSkipPatterns")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DiscardStagedScanningSkipPatterns'
type MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call struct {
	*mock.Call
}

// DiscardStagedScanningSkipPatterns is a helper method to define mock.On call
func (_e *MockScanossSettingsService_Expecter) DiscardStagedScanningSkipPatterns() *MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call {
	return &MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call{Call: _e.mock.On("DiscardStagedScanningSkipPatterns")}
}

func (_c *MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call) Run(run func()) *MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call) Return(_a0 error) *MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call) RunAndReturn(run func() error) *MockScanossSettingsService_DiscardStagedScanningSkipPatterns_Call {
	_c.Call.Return(run)
	return _c
}

// HasStagedScanningSkipPatternChanges provides a mock function with given fields:
func (_m *MockScanossSettingsService) HasStagedScanningSkipPatternChanges() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HasStagedScanningSkipPatternChanges")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasStagedScanningSkipPatternChanges'
type MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call struct {
	*mock.Call
}

// HasStagedScanningSkipPatternChanges is a helper method to define mock.On call
func (_e *MockScanossSettingsService_Expecter) HasStagedScanningSkipPatternChanges() *MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call {
	return &MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call{Call: _e.mock.On("HasStagedScanningSkipPatternChanges")}
}

func (_c *MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call) Run(run func()) *MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call) Return(_a0 bool) *MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call) RunAndReturn(run func() bool) *MockScanossSettingsService_HasStagedScanningSkipPatternChanges_Call {
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

// RemoveStagedScanningSkipPattern provides a mock function with given fields: path, pattern
func (_m *MockScanossSettingsService) RemoveStagedScanningSkipPattern(path string, pattern string) error {
	ret := _m.Called(path, pattern)

	if len(ret) == 0 {
		panic("no return value specified for RemoveStagedScanningSkipPattern")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(path, pattern)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveStagedScanningSkipPattern'
type MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call struct {
	*mock.Call
}

// RemoveStagedScanningSkipPattern is a helper method to define mock.On call
//   - path string
//   - pattern string
func (_e *MockScanossSettingsService_Expecter) RemoveStagedScanningSkipPattern(path interface{}, pattern interface{}) *MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call {
	return &MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call{Call: _e.mock.On("RemoveStagedScanningSkipPattern", path, pattern)}
}

func (_c *MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call) Run(run func(path string, pattern string)) *MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call) Return(_a0 error) *MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call) RunAndReturn(run func(string, string) error) *MockScanossSettingsService_RemoveStagedScanningSkipPattern_Call {
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

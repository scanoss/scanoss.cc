// Code generated by mockery v2.46.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	mock "github.com/stretchr/testify/mock"
)

// ResultController is an autogenerated mock type for the ResultController type
type ResultController struct {
	mock.Mock
}

type ResultController_Expecter struct {
	mock *mock.Mock
}

func (_m *ResultController) EXPECT() *ResultController_Expecter {
	return &ResultController_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: dto
func (_m *ResultController) GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	ret := _m.Called(dto)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []entities.ResultDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(*entities.RequestResultDTO) ([]entities.ResultDTO, error)); ok {
		return rf(dto)
	}
	if rf, ok := ret.Get(0).(func(*entities.RequestResultDTO) []entities.ResultDTO); ok {
		r0 = rf(dto)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.ResultDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(*entities.RequestResultDTO) error); ok {
		r1 = rf(dto)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResultController_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type ResultController_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - dto *entities.RequestResultDTO
func (_e *ResultController_Expecter) GetAll(dto interface{}) *ResultController_GetAll_Call {
	return &ResultController_GetAll_Call{Call: _e.mock.On("GetAll", dto)}
}

func (_c *ResultController_GetAll_Call) Run(run func(dto *entities.RequestResultDTO)) *ResultController_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*entities.RequestResultDTO))
	})
	return _c
}

func (_c *ResultController_GetAll_Call) Return(_a0 []entities.ResultDTO, _a1 error) *ResultController_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResultController_GetAll_Call) RunAndReturn(run func(*entities.RequestResultDTO) ([]entities.ResultDTO, error)) *ResultController_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewResultController creates a new instance of ResultController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResultController(t interface {
	mock.TestingT
	Cleanup(func())
}) *ResultController {
	mock := &ResultController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

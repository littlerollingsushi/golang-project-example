// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Timer is an autogenerated mock type for the Timer type
type Timer struct {
	mock.Mock
}

// NowInUTC provides a mock function with given fields:
func (_m *Timer) NowInUTC() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

type mockConstructorTestingTNewTimer interface {
	mock.TestingT
	Cleanup(func())
}

// NewTimer creates a new instance of Timer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTimer(t mockConstructorTestingTNewTimer) *Timer {
	mock := &Timer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

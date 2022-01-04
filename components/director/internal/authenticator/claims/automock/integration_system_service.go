// Code generated by mockery 2.9.0. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IntegrationSystemService is an autogenerated mock type for the IntegrationSystemService type
type IntegrationSystemService struct {
	mock.Mock
}

// Exists provides a mock function with given fields: _a0, _a1
func (_m *IntegrationSystemService) Exists(_a0 context.Context, _a1 string) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

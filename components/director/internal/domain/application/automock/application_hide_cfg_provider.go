// Code generated by mockery. DO NOT EDIT.

package automock

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// ApplicationHideCfgProvider is an autogenerated mock type for the ApplicationHideCfgProvider type
type ApplicationHideCfgProvider struct {
	mock.Mock
}

// GetApplicationHideSelectors provides a mock function with given fields:
func (_m *ApplicationHideCfgProvider) GetApplicationHideSelectors() (map[string][]string, error) {
	ret := _m.Called()

	var r0 map[string][]string
	if rf, ok := ret.Get(0).(func() map[string][]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewApplicationHideCfgProvider creates a new instance of ApplicationHideCfgProvider. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewApplicationHideCfgProvider(t testing.TB) *ApplicationHideCfgProvider {
	mock := &ApplicationHideCfgProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

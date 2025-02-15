// Code generated by mockery. DO NOT EDIT.

package automock

import (
	accessstrategy "github.com/kyma-incubator/compass/components/director/pkg/accessstrategy"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ExecutorProvider is an autogenerated mock type for the ExecutorProvider type
type ExecutorProvider struct {
	mock.Mock
}

// Provide provides a mock function with given fields: accessStrategyType
func (_m *ExecutorProvider) Provide(accessStrategyType accessstrategy.Type) (accessstrategy.Executor, error) {
	ret := _m.Called(accessStrategyType)

	var r0 accessstrategy.Executor
	if rf, ok := ret.Get(0).(func(accessstrategy.Type) accessstrategy.Executor); ok {
		r0 = rf(accessStrategyType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(accessstrategy.Executor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(accessstrategy.Type) error); ok {
		r1 = rf(accessStrategyType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExecutorProvider creates a new instance of ExecutorProvider. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewExecutorProvider(t testing.TB) *ExecutorProvider {
	mock := &ExecutorProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

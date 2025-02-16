// Code generated by mockery. DO NOT EDIT.

package automock

import (
	runtime "github.com/kyma-incubator/compass/components/director/internal/domain/runtime"
	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *EntityConverter) FromEntity(entity *runtime.Runtime) *model.Runtime {
	ret := _m.Called(entity)

	var r0 *model.Runtime
	if rf, ok := ret.Get(0).(func(*runtime.Runtime) *model.Runtime); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Runtime)
		}
	}

	return r0
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in *model.Runtime) (*runtime.Runtime, error) {
	ret := _m.Called(in)

	var r0 *runtime.Runtime
	if rf, ok := ret.Get(0).(func(*model.Runtime) *runtime.Runtime); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtime.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Runtime) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEntityConverter creates a new instance of EntityConverter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewEntityConverter(t testing.TB) *EntityConverter {
	mock := &EntityConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

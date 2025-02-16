// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ExternalTenantsService is an autogenerated mock type for the ExternalTenantsService type
type ExternalTenantsService struct {
	mock.Mock
}

// GetTenantByID provides a mock function with given fields: ctx, id
func (_m *ExternalTenantsService) GetTenantByID(ctx context.Context, id string) (*model.BusinessTenantMapping, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.BusinessTenantMapping
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.BusinessTenantMapping); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.BusinessTenantMapping)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExternalTenantsService creates a new instance of ExternalTenantsService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewExternalTenantsService(t testing.TB) *ExternalTenantsService {
	mock := &ExternalTenantsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

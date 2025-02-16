// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// FetchRequestRepository is an autogenerated mock type for the FetchRequestRepository type
type FetchRequestRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, tenant, item
func (_m *FetchRequestRepository) Create(ctx context.Context, tenant string, item *model.FetchRequest) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.FetchRequest) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByReferenceObjectID provides a mock function with given fields: ctx, tenant, objectType, objectID
func (_m *FetchRequestRepository) DeleteByReferenceObjectID(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectID string) error {
	ret := _m.Called(ctx, tenant, objectType, objectID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FetchRequestReferenceObjectType, string) error); ok {
		r0 = rf(ctx, tenant, objectType, objectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByReferenceObjectID provides a mock function with given fields: ctx, tenant, objectType, objectID
func (_m *FetchRequestRepository) GetByReferenceObjectID(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectID string) (*model.FetchRequest, error) {
	ret := _m.Called(ctx, tenant, objectType, objectID)

	var r0 *model.FetchRequest
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FetchRequestReferenceObjectType, string) *model.FetchRequest); ok {
		r0 = rf(ctx, tenant, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FetchRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.FetchRequestReferenceObjectType, string) error); ok {
		r1 = rf(ctx, tenant, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByReferenceObjectIDs provides a mock function with given fields: ctx, tenant, objectType, objectIDs
func (_m *FetchRequestRepository) ListByReferenceObjectIDs(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectIDs []string) ([]*model.FetchRequest, error) {
	ret := _m.Called(ctx, tenant, objectType, objectIDs)

	var r0 []*model.FetchRequest
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FetchRequestReferenceObjectType, []string) []*model.FetchRequest); ok {
		r0 = rf(ctx, tenant, objectType, objectIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FetchRequest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.FetchRequestReferenceObjectType, []string) error); ok {
		r1 = rf(ctx, tenant, objectType, objectIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFetchRequestRepository creates a new instance of FetchRequestRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFetchRequestRepository(t testing.TB) *FetchRequestRepository {
	mock := &FetchRequestRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

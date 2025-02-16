// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// BundleReferenceRepository is an autogenerated mock type for the BundleReferenceRepository type
type BundleReferenceRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *BundleReferenceRepository) Create(ctx context.Context, item *model.BundleReference) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.BundleReference) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByReferenceObjectID provides a mock function with given fields: ctx, bundleID, objectType, objectID
func (_m *BundleReferenceRepository) DeleteByReferenceObjectID(ctx context.Context, bundleID string, objectType model.BundleReferenceObjectType, objectID string) error {
	ret := _m.Called(ctx, bundleID, objectType, objectID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.BundleReferenceObjectType, string) error); ok {
		r0 = rf(ctx, bundleID, objectType, objectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBundleIDsForObject provides a mock function with given fields: ctx, objectType, objectID
func (_m *BundleReferenceRepository) GetBundleIDsForObject(ctx context.Context, objectType model.BundleReferenceObjectType, objectID *string) ([]string, error) {
	ret := _m.Called(ctx, objectType, objectID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, model.BundleReferenceObjectType, *string) []string); ok {
		r0 = rf(ctx, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.BundleReferenceObjectType, *string) error); ok {
		r1 = rf(ctx, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, objectType, objectID, bundleID
func (_m *BundleReferenceRepository) GetByID(ctx context.Context, objectType model.BundleReferenceObjectType, objectID *string, bundleID *string) (*model.BundleReference, error) {
	ret := _m.Called(ctx, objectType, objectID, bundleID)

	var r0 *model.BundleReference
	if rf, ok := ret.Get(0).(func(context.Context, model.BundleReferenceObjectType, *string, *string) *model.BundleReference); ok {
		r0 = rf(ctx, objectType, objectID, bundleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.BundleReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.BundleReferenceObjectType, *string, *string) error); ok {
		r1 = rf(ctx, objectType, objectID, bundleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByBundleIDs provides a mock function with given fields: ctx, objectType, bundleIDs, pageSize, cursor
func (_m *BundleReferenceRepository) ListByBundleIDs(ctx context.Context, objectType model.BundleReferenceObjectType, bundleIDs []string, pageSize int, cursor string) ([]*model.BundleReference, map[string]int, error) {
	ret := _m.Called(ctx, objectType, bundleIDs, pageSize, cursor)

	var r0 []*model.BundleReference
	if rf, ok := ret.Get(0).(func(context.Context, model.BundleReferenceObjectType, []string, int, string) []*model.BundleReference); ok {
		r0 = rf(ctx, objectType, bundleIDs, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.BundleReference)
		}
	}

	var r1 map[string]int
	if rf, ok := ret.Get(1).(func(context.Context, model.BundleReferenceObjectType, []string, int, string) map[string]int); ok {
		r1 = rf(ctx, objectType, bundleIDs, pageSize, cursor)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]int)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, model.BundleReferenceObjectType, []string, int, string) error); ok {
		r2 = rf(ctx, objectType, bundleIDs, pageSize, cursor)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Update provides a mock function with given fields: ctx, item
func (_m *BundleReferenceRepository) Update(ctx context.Context, item *model.BundleReference) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.BundleReference) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBundleReferenceRepository creates a new instance of BundleReferenceRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewBundleReferenceRepository(t testing.TB) *BundleReferenceRepository {
	mock := &BundleReferenceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

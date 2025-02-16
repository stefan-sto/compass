// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// LabelRepository is an autogenerated mock type for the LabelRepository type
type LabelRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, tenant, objectType, objectID, key
func (_m *LabelRepository) Delete(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) error {
	ret := _m.Called(ctx, tenant, objectType, objectID, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string, string) error); ok {
		r0 = rf(ctx, tenant, objectType, objectID, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByKey provides a mock function with given fields: ctx, tenant, key
func (_m *LabelRepository) DeleteByKey(ctx context.Context, tenant string, key string) error {
	ret := _m.Called(ctx, tenant, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenant, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByKey provides a mock function with given fields: ctx, tenant, objectType, objectID, key
func (_m *LabelRepository) GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) (*model.Label, error) {
	ret := _m.Called(ctx, tenant, objectType, objectID, key)

	var r0 *model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string, string) *model.Label); ok {
		r0 = rf(ctx, tenant, objectType, objectID, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.LabelableObject, string, string) error); ok {
		r1 = rf(ctx, tenant, objectType, objectID, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upsert provides a mock function with given fields: ctx, tenant, label
func (_m *LabelRepository) Upsert(ctx context.Context, tenant string, label *model.Label) error {
	ret := _m.Called(ctx, tenant, label)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Label) error); ok {
		r0 = rf(ctx, tenant, label)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLabelRepository creates a new instance of LabelRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewLabelRepository(t testing.TB) *LabelRepository {
	mock := &LabelRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	labelfilter "github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// RuntimeRepository is an autogenerated mock type for the RuntimeRepository type
type RuntimeRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *RuntimeRepository) Create(ctx context.Context, item *model.Runtime) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Runtime) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, tenant, id
func (_m *RuntimeRepository) Delete(ctx context.Context, tenant string, id string) error {
	ret := _m.Called(ctx, tenant, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, tenant, id
func (_m *RuntimeRepository) Exists(ctx context.Context, tenant string, id string) (bool, error) {
	ret := _m.Called(ctx, tenant, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByFiltersGlobal provides a mock function with given fields: ctx, filter
func (_m *RuntimeRepository) GetByFiltersGlobal(ctx context.Context, filter []*labelfilter.LabelFilter) (*model.Runtime, error) {
	ret := _m.Called(ctx, filter)

	var r0 *model.Runtime
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter) *model.Runtime); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, tenant, id
func (_m *RuntimeRepository) GetByID(ctx context.Context, tenant string, id string) (*model.Runtime, error) {
	ret := _m.Called(ctx, tenant, id)

	var r0 *model.Runtime
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Runtime); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, tenant, filter, pageSize, cursor
func (_m *RuntimeRepository) List(ctx context.Context, tenant string, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimePage, error) {
	ret := _m.Called(ctx, tenant, filter, pageSize, cursor)

	var r0 *model.RuntimePage
	if rf, ok := ret.Get(0).(func(context.Context, string, []*labelfilter.LabelFilter, int, string) *model.RuntimePage); ok {
		r0 = rf(ctx, tenant, filter, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RuntimePage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []*labelfilter.LabelFilter, int, string) error); ok {
		r1 = rf(ctx, tenant, filter, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByFiltersGlobal provides a mock function with given fields: _a0, _a1
func (_m *RuntimeRepository) ListByFiltersGlobal(_a0 context.Context, _a1 []*labelfilter.LabelFilter) ([]*model.Runtime, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*model.Runtime
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter) []*model.Runtime); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Runtime)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, item
func (_m *RuntimeRepository) Update(ctx context.Context, item *model.Runtime) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Runtime) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTenantID provides a mock function with given fields: ctx, runtimeID, newTenantID
func (_m *RuntimeRepository) UpdateTenantID(ctx context.Context, runtimeID string, newTenantID string) error {
	ret := _m.Called(ctx, runtimeID, newTenantID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, runtimeID, newTenantID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

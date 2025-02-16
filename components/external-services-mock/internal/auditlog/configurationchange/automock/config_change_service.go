// Code generated by mockery. DO NOT EDIT.

package automock

import (
	testing "testing"

	model "github.com/kyma-incubator/compass/components/gateway/pkg/auditlog/model"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ConfigChangeService is an autogenerated mock type for the ConfigChangeService type
type ConfigChangeService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *ConfigChangeService) Delete(id string) {
	_m.Called(id)
}

// Get provides a mock function with given fields: id
func (_m *ConfigChangeService) Get(id string) *model.ConfigurationChange {
	ret := _m.Called(id)

	var r0 *model.ConfigurationChange
	if rf, ok := ret.Get(0).(func(string) *model.ConfigurationChange); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ConfigurationChange)
		}
	}

	return r0
}

// List provides a mock function with given fields:
func (_m *ConfigChangeService) List() []model.ConfigurationChange {
	ret := _m.Called()

	var r0 []model.ConfigurationChange
	if rf, ok := ret.Get(0).(func() []model.ConfigurationChange); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ConfigurationChange)
		}
	}

	return r0
}

// Save provides a mock function with given fields: change
func (_m *ConfigChangeService) Save(change model.ConfigurationChange) (string, error) {
	ret := _m.Called(change)

	var r0 string
	if rf, ok := ret.Get(0).(func(model.ConfigurationChange) string); ok {
		r0 = rf(change)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.ConfigurationChange) error); ok {
		r1 = rf(change)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchByTimestamp provides a mock function with given fields: timeFrom, timeTo
func (_m *ConfigChangeService) SearchByTimestamp(timeFrom time.Time, timeTo time.Time) []model.ConfigurationChange {
	ret := _m.Called(timeFrom, timeTo)

	var r0 []model.ConfigurationChange
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) []model.ConfigurationChange); ok {
		r0 = rf(timeFrom, timeTo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ConfigurationChange)
		}
	}

	return r0
}

// NewConfigChangeService creates a new instance of ConfigChangeService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfigChangeService(t testing.TB) *ConfigChangeService {
	mock := &ConfigChangeService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

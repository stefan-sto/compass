// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	systemfetcher "github.com/kyma-incubator/compass/components/director/internal/systemfetcher"

	testing "testing"
)

// TemplateRenderer is an autogenerated mock type for the templateRenderer type
type TemplateRenderer struct {
	mock.Mock
}

// ApplicationRegisterInputFromTemplate provides a mock function with given fields: ctx, sc
func (_m *TemplateRenderer) ApplicationRegisterInputFromTemplate(ctx context.Context, sc systemfetcher.System) (*model.ApplicationRegisterInput, error) {
	ret := _m.Called(ctx, sc)

	var r0 *model.ApplicationRegisterInput
	if rf, ok := ret.Get(0).(func(context.Context, systemfetcher.System) *model.ApplicationRegisterInput); ok {
		r0 = rf(ctx, sc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationRegisterInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, systemfetcher.System) error); ok {
		r1 = rf(ctx, sc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTemplateRenderer creates a new instance of TemplateRenderer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTemplateRenderer(t testing.TB) *TemplateRenderer {
	mock := &TemplateRenderer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

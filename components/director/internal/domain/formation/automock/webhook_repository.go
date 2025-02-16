// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// WebhookRepository is an autogenerated mock type for the webhookRepository type
type WebhookRepository struct {
	mock.Mock
}

// GetByIDAndWebhookType provides a mock function with given fields: ctx, tenant, objectID, objectType, webhookType
func (_m *WebhookRepository) GetByIDAndWebhookType(ctx context.Context, tenant string, objectID string, objectType model.WebhookReferenceObjectType, webhookType model.WebhookType) (*model.Webhook, error) {
	ret := _m.Called(ctx, tenant, objectID, objectType, webhookType)

	var r0 *model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string, string, model.WebhookReferenceObjectType, model.WebhookType) *model.Webhook); ok {
		r0 = rf(ctx, tenant, objectID, objectType, webhookType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, model.WebhookReferenceObjectType, model.WebhookType) error); ok {
		r1 = rf(ctx, tenant, objectID, objectType, webhookType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByReferenceObjectTypeAndWebhookType provides a mock function with given fields: ctx, tenant, whType, objType
func (_m *WebhookRepository) ListByReferenceObjectTypeAndWebhookType(ctx context.Context, tenant string, whType model.WebhookType, objType model.WebhookReferenceObjectType) ([]*model.Webhook, error) {
	ret := _m.Called(ctx, tenant, whType, objType)

	var r0 []*model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string, model.WebhookType, model.WebhookReferenceObjectType) []*model.Webhook); ok {
		r0 = rf(ctx, tenant, whType, objType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.WebhookType, model.WebhookReferenceObjectType) error); ok {
		r1 = rf(ctx, tenant, whType, objType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewWebhookRepositoryT interface {
	mock.TestingT
	Cleanup(func())
}

// NewWebhookRepository creates a new instance of WebhookRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebhookRepository(t NewWebhookRepositoryT) *WebhookRepository {
	mock := &WebhookRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

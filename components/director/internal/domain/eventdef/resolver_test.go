package eventdef_test

import (
	"context"
	"testing"
	"time"

	dataloader "github.com/kyma-incubator/compass/components/director/internal/dataloaders"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	event "github.com/kyma-incubator/compass/components/director/internal/domain/eventdef"
	"github.com/kyma-incubator/compass/components/director/internal/domain/eventdef/automock"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	persistenceautomock "github.com/kyma-incubator/compass/components/director/pkg/persistence/automock"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence/txtest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolver_AddEventToBundle(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "bar"

	modelEvent, spec, bundleRef := fixFullEventDefinitionModel("test")
	modelBndl := &model.Bundle{
		ApplicationID: appID,
		BaseEntity: &model.BaseEntity{
			ID: bundleID,
		},
	}
	gqlEvent := fixFullGQLEventDefinition("test")
	gqlEventInput := fixGQLEventDefinitionInput("name", "foo", "bar")
	modelEventInput, specInput := fixModelEventDefinitionInput("name", "foo", "bar")

	txGen := txtest.NewTransactionContextGenerator(testErr)

	testCases := []struct {
		Name                   string
		TransactionerFn        func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		ServiceFn              func() *automock.EventDefService
		BndlServiceFn          func() *automock.BundleService
		BndlReferenceServiceFn func() *automock.BundleReferenceService
		SpecServiceFn          func() *automock.SpecService
		ConverterFn            func() *automock.EventDefConverter
		ExpectedEvent          *graphql.EventDefinition
		ExpectedErr            error
	}{
		{
			Name:            "Success",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("CreateInBundle", txtest.CtxWithDBMatcher(), appID, bundleID, *modelEventInput, specInput).Return(id, nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEvent, nil).Once()
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(modelBndl, nil)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &gqlEvent.ID, str.Ptr(bundleID)).Return(&bundleRef, nil).Once()
				return svc
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, gqlEvent.ID).Return(&spec, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				conv.On("ToGraphQL", &modelEvent, &spec, &bundleRef).Return(gqlEvent, nil).Once()
				return conv
			},
			ExpectedEvent: gqlEvent,
			ExpectedErr:   nil,
		},
		{
			Name:            "Returns error when starting transaction",
			TransactionerFn: txGen.ThatFailsOnBegin,
			ServiceFn: func() *automock.EventDefService {
				return &automock.EventDefService{}
			},
			BndlServiceFn: func() *automock.BundleService {
				return &automock.BundleService{}
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				return &automock.EventDefConverter{}
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when bundle not exist",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				return &automock.EventDefService{}
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(nil, apperrors.NewNotFoundError(resource.Bundle, bundleID))
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   errors.New("cannot add Event Definition to not existing Bundle"),
		},
		{
			Name:            "Returns error when bundle existence check failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(nil, testErr)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when Event creation failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("CreateInBundle", txtest.CtxWithDBMatcher(), appID, bundleID, *modelEventInput, specInput).Return("", testErr).Once()
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(modelBndl, nil)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when Event retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("CreateInBundle", txtest.CtxWithDBMatcher(), appID, bundleID, *modelEventInput, specInput).Return(id, nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(nil, testErr).Once()
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(modelBndl, nil)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when Spec retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("CreateInBundle", txtest.CtxWithDBMatcher(), appID, bundleID, *modelEventInput, specInput).Return(id, nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEvent, nil).Once()
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(modelBndl, nil)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, gqlEvent.ID).Return(nil, testErr).Once()
				return svc
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when BundleReference retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("CreateInBundle", txtest.CtxWithDBMatcher(), appID, bundleID, *modelEventInput, specInput).Return(id, nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEvent, nil).Once()
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(modelBndl, nil)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &gqlEvent.ID, str.Ptr(bundleID)).Return(nil, testErr).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, gqlEvent.ID).Return(&spec, nil).Once()
				return svc
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Success",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("CreateInBundle", txtest.CtxWithDBMatcher(), appID, bundleID, *modelEventInput, specInput).Return(id, nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEvent, nil).Once()
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(modelBndl, nil)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &gqlEvent.ID, str.Ptr(bundleID)).Return(&bundleRef, nil).Once()
				return svc
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, gqlEvent.ID).Return(&spec, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				conv.On("ToGraphQL", &modelEvent, &spec, &bundleRef).Return(gqlEvent, testErr).Once()
				return conv
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when commit transaction failed",
			TransactionerFn: txGen.ThatFailsOnCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("CreateInBundle", txtest.CtxWithDBMatcher(), appID, bundleID, *modelEventInput, specInput).Return(id, nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEvent, nil).Once()
				return svc
			},
			BndlServiceFn: func() *automock.BundleService {
				appSvc := &automock.BundleService{}
				appSvc.On("Get", txtest.CtxWithDBMatcher(), bundleID).Return(modelBndl, nil)
				return appSvc
			},
			BndlReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &gqlEvent.ID, str.Ptr(bundleID)).Return(&bundleRef, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventInput).Return(modelEventInput, specInput, nil).Once()
				conv.On("ToGraphQL", &modelEvent, &spec, &bundleRef).Return(gqlEvent, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, gqlEvent.ID).Return(&spec, nil).Once()
				return svc
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// GIVEN
			persist, transact := testCase.TransactionerFn()
			svc := testCase.ServiceFn()
			converter := testCase.ConverterFn()
			bndlSvc := testCase.BndlServiceFn()
			specSvc := testCase.SpecServiceFn()
			bndlRefSvc := testCase.BndlReferenceServiceFn()

			resolver := event.NewResolver(transact, svc, bndlSvc, bndlRefSvc, converter, nil, specSvc, nil)

			// WHEN
			result, err := resolver.AddEventDefinitionToBundle(context.TODO(), bundleID, *gqlEventInput)

			// then
			assert.Equal(t, testCase.ExpectedEvent, result)
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				require.Nil(t, err)
			}

			persist.AssertExpectations(t)
			transact.AssertExpectations(t)
			svc.AssertExpectations(t)
			bndlSvc.AssertExpectations(t)
			bndlRefSvc.AssertExpectations(t)
			specSvc.AssertExpectations(t)
			converter.AssertExpectations(t)
		})
	}
}

func TestResolver_DeleteEvent(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "bar"
	modelEventDefinition, spec, bundleRef := fixFullEventDefinitionModel("test")
	gqlEventDefinition := fixFullGQLEventDefinition("test")

	var nilBundleID *string

	txGen := txtest.NewTransactionContextGenerator(testErr)

	testCases := []struct {
		Name              string
		TransactionerFn   func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		ServiceFn         func() *automock.EventDefService
		ConverterFn       func() *automock.EventDefConverter
		SpecServiceFn     func() *automock.SpecService
		BundleReferenceFn func() *automock.BundleReferenceService
		ExpectedEvent     *graphql.EventDefinition
		ExpectedErr       error
	}{
		{
			Name:            "Success",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				svc.On("Delete", txtest.CtxWithDBMatcher(), id).Return(nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("ToGraphQL", &modelEventDefinition, &spec, &bundleRef).Return(gqlEventDefinition, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&spec, nil).Once()
				return svc
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(&bundleRef, nil).Once()
				return svc
			},
			ExpectedEvent: gqlEventDefinition,
			ExpectedErr:   nil,
		},
		{
			Name:            "Return error when starting transaction fails",
			TransactionerFn: txGen.ThatFailsOnBegin,
			ServiceFn: func() *automock.EventDefService {
				return &automock.EventDefService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				return &automock.EventDefConverter{}
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when Event retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(nil, testErr).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				return &automock.EventDefConverter{}
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when Spec retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				return &automock.EventDefConverter{}
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(nil, testErr).Once()
				return svc
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when BundleReference retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				return &automock.EventDefConverter{}
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&spec, nil).Once()
				return svc
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(nil, testErr).Once()
				return svc
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when Event conversion failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("ToGraphQL", &modelEventDefinition, &spec, &bundleRef).Return(nil, testErr).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&spec, nil).Once()
				return svc
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(&bundleRef, nil).Once()
				return svc
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Returns error when Event deletion failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				svc.On("Delete", txtest.CtxWithDBMatcher(), id).Return(testErr).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("ToGraphQL", &modelEventDefinition, &spec, &bundleRef).Return(gqlEventDefinition, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&spec, nil).Once()
				return svc
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(&bundleRef, nil).Once()
				return svc
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
		{
			Name:            "Return error when commit transaction fails",
			TransactionerFn: txGen.ThatFailsOnCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				svc.On("Delete", txtest.CtxWithDBMatcher(), id).Return(nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("ToGraphQL", &modelEventDefinition, &spec, &bundleRef).Return(gqlEventDefinition, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&spec, nil).Once()
				return svc
			},
			BundleReferenceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(&bundleRef, nil).Once()
				return svc
			},
			ExpectedEvent: nil,
			ExpectedErr:   testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// GIVEN
			persist, transact := testCase.TransactionerFn()
			svc := testCase.ServiceFn()
			specService := testCase.SpecServiceFn()
			converter := testCase.ConverterFn()
			bndlRefService := testCase.BundleReferenceFn()

			resolver := event.NewResolver(transact, svc, nil, bndlRefService, converter, nil, specService, nil)

			// WHEN
			result, err := resolver.DeleteEventDefinition(context.TODO(), id)

			// then
			assert.Equal(t, testCase.ExpectedEvent, result)
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				require.Nil(t, err)
			}

			svc.AssertExpectations(t)
			specService.AssertExpectations(t)
			bndlRefService.AssertExpectations(t)
			converter.AssertExpectations(t)
			transact.AssertExpectations(t)
			persist.AssertExpectations(t)
		})
	}
}

func TestResolver_UpdateEvent(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "bar"
	gqlEventDefinitionInput := fixGQLEventDefinitionInput(id, "foo", "bar")
	modelEventDefinitionInput, modelSpecInput := fixModelEventDefinitionInput(id, "foo", "bar")
	gqlEventDefinition := fixFullGQLEventDefinition("test")
	modelEventDefinition, modelSpec, modelBundleReference := fixFullEventDefinitionModel("test")

	var nilBundleID *string

	txGen := txtest.NewTransactionContextGenerator(testErr)

	testCases := []struct {
		Name                     string
		TransactionerFn          func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		ServiceFn                func() *automock.EventDefService
		ConverterFn              func() *automock.EventDefConverter
		SpecServiceFn            func() *automock.SpecService
		BundleReferenceServiceFn func() *automock.BundleReferenceService
		InputWebhookID           string
		InputEvent               graphql.EventDefinitionInput
		ExpectedEventDefinition  *graphql.EventDefinition
		ExpectedErr              error
	}{
		{
			Name:            "Success",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Update", txtest.CtxWithDBMatcher(), id, *modelEventDefinitionInput, modelSpecInput).Return(nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(modelEventDefinitionInput, modelSpecInput, nil).Once()
				conv.On("ToGraphQL", &modelEventDefinition, &modelSpec, &modelBundleReference).Return(gqlEventDefinition, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&modelSpec, nil).Once()
				return svc
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(&modelBundleReference, nil).Once()
				return svc
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: gqlEventDefinition,
			ExpectedErr:             nil,
		},
		{
			Name:            "Returns error when starting transaction failed",
			TransactionerFn: txGen.ThatFailsOnBegin,
			ServiceFn: func() *automock.EventDefService {
				return &automock.EventDefService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				return &automock.EventDefConverter{}
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
		{
			Name:            "Returns error when converting input to GraphQL fails",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				return &automock.EventDefService{}
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(nil, nil, testErr).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
		{
			Name:            "Returns error when Event update failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Update", txtest.CtxWithDBMatcher(), id, *modelEventDefinitionInput, modelSpecInput).Return(testErr).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(modelEventDefinitionInput, modelSpecInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
		{
			Name:            "Returns error when Event retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Update", txtest.CtxWithDBMatcher(), id, *modelEventDefinitionInput, modelSpecInput).Return(nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(nil, testErr).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(modelEventDefinitionInput, modelSpecInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
		{
			Name:            "Returns error when Spec retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Update", txtest.CtxWithDBMatcher(), id, *modelEventDefinitionInput, modelSpecInput).Return(nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(modelEventDefinitionInput, modelSpecInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(nil, testErr).Once()
				return svc
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				return &automock.BundleReferenceService{}
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
		{
			Name:            "Returns error when BundlerReference retrieval failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Update", txtest.CtxWithDBMatcher(), id, *modelEventDefinitionInput, modelSpecInput).Return(nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(modelEventDefinitionInput, modelSpecInput, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&modelSpec, nil).Once()
				return svc
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(nil, testErr).Once()
				return svc
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
		{
			Name:            "Returns error when converting to GraphQL failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Update", txtest.CtxWithDBMatcher(), id, *modelEventDefinitionInput, modelSpecInput).Return(nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(modelEventDefinitionInput, modelSpecInput, nil).Once()
				conv.On("ToGraphQL", &modelEventDefinition, &modelSpec, &modelBundleReference).Return(nil, testErr).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&modelSpec, nil).Once()
				return svc
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(&modelBundleReference, nil).Once()
				return svc
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
		{
			Name:            "Returns error when commit transaction failed",
			TransactionerFn: txGen.ThatFailsOnCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("Update", txtest.CtxWithDBMatcher(), id, *modelEventDefinitionInput, modelSpecInput).Return(nil).Once()
				svc.On("Get", txtest.CtxWithDBMatcher(), id).Return(&modelEventDefinition, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.EventDefConverter {
				conv := &automock.EventDefConverter{}
				conv.On("InputFromGraphQL", gqlEventDefinitionInput).Return(modelEventDefinitionInput, modelSpecInput, nil).Once()
				conv.On("ToGraphQL", &modelEventDefinition, &modelSpec, &modelBundleReference).Return(gqlEventDefinition, nil).Once()
				return conv
			},
			SpecServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, modelEventDefinition.ID).Return(&modelSpec, nil).Once()
				return svc
			},
			BundleReferenceServiceFn: func() *automock.BundleReferenceService {
				svc := &automock.BundleReferenceService{}
				svc.On("GetForBundle", txtest.CtxWithDBMatcher(), model.BundleEventReference, &modelEventDefinition.ID, nilBundleID).Return(&modelBundleReference, nil).Once()
				return svc
			},
			InputWebhookID:          id,
			InputEvent:              *gqlEventDefinitionInput,
			ExpectedEventDefinition: nil,
			ExpectedErr:             testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// GIVEN
			persist, transact := testCase.TransactionerFn()
			svc := testCase.ServiceFn()
			converter := testCase.ConverterFn()
			specService := testCase.SpecServiceFn()
			bndlRefService := testCase.BundleReferenceServiceFn()

			resolver := event.NewResolver(transact, svc, nil, bndlRefService, converter, nil, specService, nil)

			// WHEN
			result, err := resolver.UpdateEventDefinition(context.TODO(), id, *gqlEventDefinitionInput)

			// then
			assert.Equal(t, testCase.ExpectedEventDefinition, result)
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				require.Nil(t, err)
			}

			persist.AssertExpectations(t)
			transact.AssertExpectations(t)
			svc.AssertExpectations(t)
			specService.AssertExpectations(t)
			bndlRefService.AssertExpectations(t)
			converter.AssertExpectations(t)
		})
	}
}

func TestResolver_RefetchEventSpec(t *testing.T) {
	// GIVEN
	testErr := errors.New("test error")

	eventID := "eventID"
	specID := "specID"

	dataBytes := "data"
	modelSpec := &model.Spec{
		ID:   specID,
		Data: &dataBytes,
	}

	clob := graphql.CLOB(dataBytes)
	gqlEventSpec := &graphql.EventSpec{
		Data: &clob,
	}

	txGen := txtest.NewTransactionContextGenerator(testErr)

	testCases := []struct {
		Name              string
		TransactionerFn   func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		ServiceFn         func() *automock.SpecService
		ConvFn            func() *automock.SpecConverter
		ExpectedEventSpec *graphql.EventSpec
		ExpectedErr       error
	}{
		{
			Name:            "Success",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, eventID).Return(modelSpec, nil).Once()
				svc.On("RefetchSpec", txtest.CtxWithDBMatcher(), specID, model.EventSpecReference).Return(modelSpec, nil).Once()
				return svc
			},
			ConvFn: func() *automock.SpecConverter {
				conv := &automock.SpecConverter{}
				conv.On("ToGraphQLEventSpec", modelSpec).Return(gqlEventSpec, nil).Once()
				return conv
			},
			ExpectedEventSpec: gqlEventSpec,
			ExpectedErr:       nil,
		},
		{
			Name:            "Returns error when starting transaction",
			TransactionerFn: txGen.ThatFailsOnBegin,
			ServiceFn: func() *automock.SpecService {
				return &automock.SpecService{}
			},
			ConvFn: func() *automock.SpecConverter {
				return &automock.SpecConverter{}
			},
			ExpectedEventSpec: nil,
			ExpectedErr:       testErr,
		},
		{
			Name:            "Returns error when getting spec failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, eventID).Return(nil, testErr).Once()
				return svc
			},
			ConvFn: func() *automock.SpecConverter {
				return &automock.SpecConverter{}
			},
			ExpectedEventSpec: nil,
			ExpectedErr:       testErr,
		},
		{
			Name:            "Returns error when spec not found",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, eventID).Return(nil, nil).Once()
				return svc
			},
			ConvFn: func() *automock.SpecConverter {
				return &automock.SpecConverter{}
			},
			ExpectedEventSpec: nil,
			ExpectedErr:       errors.Errorf("spec for Event with id %q not found", eventID),
		},
		{
			Name:            "Returns error when refetching event spec failed",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, eventID).Return(modelSpec, nil).Once()
				svc.On("RefetchSpec", txtest.CtxWithDBMatcher(), specID, model.EventSpecReference).Return(nil, testErr).Once()
				return svc
			},
			ConvFn: func() *automock.SpecConverter {
				return &automock.SpecConverter{}
			},
			ExpectedEventSpec: nil,
			ExpectedErr:       testErr,
		},
		{
			Name:            "Returns error converting to GraphQL fails",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, eventID).Return(modelSpec, nil).Once()
				svc.On("RefetchSpec", txtest.CtxWithDBMatcher(), specID, model.EventSpecReference).Return(modelSpec, nil).Once()
				return svc
			},
			ConvFn: func() *automock.SpecConverter {
				conv := &automock.SpecConverter{}
				conv.On("ToGraphQLEventSpec", modelSpec).Return(nil, testErr).Once()
				return conv
			},
			ExpectedEventSpec: nil,
			ExpectedErr:       testErr,
		},
		{
			Name:            "Returns error when commit transaction failed",
			TransactionerFn: txGen.ThatFailsOnCommit,
			ServiceFn: func() *automock.SpecService {
				svc := &automock.SpecService{}
				svc.On("GetByReferenceObjectID", txtest.CtxWithDBMatcher(), model.EventSpecReference, eventID).Return(modelSpec, nil).Once()
				svc.On("RefetchSpec", txtest.CtxWithDBMatcher(), specID, model.EventSpecReference).Return(modelSpec, nil).Once()
				return svc
			},
			ConvFn: func() *automock.SpecConverter {
				conv := &automock.SpecConverter{}
				conv.On("ToGraphQLEventSpec", modelSpec).Return(gqlEventSpec, nil).Once()
				return conv
			},
			ExpectedEventSpec: nil,
			ExpectedErr:       testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// GIVEN
			svc := testCase.ServiceFn()
			conv := testCase.ConvFn()
			persist, transact := testCase.TransactionerFn()
			resolver := event.NewResolver(transact, nil, nil, nil, nil, nil, svc, conv)

			// WHEN
			result, err := resolver.RefetchEventDefinitionSpec(context.TODO(), eventID)

			// then
			assert.Equal(t, testCase.ExpectedEventSpec, result)
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				require.Nil(t, err)
			}

			persist.AssertExpectations(t)
			transact.AssertExpectations(t)
			svc.AssertExpectations(t)
			conv.AssertExpectations(t)
		})
	}
}

func TestResolver_FetchRequest(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	firstSpecID := "specID"
	secondSpecID := "specID2"
	specIDs := []string{firstSpecID, secondSpecID}
	firstFRID := "frID"
	secondFRID := "frID2"
	frURL := "foo.bar"
	timestamp := time.Now()

	frFirstSpec := fixModelFetchRequest(firstFRID, frURL, timestamp)
	frSecondSpec := fixModelFetchRequest(secondFRID, frURL, timestamp)
	fetchRequests := []*model.FetchRequest{frFirstSpec, frSecondSpec}

	gqlFRFirstSpec := fixGQLFetchRequest(frURL, timestamp)
	gqlFRSecondSpec := fixGQLFetchRequest(frURL, timestamp)
	gqlFetchRequests := []*graphql.FetchRequest{gqlFRFirstSpec, gqlFRSecondSpec}

	txGen := txtest.NewTransactionContextGenerator(testErr)

	testCases := []struct {
		Name            string
		TransactionerFn func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		ServiceFn       func() *automock.EventDefService
		ConverterFn     func() *automock.FetchRequestConverter
		ExpectedResult  []*graphql.FetchRequest
		ExpectedErr     []error
	}{
		{
			Name:            "Success",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("ListFetchRequests", txtest.CtxWithDBMatcher(), specIDs).Return(fetchRequests, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				conv.On("ToGraphQL", frFirstSpec).Return(gqlFRFirstSpec, nil).Once()
				conv.On("ToGraphQL", frSecondSpec).Return(gqlFRSecondSpec, nil).Once()
				return conv
			},
			ExpectedResult: gqlFetchRequests,
			ExpectedErr:    nil,
		},
		{
			Name:            "Returns error when starting transaction failed",
			TransactionerFn: txGen.ThatFailsOnBegin,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				return svc
			},
			ConverterFn: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				return conv
			},
			ExpectedResult: nil,
			ExpectedErr:    []error{testErr},
		},
		{
			Name:            "FetchRequest doesn't exist",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("ListFetchRequests", txtest.CtxWithDBMatcher(), specIDs).Return(nil, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				return conv
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name:            "Error when listing Event FetchRequests",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("ListFetchRequests", txtest.CtxWithDBMatcher(), specIDs).Return(nil, testErr).Once()
				return svc
			},
			ConverterFn: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				return conv
			},
			ExpectedResult: nil,
			ExpectedErr:    []error{testErr},
		},
		{
			Name:            "Error when converting FetchRequest to graphql",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("ListFetchRequests", txtest.CtxWithDBMatcher(), specIDs).Return(fetchRequests, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				conv.On("ToGraphQL", frFirstSpec).Return(nil, testErr).Once()
				return conv
			},
			ExpectedResult: nil,
			ExpectedErr:    []error{testErr},
		},
		{
			Name:            "Returns error when commit transaction fails",
			TransactionerFn: txGen.ThatFailsOnCommit,
			ServiceFn: func() *automock.EventDefService {
				svc := &automock.EventDefService{}
				svc.On("ListFetchRequests", txtest.CtxWithDBMatcher(), specIDs).Return(fetchRequests, nil).Once()
				return svc
			},
			ConverterFn: func() *automock.FetchRequestConverter {
				conv := &automock.FetchRequestConverter{}
				conv.On("ToGraphQL", frFirstSpec).Return(gqlFRFirstSpec, nil).Once()
				conv.On("ToGraphQL", frSecondSpec).Return(gqlFRSecondSpec, nil).Once()
				return conv
			},
			ExpectedResult: nil,
			ExpectedErr:    []error{testErr},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			persist, transact := testCase.TransactionerFn()
			svc := testCase.ServiceFn()
			converter := testCase.ConverterFn()

			firstFRParams := dataloader.ParamFetchRequestEventDef{ID: firstSpecID, Ctx: context.TODO()}
			secondFRParams := dataloader.ParamFetchRequestEventDef{ID: secondSpecID, Ctx: context.TODO()}
			keys := []dataloader.ParamFetchRequestEventDef{firstFRParams, secondFRParams}
			resolver := event.NewResolver(transact, svc, nil, nil, nil, converter, nil, nil)

			// WHEN
			result, err := resolver.FetchRequestEventDefDataLoader(keys)

			// then
			assert.Equal(t, testCase.ExpectedResult, result)
			assert.Equal(t, testCase.ExpectedErr, err)

			svc.AssertExpectations(t)
			converter.AssertExpectations(t)
			persist.AssertExpectations(t)
			transact.AssertExpectations(t)
		})
	}
	t.Run("Returns error when there are no Specs", func(t *testing.T) {
		resolver := event.NewResolver(nil, nil, nil, nil, nil, nil, nil, nil)
		// WHEN
		_, err := resolver.FetchRequestEventDefDataLoader([]dataloader.ParamFetchRequestEventDef{})
		// THEN
		require.Error(t, err[0])
		assert.EqualError(t, err[0], apperrors.NewInternalError("No EventDef specs found").Error())
	})

	t.Run("Returns error when Specification ID is empty", func(t *testing.T) {
		params := dataloader.ParamFetchRequestEventDef{ID: "", Ctx: context.TODO()}
		keys := []dataloader.ParamFetchRequestEventDef{params}

		resolver := event.NewResolver(nil, nil, nil, nil, nil, nil, nil, nil)
		// WHEN
		_, err := resolver.FetchRequestEventDefDataLoader(keys)
		// THEN
		require.Error(t, err[0])
		assert.EqualError(t, err[0], apperrors.NewInternalError("Cannot fetch FetchRequest. EventDefinition Spec ID is empty").Error())
	})
}

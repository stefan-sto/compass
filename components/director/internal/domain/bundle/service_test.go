package bundle_test

import (
	"context"
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal/domain/bundle"
	"github.com/kyma-incubator/compass/components/director/internal/domain/bundle/automock"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "foo"
	applicationID := appID
	name := "foo"
	desc := "bar"

	modelInput := model.BundleCreateInput{
		Name:                           name,
		Description:                    &desc,
		InstanceAuthRequestInputSchema: fixBasicSchema(),
		DefaultInstanceAuth:            &model.AuthInput{},
		Documents: []*model.DocumentInput{
			{Title: "foo", Description: "test", FetchRequest: &model.FetchRequestInput{URL: "doc.foo.bar"}},
			{Title: "bar", Description: "test"},
		},
		APIDefinitions: []*model.APIDefinitionInput{
			{
				Name: "foo",
			},
			{
				Name: "bar",
			},
		},
		APISpecs: []*model.SpecInput{
			{
				FetchRequest: &model.FetchRequestInput{URL: "api.foo.bar"},
			},
			nil,
		},
		EventDefinitions: []*model.EventDefinitionInput{
			{
				Name: "foo",
			},
			{
				Name: "bar",
			},
		},
		EventSpecs: []*model.SpecInput{
			{
				FetchRequest: &model.FetchRequestInput{URL: "eventapi.foo.bar"},
			},
			nil,
		},
	}

	modelBundle := &model.Bundle{
		ApplicationID:                  applicationID,
		Name:                           name,
		Description:                    &desc,
		InstanceAuthRequestInputSchema: fixBasicSchema(),
		DefaultInstanceAuth:            &model.Auth{},
		BaseEntity: &model.BaseEntity{
			ID:    id,
			Ready: true,
		},
	}
	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name              string
		RepositoryFn      func() *automock.BundleRepository
		APIServiceFn      func() *automock.APIService
		EventServiceFn    func() *automock.EventService
		DocumentServiceFn func() *automock.DocumentService
		UIDServiceFn      func() *automock.UIDService
		Input             model.BundleCreateInput
		ExpectedErr       error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("Create", ctx, tenantID, modelBundle).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id)
				return svc
			},
			APIServiceFn: func() *automock.APIService {
				svc := &automock.APIService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.APIDefinitions[0], modelInput.APISpecs[0]).Return("", nil).Once()
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.APIDefinitions[1], modelInput.APISpecs[1]).Return("", nil).Once()
				return svc
			},
			EventServiceFn: func() *automock.EventService {
				svc := &automock.EventService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.EventDefinitions[0], modelInput.EventSpecs[0]).Return("", nil).Once()
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.EventDefinitions[1], modelInput.EventSpecs[1]).Return("", nil).Once()
				return svc
			},
			DocumentServiceFn: func() *automock.DocumentService {
				svc := &automock.DocumentService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.Documents[0]).Return("", nil).Once()
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.Documents[1]).Return("", nil).Once()
				return svc
			},
			Input:       modelInput,
			ExpectedErr: nil,
		},
		{
			Name: "Error - Bundle creation",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("Create", ctx, tenantID, modelBundle).Return(testErr).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id).Once()
				return svc
			},
			APIServiceFn: func() *automock.APIService {
				return &automock.APIService{}
			},
			EventServiceFn: func() *automock.EventService {
				return &automock.EventService{}
			},
			DocumentServiceFn: func() *automock.DocumentService {
				return &automock.DocumentService{}
			},
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Error - API creation",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("Create", ctx, tenantID, modelBundle).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id)
				return svc
			},
			APIServiceFn: func() *automock.APIService {
				svc := &automock.APIService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.APIDefinitions[0], modelInput.APISpecs[0]).Return("", testErr).Once()
				return svc
			},
			EventServiceFn: func() *automock.EventService {
				svc := &automock.EventService{}
				return svc
			},
			DocumentServiceFn: func() *automock.DocumentService {
				svc := &automock.DocumentService{}
				return svc
			},
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Error - Event creation",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("Create", ctx, tenantID, modelBundle).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id)
				return svc
			},
			APIServiceFn: func() *automock.APIService {
				svc := &automock.APIService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.APIDefinitions[0], modelInput.APISpecs[0]).Return("", nil).Once()
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.APIDefinitions[1], modelInput.APISpecs[1]).Return("", nil).Once()
				return svc
			},
			EventServiceFn: func() *automock.EventService {
				svc := &automock.EventService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.EventDefinitions[0], modelInput.EventSpecs[0]).Return("", testErr).Once()
				return svc
			},
			DocumentServiceFn: func() *automock.DocumentService {
				repo := &automock.DocumentService{}
				return repo
			},
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Error - Document creation",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("Create", ctx, tenantID, modelBundle).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id)
				return svc
			},
			APIServiceFn: func() *automock.APIService {
				svc := &automock.APIService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.APIDefinitions[0], modelInput.APISpecs[0]).Return("", nil).Once()
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.APIDefinitions[1], modelInput.APISpecs[1]).Return("", nil).Once()
				return svc
			},
			EventServiceFn: func() *automock.EventService {
				svc := &automock.EventService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.EventDefinitions[0], modelInput.EventSpecs[0]).Return("", nil).Once()
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.EventDefinitions[1], modelInput.EventSpecs[1]).Return("", nil).Once()
				return svc
			},
			DocumentServiceFn: func() *automock.DocumentService {
				svc := &automock.DocumentService{}
				svc.On("CreateInBundle", ctx, appID, id, *modelInput.Documents[0]).Return("", testErr).Once()
				return svc
			},
			Input:       modelInput,
			ExpectedErr: testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// GIVEN
			repo := testCase.RepositoryFn()
			uidService := testCase.UIDServiceFn()

			apiSvc := testCase.APIServiceFn()
			eventSvc := testCase.EventServiceFn()
			documentSvc := testCase.DocumentServiceFn()
			svc := bundle.NewService(repo, apiSvc, eventSvc, documentSvc, uidService)

			// WHEN
			result, err := svc.Create(ctx, applicationID, testCase.Input)

			// THEN
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				assert.IsType(t, "string", result)
			}

			mock.AssertExpectationsForObjects(t, repo, apiSvc, eventSvc, documentSvc, uidService)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		_, err := svc.Create(context.TODO(), "", model.BundleCreateInput{})
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_Update(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "foo"
	name := "bar"
	desc := "baz"

	modelInput := model.BundleUpdateInput{
		Name:                           name,
		Description:                    &desc,
		InstanceAuthRequestInputSchema: fixBasicSchema(),
		DefaultInstanceAuth:            &model.AuthInput{},
	}

	inputBundleModel := mock.MatchedBy(func(bndl *model.Bundle) bool {
		return bndl.Name == modelInput.Name
	})

	bundleModel := &model.Bundle{
		ApplicationID:                  "id",
		Name:                           name,
		Description:                    &desc,
		InstanceAuthRequestInputSchema: fixBasicSchema(),
		DefaultInstanceAuth:            &model.Auth{},
		BaseEntity:                     &model.BaseEntity{ID: id},
	}

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name         string
		RepositoryFn func() *automock.BundleRepository
		Input        model.BundleUpdateInput
		InputID      string
		ExpectedErr  error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(bundleModel, nil).Once()
				repo.On("Update", ctx, tenantID, inputBundleModel).Return(nil).Once()
				return repo
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: nil,
		},
		{
			Name: "Update Error",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("GetByID", ctx, tenantID, "foo").Return(bundleModel, nil).Once()
				repo.On("Update", ctx, tenantID, inputBundleModel).Return(testErr).Once()
				return repo
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Get Error",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("GetByID", ctx, tenantID, "foo").Return(nil, testErr).Once()
				return repo
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// GIVEN
			repo := testCase.RepositoryFn()

			svc := bundle.NewService(repo, nil, nil, nil, nil)

			// WHEN
			err := svc.Update(ctx, testCase.InputID, testCase.Input)

			// THEN
			if testCase.ExpectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		err := svc.Update(context.TODO(), "", model.BundleUpdateInput{})
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_Delete(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "foo"

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name         string
		RepositoryFn func() *automock.BundleRepository
		Input        model.BundleCreateInput
		InputID      string
		ExpectedErr  error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("Delete", ctx, tenantID, id).Return(nil).Once()
				return repo
			},
			InputID:     id,
			ExpectedErr: nil,
		},
		{
			Name: "Delete Error",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("Delete", ctx, tenantID, id).Return(testErr).Once()
				return repo
			},
			InputID:     id,
			ExpectedErr: testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// GIVEN
			repo := testCase.RepositoryFn()

			svc := bundle.NewService(repo, nil, nil, nil, nil)

			// WHEN
			err := svc.Delete(ctx, testCase.InputID)

			// THEN
			if testCase.ExpectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		err := svc.Delete(context.TODO(), "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_Exist(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")
	ctx := tenant.SaveToContext(context.TODO(), tenantID, externalTenantID)
	id := "foo"

	testCases := []struct {
		Name           string
		RepoFn         func() *automock.BundleRepository
		ExpectedError  error
		ExpectedOutput bool
	}{
		{
			Name: "Success",
			RepoFn: func() *automock.BundleRepository {
				bndlRepo := &automock.BundleRepository{}
				bndlRepo.On("Exists", ctx, tenantID, id).Return(true, nil).Once()
				return bndlRepo
			},
			ExpectedOutput: true,
		},
		{
			Name: "Error when getting Bundle",
			RepoFn: func() *automock.BundleRepository {
				bndlRepo := &automock.BundleRepository{}
				bndlRepo.On("Exists", ctx, tenantID, id).Return(false, testErr).Once()
				return bndlRepo
			},
			ExpectedError:  testErr,
			ExpectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			bndlRepo := testCase.RepoFn()
			svc := bundle.NewService(bndlRepo, nil, nil, nil, nil)

			// WHEN
			result, err := svc.Exist(ctx, id)

			// THEN
			if testCase.ExpectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.ExpectedOutput, result)

			bndlRepo.AssertExpectations(t)
		})
	}

	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		_, err := svc.Exist(context.TODO(), "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_Get(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "foo"
	name := "foo"
	desc := "bar"

	bndl := fixBundleModel(name, desc)

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name               string
		RepositoryFn       func() *automock.BundleRepository
		Input              model.BundleCreateInput
		InputID            string
		ExpectedBundle     *model.Bundle
		ExpectedErrMessage string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(bndl, nil).Once()
				return repo
			},
			InputID:            id,
			ExpectedBundle:     bndl,
			ExpectedErrMessage: "",
		},
		{
			Name: "Returns error when Bundle retrieval failed",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(nil, testErr).Once()
				return repo
			},
			InputID:            id,
			ExpectedBundle:     bndl,
			ExpectedErrMessage: testErr.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()
			svc := bundle.NewService(repo, nil, nil, nil, nil)

			// WHEN
			bndl, err := svc.Get(ctx, testCase.InputID)

			// THEN
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedBundle, bndl)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		_, err := svc.Get(context.TODO(), "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_GetForApplication(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	id := "foo"
	appID := "bar"
	name := "foo"
	desc := "bar"

	bndl := fixBundleModel(name, desc)

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name               string
		RepositoryFn       func() *automock.BundleRepository
		Input              model.BundleCreateInput
		InputID            string
		ApplicationID      string
		ExpectedBundle     *model.Bundle
		ExpectedErrMessage string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("GetForApplication", ctx, tenantID, id, appID).Return(bndl, nil).Once()
				return repo
			},
			InputID:            id,
			ApplicationID:      appID,
			ExpectedBundle:     bndl,
			ExpectedErrMessage: "",
		},
		{
			Name: "Returns error when Bundle retrieval failed",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("GetForApplication", ctx, tenantID, id, appID).Return(nil, testErr).Once()
				return repo
			},
			InputID:            id,
			ApplicationID:      appID,
			ExpectedBundle:     nil,
			ExpectedErrMessage: testErr.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()
			svc := bundle.NewService(repo, nil, nil, nil, nil)

			// WHEN
			document, err := svc.GetForApplication(ctx, testCase.InputID, testCase.ApplicationID)

			// THEN
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedBundle, document)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		_, err := svc.GetForApplication(context.TODO(), "", "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_ListByApplicationIDNoPaging(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	name := "foo"
	desc := "bar"

	bundles := []*model.Bundle{
		fixBundleModel(name, desc),
		fixBundleModel(name, desc),
		fixBundleModel(name, desc),
	}

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name               string
		RepositoryFn       func() *automock.BundleRepository
		ExpectedResult     []*model.Bundle
		ExpectedErrMessage string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("ListByApplicationIDNoPaging", ctx, tenantID, appID).Return(bundles, nil).Once()
				return repo
			},
			ExpectedResult:     bundles,
			ExpectedErrMessage: "",
		},
		{
			Name: "Returns error when Bundle listing failed",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("ListByApplicationIDNoPaging", ctx, tenantID, appID).Return(nil, testErr).Once()
				return repo
			},
			ExpectedResult:     nil,
			ExpectedErrMessage: testErr.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()

			svc := bundle.NewService(repo, nil, nil, nil, nil)

			// WHEN
			docs, err := svc.ListByApplicationIDNoPaging(ctx, appID)

			// THEN
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, docs)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		_, err := svc.ListByApplicationIDNoPaging(context.TODO(), "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_ListByApplicationIDs(t *testing.T) {
	// GIVEN
	testErr := errors.New("Test error")

	firstAppID := "bar"
	secondAppID := "bar2"
	name := "foo"
	desc := "bar"
	appIDs := []string{firstAppID, secondAppID}

	bundleFirstApp := fixBundleModel(name, desc)
	bundleFirstApp.ApplicationID = firstAppID
	bundleSecondApp := fixBundleModel(name, desc)
	bundleSecondApp.ApplicationID = secondAppID

	bundlesFirstApp := []*model.Bundle{bundleFirstApp}
	bundlesSecondApp := []*model.Bundle{bundleSecondApp}

	bundlePageFirstApp := &model.BundlePage{
		Data:       bundlesFirstApp,
		TotalCount: len(bundlesFirstApp),
		PageInfo: &pagination.Page{
			HasNextPage: false,
			EndCursor:   "end",
			StartCursor: "start",
		},
	}
	bundlePageSecondApp := &model.BundlePage{
		Data:       bundlesSecondApp,
		TotalCount: len(bundlesSecondApp),
		PageInfo: &pagination.Page{
			HasNextPage: false,
			EndCursor:   "end",
			StartCursor: "start",
		},
	}

	bundlePages := []*model.BundlePage{bundlePageFirstApp, bundlePageSecondApp}

	after := "test"

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name               string
		PageSize           int
		RepositoryFn       func() *automock.BundleRepository
		ExpectedResult     []*model.BundlePage
		ExpectedErrMessage string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("ListByApplicationIDs", ctx, tenantID, appIDs, 2, after).Return(bundlePages, nil).Once()
				return repo
			},
			PageSize:           2,
			ExpectedResult:     bundlePages,
			ExpectedErrMessage: "",
		},
		{
			Name: "Return error when page size is less than 1",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				return repo
			},
			PageSize:           0,
			ExpectedResult:     bundlePages,
			ExpectedErrMessage: "page size must be between 1 and 200",
		},
		{
			Name: "Return error when page size is bigger than 200",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				return repo
			},
			PageSize:           201,
			ExpectedResult:     bundlePages,
			ExpectedErrMessage: "page size must be between 1 and 200",
		},
		{
			Name: "Returns error when Bundle listing failed",
			RepositoryFn: func() *automock.BundleRepository {
				repo := &automock.BundleRepository{}
				repo.On("ListByApplicationIDs", ctx, tenantID, appIDs, 2, after).Return(nil, testErr).Once()
				return repo
			},
			PageSize:           2,
			ExpectedResult:     nil,
			ExpectedErrMessage: testErr.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()

			svc := bundle.NewService(repo, nil, nil, nil, nil)

			// WHEN
			bndls, err := svc.ListByApplicationIDs(ctx, appIDs, testCase.PageSize, after)

			// then
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, bndls)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := bundle.NewService(nil, nil, nil, nil, nil)
		// WHEN
		_, err := svc.ListByApplicationIDs(context.TODO(), nil, 5, "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

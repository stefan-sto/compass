package oauth20

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/internal/model"

	pkgmodel "github.com/kyma-incubator/compass/components/director/pkg/model"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/hashicorp/go-multierror"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/pkg/errors"
)

// SystemAuthService missing godoc
//go:generate mockery --name=SystemAuthService --output=automock --outpkg=automock --case=underscore --disable-version-string
type SystemAuthService interface {
	CreateWithCustomID(ctx context.Context, id string, objectType pkgmodel.SystemAuthReferenceObjectType, objectID string, authInput *model.AuthInput) (string, error)
	GetByIDForObject(ctx context.Context, objectType pkgmodel.SystemAuthReferenceObjectType, authID string) (*pkgmodel.SystemAuth, error)
}

// ApplicationService missing godoc
//go:generate mockery --name=ApplicationService --output=automock --outpkg=automock --case=underscore --disable-version-string
type ApplicationService interface {
	Exist(ctx context.Context, id string) (bool, error)
}

// RuntimeService missing godoc
//go:generate mockery --name=RuntimeService --output=automock --outpkg=automock --case=underscore --disable-version-string
type RuntimeService interface {
	Exist(ctx context.Context, id string) (bool, error)
}

// IntegrationSystemService missing godoc
//go:generate mockery --name=IntegrationSystemService --output=automock --outpkg=automock --case=underscore --disable-version-string
type IntegrationSystemService interface {
	Exists(ctx context.Context, id string) (bool, error)
}

// SystemAuthConverter missing godoc
//go:generate mockery --name=SystemAuthConverter --output=automock --outpkg=automock --case=underscore --disable-version-string
type SystemAuthConverter interface {
	ToGraphQL(model *pkgmodel.SystemAuth) (graphql.SystemAuth, error)
}

// Service missing godoc
//go:generate mockery --name=Service --output=automock --outpkg=automock --case=underscore --disable-version-string
type Service interface {
	CreateClientCredentials(ctx context.Context, objectType pkgmodel.SystemAuthReferenceObjectType) (*model.OAuthCredentialDataInput, error)
	DeleteClientCredentials(ctx context.Context, clientID string) error
}

// Resolver missing godoc
type Resolver struct {
	transact       persistence.Transactioner
	svc            Service
	systemAuthSvc  SystemAuthService
	systemAuthConv SystemAuthConverter
	appSvc         ApplicationService
	rtmSvc         RuntimeService
	isSvc          IntegrationSystemService
}

// NewResolver missing godoc
func NewResolver(transactioner persistence.Transactioner, svc Service, appSvc ApplicationService, rtmSvc RuntimeService, isSvc IntegrationSystemService, systemAuthSvc SystemAuthService, systemAuthConv SystemAuthConverter) *Resolver {
	return &Resolver{transact: transactioner, svc: svc, appSvc: appSvc, rtmSvc: rtmSvc, systemAuthSvc: systemAuthSvc, isSvc: isSvc, systemAuthConv: systemAuthConv}
}

// RequestClientCredentialsForRuntime missing godoc
func (r *Resolver) RequestClientCredentialsForRuntime(ctx context.Context, id string) (graphql.SystemAuth, error) {
	return r.generateClientCredentials(ctx, pkgmodel.RuntimeReference, id)
}

// RequestClientCredentialsForApplication missing godoc
func (r *Resolver) RequestClientCredentialsForApplication(ctx context.Context, id string) (graphql.SystemAuth, error) {
	return r.generateClientCredentials(ctx, pkgmodel.ApplicationReference, id)
}

// RequestClientCredentialsForIntegrationSystem missing godoc
func (r *Resolver) RequestClientCredentialsForIntegrationSystem(ctx context.Context, id string) (graphql.SystemAuth, error) {
	return r.generateClientCredentials(ctx, pkgmodel.IntegrationSystemReference, id)
}

func (r *Resolver) generateClientCredentials(ctx context.Context, objType pkgmodel.SystemAuthReferenceObjectType, objID string) (graphql.SystemAuth, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(ctx, tx)

	log.C(ctx).Infof("Requesting creation of client credentials for %s with id %s", objType, objID)
	ctx = persistence.SaveToContext(ctx, tx)

	exists, err := r.checkObjectExist(ctx, objType, objID)
	if err != nil {
		return nil, errors.Wrapf(err, "while checking if %s with ID '%s' exists", objType, objID)
	}
	if !exists {
		return nil, fmt.Errorf("%s with ID '%s' not found", objType, objID)
	}

	log.C(ctx).Debugf("Generating client credentials for %s with id %s by Director", objType, objID)
	clientCreds, err := r.svc.CreateClientCredentials(ctx, objType)
	if err != nil {
		return nil, errors.Wrapf(err, "while creating client credentials for %s with id %s", objType, objID)
	}
	if clientCreds == nil {
		return nil, apperrors.NewInvalidDataError("client credentials cannot be empty")
	}
	log.C(ctx).Debugf("Client credentials for %s with id %s are successfully generated by Director", objType, objID)
	cleanupOnError := func(originalErr error) error {
		cleanupErr := r.svc.DeleteClientCredentials(ctx, clientCreds.ClientID)
		if cleanupErr != nil {
			return multierror.Append(err, cleanupErr)
		}

		return originalErr
	}

	id := clientCreds.ClientID
	log.C(ctx).Debugf("Creating SystemAuth for the client credentials for %s with id %s", objType, objID)
	_, err = r.systemAuthSvc.CreateWithCustomID(ctx, id, objType, objID, &model.AuthInput{
		Credential: &model.CredentialDataInput{
			Oauth: clientCreds,
		},
	})
	if err != nil {
		finalErr := cleanupOnError(err)
		return nil, errors.Wrapf(finalErr, "error occurred while creating SystemAuth for %s with id %s", objType, objID)
	}
	log.C(ctx).Debugf("Successfully created SystemAuth for the client credentials for %s with id %s", objType, objID)

	sysAuth, err := r.systemAuthSvc.GetByIDForObject(ctx, objType, id)
	if err != nil {
		finalErr := cleanupOnError(err)
		return nil, finalErr
	}

	err = tx.Commit()
	if err != nil {
		finalErr := cleanupOnError(err)
		return nil, finalErr
	}

	log.C(ctx).Infof("Successfully created client credentials with client_id %s for %s with id %s", id, objType, objID)
	return r.systemAuthConv.ToGraphQL(sysAuth)
}

func (r *Resolver) checkObjectExist(ctx context.Context, objType pkgmodel.SystemAuthReferenceObjectType, objID string) (bool, error) {
	switch objType {
	case pkgmodel.RuntimeReference:
		return r.rtmSvc.Exist(ctx, objID)
	case pkgmodel.ApplicationReference:
		return r.appSvc.Exist(ctx, objID)
	case pkgmodel.IntegrationSystemReference:
		return r.isSvc.Exists(ctx, objID)
	}

	return false, fmt.Errorf("invalid object type %s", objType)
}

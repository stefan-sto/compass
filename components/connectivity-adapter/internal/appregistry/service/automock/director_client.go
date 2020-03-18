// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"
)

// DirectorClient is an autogenerated mock type for the DirectorClient type
type DirectorClient struct {
	mock.Mock
}

// CreateAPIDefinition provides a mock function with given fields: packageID, apiDefinitionInput
func (_m *DirectorClient) CreateAPIDefinition(packageID string, apiDefinitionInput graphql.APIDefinitionInput) (string, error) {
	ret := _m.Called(packageID, apiDefinitionInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, graphql.APIDefinitionInput) string); ok {
		r0 = rf(packageID, apiDefinitionInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, graphql.APIDefinitionInput) error); ok {
		r1 = rf(packageID, apiDefinitionInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateDocument provides a mock function with given fields: packageID, documentInput
func (_m *DirectorClient) CreateDocument(packageID string, documentInput graphql.DocumentInput) (string, error) {
	ret := _m.Called(packageID, documentInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, graphql.DocumentInput) string); ok {
		r0 = rf(packageID, documentInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, graphql.DocumentInput) error); ok {
		r1 = rf(packageID, documentInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateEventDefinition provides a mock function with given fields: packageID, eventDefinitionInput
func (_m *DirectorClient) CreateEventDefinition(packageID string, eventDefinitionInput graphql.EventDefinitionInput) (string, error) {
	ret := _m.Called(packageID, eventDefinitionInput)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, graphql.EventDefinitionInput) string); ok {
		r0 = rf(packageID, eventDefinitionInput)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, graphql.EventDefinitionInput) error); ok {
		r1 = rf(packageID, eventDefinitionInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePackage provides a mock function with given fields: appID, in
func (_m *DirectorClient) CreatePackage(appID string, in graphql.PackageCreateInput) (string, error) {
	ret := _m.Called(appID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, graphql.PackageCreateInput) string); ok {
		r0 = rf(appID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, graphql.PackageCreateInput) error); ok {
		r1 = rf(appID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAPIDefinition provides a mock function with given fields: apiID
func (_m *DirectorClient) DeleteAPIDefinition(apiID string) error {
	ret := _m.Called(apiID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(apiID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDocument provides a mock function with given fields: documentID
func (_m *DirectorClient) DeleteDocument(documentID string) error {
	ret := _m.Called(documentID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(documentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEventDefinition provides a mock function with given fields: eventID
func (_m *DirectorClient) DeleteEventDefinition(eventID string) error {
	ret := _m.Called(eventID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePackage provides a mock function with given fields: packageID
func (_m *DirectorClient) DeletePackage(packageID string) error {
	ret := _m.Called(packageID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(packageID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetPackage provides a mock function with given fields: appID, packageID
func (_m *DirectorClient) GetPackage(appID string, packageID string) (graphql.PackageExt, error) {
	ret := _m.Called(appID, packageID)

	var r0 graphql.PackageExt
	if rf, ok := ret.Get(0).(func(string, string) graphql.PackageExt); ok {
		r0 = rf(appID, packageID)
	} else {
		r0 = ret.Get(0).(graphql.PackageExt)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(appID, packageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPackages provides a mock function with given fields: appID
func (_m *DirectorClient) ListPackages(appID string) ([]*graphql.PackageExt, error) {
	ret := _m.Called(appID)

	var r0 []*graphql.PackageExt
	if rf, ok := ret.Get(0).(func(string) []*graphql.PackageExt); ok {
		r0 = rf(appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.PackageExt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePackage provides a mock function with given fields: packageID, in
func (_m *DirectorClient) UpdatePackage(packageID string, in graphql.PackageUpdateInput) error {
	ret := _m.Called(packageID, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, graphql.PackageUpdateInput) error); ok {
		r0 = rf(packageID, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

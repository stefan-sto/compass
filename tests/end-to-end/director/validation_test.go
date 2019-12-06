package director

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/end-to-end/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Runtime Validation

func TestCreateRuntime_Validation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	invalidInput := graphql.RuntimeInput{
		Name: "0invalid",
	}
	inputString, err := tc.graphqlizer.RuntimeInputToGQL(invalidInput)
	require.NoError(t, err)
	var result graphql.Runtime
	request := fixCreateRuntimeRequest(inputString)

	// WHEN
	err = tc.RunOperation(ctx, request, &result)

	// THEN
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation error for type RuntimeInput")
}

func TestUpdateRuntime_Validation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	rtm := createRuntime(t, ctx, "validation-test-rtm")
	defer deleteRuntime(t, rtm.ID)

	invalidInput := graphql.RuntimeInput{
		Name: "0invalid",
	}
	inputString, err := tc.graphqlizer.RuntimeInputToGQL(invalidInput)
	require.NoError(t, err)
	var result graphql.Runtime
	request := fixUpdateRuntimeRequest(rtm.ID, inputString)

	// WHEN
	err = tc.RunOperation(ctx, request, &result)

	// THEN
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation error for type RuntimeInput")
}

// Label Definition Validation

func TestCreateLabelDefinition_Validation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	invalidInput := graphql.LabelDefinitionInput{
		Key: "",
	}
	inputString, err := tc.graphqlizer.LabelDefinitionInputToGQL(invalidInput)
	require.NoError(t, err)
	var result graphql.Runtime
	request := fixCreateLabelDefinitionRequest(inputString)

	// WHEN
	err = tc.RunOperation(ctx, request, &result)

	// THEN
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation error for type LabelDefinitionInput")
}

func TestUpdateLabelDefinition_Validation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	key := "test-validation-ld"
	ld := createLabelDefinitionWithinTenant(t, ctx, key, map[string]string{"type": "string"}, defaultTenant)
	defer deleteLabelDefinitionWithinTenant(t, ctx, ld.Key, true, defaultTenant)
	invalidSchema := graphql.JSONSchema(`"{\"test\":}"`)
	invalidInput := graphql.LabelDefinitionInput{
		Key:    key,
		Schema: &invalidSchema,
	}
	inputString, err := tc.graphqlizer.LabelDefinitionInputToGQL(invalidInput)
	require.NoError(t, err)
	var result graphql.Runtime
	request := fixUpdateLabelDefinitionRequest(inputString)

	// WHEN
	err = tc.RunOperation(ctx, request, &result)

	// THEN
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation error for type LabelDefinitionInput")
}

// Label Validation

func TestSetApplicationLabel_Validation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	app := createApplication(t, ctx, "validation-test-app")
	defer deleteApplication(t, app.ID)

	request := fixSetApplicationLabelRequest(app.ID, strings.Repeat("x", 257), "")
	var result graphql.Label

	// WHEN
	err := tc.RunOperation(ctx, request, &result)

	// THEN
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation error for type LabelInput")
}

func TestSetRuntimeLabel_Validation(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	rtm := createRuntime(t, ctx, "validation-test-rtm")
	defer deleteRuntime(t, rtm.ID)

	request := fixSetRuntimeLabelRequest(rtm.ID, strings.Repeat("x", 257), "")
	var result graphql.Label

	// WHEN
	err := tc.RunOperation(ctx, request, &result)

	// THEN
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation error for type LabelInput")
}

// Auth Validation

const longDescErrorMsg = "graphql: validation error for type %s: description: the length must be no more than 128."

func TestCreateApplicationInput_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	app := fixSampleApplicationCreateInputWithName("placeholder", "name")
	longDesc := strings.Repeat("a", 129)
	app.Description = &longDesc

	appInputGQL, err := tc.graphqlizer.ApplicationCreateInputToGQL(app)
	require.NoError(t, err)
	createRequest := fixCreateApplicationRequest(appInputGQL)

	//WHEN
	err = tc.RunOperation(ctx, createRequest, nil)

	//THEN
	require.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf(longDescErrorMsg, "ApplicationCreateInput"))
}

func TestCreateApplicationUpdateInput_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	app := createApplication(t, ctx, "app-name")
	defer deleteApplication(t, app.ID)

	longDesc := strings.Repeat("a", 129)
	appUpdate := graphql.ApplicationUpdateInput{Name: "name", Description: &longDesc}
	appInputGQL, err := tc.graphqlizer.ApplicationUpdateInputToGQL(appUpdate)
	require.NoError(t, err)
	updateRequest := fixUpdateApplicationRequest(app.ID, appInputGQL)

	//WHEN
	err = tc.RunOperation(ctx, updateRequest, nil)

	//THEN
	require.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf(longDescErrorMsg, "ApplicationUpdateInput"))
}

func TestAddDocument_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	app := createApplication(t, ctx, "app-name")
	defer deleteApplication(t, app.ID)

	doc := fixDocumentInput()
	doc.DisplayName = strings.Repeat("a", 129)
	docInputGQL, err := tc.graphqlizer.DocumentInputToGQL(&doc)
	require.NoError(t, err)
	createRequest := fixAddDocumentRequest(app.ID, docInputGQL)

	//WHEN
	err = tc.RunOperation(ctx, createRequest, nil)

	//THEN
	require.Error(t, err)
	assert.EqualError(t, err, "graphql: validation error for type DocumentInput: displayName: the length must be between 1 and 128.")
}

func TestCreateIntegrationSystem_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	intSys := graphql.IntegrationSystemInput{Name: "valid-name"}
	longDesc := strings.Repeat("a", 129)
	intSys.Description = &longDesc

	isInputGQL, err := tc.graphqlizer.IntegrationSystemInputToGQL(intSys)
	require.NoError(t, err)
	createRequest := fixCreateIntegrationSystemRequest(isInputGQL)

	//WHEN
	err = tc.RunOperation(ctx, createRequest, nil)

	//THEN
	require.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf(longDescErrorMsg, "IntegrationSystemInput"))
}

func TestUpdateIntegrationSystem_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	intSys := createIntegrationSystem(t, ctx, "integration-system")
	defer deleteIntegrationSystem(t, ctx, intSys.ID)
	longDesc := strings.Repeat("a", 256)
	intSysUpdate := graphql.IntegrationSystemInput{Name: "name", Description: &longDesc}
	isUpdateGQL, err := tc.graphqlizer.IntegrationSystemInputToGQL(intSysUpdate)
	require.NoError(t, err)
	update := fixUpdateIntegrationSystemRequest(intSys.ID, isUpdateGQL)

	//WHEN
	err = tc.RunOperation(ctx, update, nil)

	//THEN
	require.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf(longDescErrorMsg, "IntegrationSystemInput"))
}

func TestAddAPI_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	app := createApplication(t, ctx, "name")
	defer deleteApplication(t, app.ID)

	api := graphql.APIDefinitionInput{Name: "name", TargetURL: "https://kyma project.io"}
	apiGQL, err := tc.graphqlizer.APIDefinitionInputToGQL(api)
	require.NoError(t, err)
	addAPIRequest := fixAddAPIRequest(app.ID, apiGQL)

	//WHEN
	err = tc.RunOperation(ctx, addAPIRequest, nil)

	//THEN
	require.Error(t, err)
	require.EqualError(t, err, "graphql: validation error for type APIDefinitionInput: targetURL: must be a valid URL.")
}

func TestUpdateAPI_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	app := createApplication(t, ctx, "name")
	defer deleteApplication(t, app.ID)

	api := graphql.APIDefinitionInput{Name: "name", TargetURL: "https://kyma-project.io"}
	addAPI(t, ctx, app.ID, api)

	api.TargetURL = "invalid URL"
	apiGQL, err := tc.graphqlizer.APIDefinitionInputToGQL(api)
	require.NoError(t, err)
	updateAPIRequest := fixUpdateAPIRequest(app.ID, apiGQL)

	//WHEN
	err = tc.RunOperation(ctx, updateAPIRequest, nil)

	//THEN
	require.Error(t, err)
	require.EqualError(t, err, "graphql: validation error for type APIDefinitionInput: targetURL: is not valid URL.")
}

func TestAddEventAPI_Validation(t *testing.T) {
	//GIVEN
	ctx := context.TODO()
	app := createApplication(t, ctx, "name")
	defer deleteApplication(t, app.ID)

	eventAPI := fixEventAPIDefinitionInput()
	longDesc := strings.Repeat("a", 129)
	eventAPI.Description = &longDesc
	evenApiGQL, err := tc.graphqlizer.EventAPIDefinitionInputToGQL(eventAPI)
	require.NoError(t, err)
	addEventAPIRequest := fixAddEventAPIRequest(app.ID, evenApiGQL)

	//WHEN
	err = tc.RunOperation(ctx, addEventAPIRequest, nil)

	//THEN
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf(longDescErrorMsg, "EventAPIDefinitionInput"))
}

func TestUpdateEventAPI_Validation(t *testing.T) {
	ctx := context.TODO()
	app := createApplication(t, ctx, "name")
	defer deleteApplication(t, app.ID)

	eventAPIUpdate := fixEventAPIDefinitionInput()
	eventAPI := addEventAPI(t, ctx, app.ID, eventAPIUpdate)

	longDesc := strings.Repeat("a", 129)
	eventAPIUpdate.Description = &longDesc
	evenApiGQL, err := tc.graphqlizer.EventAPIDefinitionInputToGQL(eventAPIUpdate)
	require.NoError(t, err)
	updateEventAPI := fixUpdateEventAPIRequest(eventAPI.ID, evenApiGQL)

	//WHEN
	err = tc.RunOperation(ctx, updateEventAPI, nil)

	//THEN
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf(longDescErrorMsg, "EventAPIDefinitionInput"))
}

func fixEventAPIDefinitionInput() graphql.EventAPIDefinitionInput {
	data := graphql.CLOB("data")
	return graphql.EventAPIDefinitionInput{Name: "name",
		Spec: &graphql.EventAPISpecInput{
			Data:          &data,
			EventSpecType: graphql.EventAPISpecTypeAsyncAPI,
			Format:        graphql.SpecFormatJSON,
		}}
}

func fixDocumentInput() graphql.DocumentInput {
	return graphql.DocumentInput{
		Title:       "Readme",
		Description: "Detailed description of project",
		Format:      graphql.DocumentFormatMarkdown,
		DisplayName: "display-name",
		FetchRequest: &graphql.FetchRequestInput{
			URL:    "kyma-project.io",
			Mode:   ptr.FetchMode(graphql.FetchModePackage),
			Filter: ptr.String("/docs/README.md"),
			Auth:   fixBasicAuth(),
		},
	}
}

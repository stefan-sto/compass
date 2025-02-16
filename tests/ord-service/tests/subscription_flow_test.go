/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tests

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/tests/pkg/subscription"

	testingx "github.com/kyma-incubator/compass/tests/pkg/testing"

	"github.com/kyma-incubator/compass/tests/pkg/certs/certprovider"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/kyma-incubator/compass/tests/pkg/ptr"
	"github.com/kyma-incubator/compass/tests/pkg/tenantfetcher"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	"github.com/kyma-incubator/compass/tests/pkg/token"
	gcli "github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestSelfRegisterFlow(t *testing.T) {
	ctx := context.Background()
	accountTenantID := conf.AccountTenantID // accountTenantID is parent of the tenant/subaccountID of the configured certificate client's tenant below

	// Register application
	app, err := fixtures.RegisterApplication(t, ctx, certSecuredGraphQLClient, "testingApp", accountTenantID)
	defer fixtures.CleanupApplication(t, ctx, certSecuredGraphQLClient, accountTenantID, &app)
	require.NoError(t, err)
	require.NotEmpty(t, app.ID)

	formationName := "sr-test-scenario"
	t.Logf("Creating formation with name %s...", formationName)
	createFormationReq := fixtures.FixCreateFormationRequest(formationName)
	executeGQLRequest(t, ctx, createFormationReq, formationName, accountTenantID)
	t.Logf("Successfully created formation: %s", formationName)

	defer func() {
		t.Logf("Deleting formation with name: %s...", formationName)
		deleteRequest := fixtures.FixDeleteFormationRequest(formationName)
		executeGQLRequest(t, ctx, deleteRequest, formationName, accountTenantID)
		t.Logf("Successfully deleted formation with name: %s...", formationName)
	}()

	t.Logf("Assign application to formation %s", formationName)
	assignToFormation(t, ctx, app.ID, "APPLICATION", formationName, accountTenantID)
	t.Logf("Successfully assigned application to formation %s", formationName)

	defer func() {
		t.Logf("Unassign application from formation %s", formationName)
		unassignFromFormation(t, ctx, app.ID, "APPLICATION", formationName, accountTenantID)
		t.Logf("Successfully unassigned application from formation %s", formationName)
	}()

	// Self register runtime
	runtimeInput := graphql.RuntimeRegisterInput{
		Name:        "selfRegisterRuntime",
		Description: ptr.String("selfRegisterRuntime-description"),
		Labels:      graphql.Labels{conf.SubscriptionConfig.SelfRegDistinguishLabelKey: conf.SubscriptionConfig.SelfRegDistinguishLabelValue, tenantfetcher.RegionKey: conf.SubscriptionConfig.SelfRegRegion},
	}
	runtime := fixtures.RegisterRuntimeFromInputWithoutTenant(t, ctx, certSecuredGraphQLClient, &runtimeInput)
	defer fixtures.CleanupRuntimeWithoutTenant(t, ctx, certSecuredGraphQLClient, &runtime)
	require.NotEmpty(t, runtime.ID)
	strLbl, ok := runtime.Labels[conf.SubscriptionConfig.SelfRegisterLabelKey].(string)
	require.True(t, ok)
	require.Contains(t, strLbl, runtime.ID)

	// Verify that the label returned cannot be modified
	setLabelRequest := fixtures.FixSetRuntimeLabelRequest(runtime.ID, conf.SubscriptionConfig.SelfRegisterLabelKey, "value")
	label := graphql.Label{}
	err = testctx.Tc.RunOperationWithoutTenant(ctx, certSecuredGraphQLClient, setLabelRequest, &label)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("could not set unmodifiable label with key %s", conf.SubscriptionConfig.SelfRegisterLabelKey))

	labelDefinitions, err := fixtures.ListLabelDefinitionsWithinTenant(t, ctx, certSecuredGraphQLClient, accountTenantID)
	require.NoError(t, err)
	numOfScenarioLabelDefinitions := 0
	for _, ld := range labelDefinitions {
		if ld.Key == scenariosLabel {
			numOfScenarioLabelDefinitions++
		}
	}
	// the parent tenant should not see child label definitions
	require.Equal(t, 1, numOfScenarioLabelDefinitions)
}

func TestConsumerProviderFlow(stdT *testing.T) {
	t := testingx.NewT(stdT)
	t.Run("ConsumerProvider flow", func(t *testing.T) {
		ctx := context.Background()
		secondaryTenant := conf.TestConsumerAccountID
		subscriptionProviderSubaccountID := conf.TestProviderSubaccountID
		subscriptionConsumerSubaccountID := conf.TestConsumerSubaccountID
		subscriptionConsumerTenantID := conf.TestConsumerTenantID

		// Prepare provider external client certificate and secret and Build graphql director client configured with certificate
		providerClientKey, providerRawCertChain := certprovider.NewExternalCertFromConfig(t, ctx, conf.ExternalCertProviderConfig)
		directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, providerClientKey, providerRawCertChain, conf.SkipSSLValidation)

		runtimeInput := graphql.RuntimeRegisterInput{
			Name:        "providerRuntime",
			Description: ptr.String("providerRuntime-description"),
			Labels:      graphql.Labels{conf.SubscriptionConfig.SelfRegDistinguishLabelKey: conf.SubscriptionConfig.SelfRegDistinguishLabelValue, tenantfetcher.RegionKey: conf.SubscriptionConfig.SelfRegRegion},
		}

		runtime := fixtures.RegisterRuntimeFromInputWithoutTenant(t, ctx, directorCertSecuredClient, &runtimeInput)
		defer fixtures.CleanupRuntimeWithoutTenant(t, ctx, directorCertSecuredClient, &runtime)
		require.NotEmpty(t, runtime.ID)

		// Register application
		app, err := fixtures.RegisterApplication(t, ctx, certSecuredGraphQLClient, "testingApp", secondaryTenant)
		defer fixtures.CleanupApplication(t, ctx, certSecuredGraphQLClient, secondaryTenant, &app)
		require.NoError(t, err)
		require.NotEmpty(t, app.ID)

		// Register consumer application
		consumerApp, err := fixtures.RegisterApplication(t, ctx, certSecuredGraphQLClient, "consumerApp", secondaryTenant)
		defer fixtures.CleanupApplication(t, ctx, certSecuredGraphQLClient, secondaryTenant, &consumerApp)
		require.NoError(t, err)
		require.NotEmpty(t, consumerApp.ID)
		require.NotEmpty(t, consumerApp.Name)

		consumerFormationName := "consumer-test-scenario"
		t.Logf("Creating formation with name %s...", consumerFormationName)
		createFormationReq := fixtures.FixCreateFormationRequest(consumerFormationName)
		executeGQLRequest(t, ctx, createFormationReq, consumerFormationName, secondaryTenant)
		t.Logf("Successfully created formation: %s", consumerFormationName)

		defer func() {
			t.Logf("Deleting formation with name: %s...", consumerFormationName)
			deleteRequest := fixtures.FixDeleteFormationRequest(consumerFormationName)
			executeGQLRequest(t, ctx, deleteRequest, consumerFormationName, secondaryTenant)
			t.Logf("Successfully deleted formation with name: %s...", consumerFormationName)
		}()

		t.Logf("Assign application to formation %s", consumerFormationName)
		assignToFormation(t, ctx, consumerApp.ID, "APPLICATION", consumerFormationName, secondaryTenant)
		t.Logf("Successfully assigned application to formation %s", consumerFormationName)

		defer func() {
			t.Logf("Unassign application from formation %s", consumerFormationName)
			unassignFromFormation(t, ctx, consumerApp.ID, "APPLICATION", consumerFormationName, secondaryTenant)
			t.Logf("Successfully unassigned application from formation %s", consumerFormationName)
		}()

		t.Logf("Assign tenant %s to formation %s...", subscriptionConsumerSubaccountID, consumerFormationName)
		assignToFormation(t, ctx, subscriptionConsumerSubaccountID, "TENANT", consumerFormationName, secondaryTenant)
		t.Logf("Successfully assigned tenant %s to formation %s", subscriptionConsumerSubaccountID, consumerFormationName)

		defer func() {
			t.Logf("Unassign tenant %s from formation %s", subscriptionConsumerSubaccountID, consumerFormationName)
			unassignFromFormation(t, ctx, subscriptionConsumerSubaccountID, "TENANT", consumerFormationName, secondaryTenant)
			t.Logf("Successfully unassigned tenant %s to formation %s", subscriptionConsumerSubaccountID, consumerFormationName)
		}()

		selfRegLabelValue, ok := runtime.Labels[conf.SubscriptionConfig.SelfRegisterLabelKey].(string)
		require.True(t, ok)
		require.Contains(t, selfRegLabelValue, conf.SubscriptionConfig.SelfRegisterLabelValuePrefix+runtime.ID)

		httpClient := &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.SkipSSLValidation},
			},
		}

		depConfigureReq, err := http.NewRequest(http.MethodPost, conf.ExternalServicesMockBaseURL+"/v1/dependencies/configure", bytes.NewBuffer([]byte(selfRegLabelValue)))
		require.NoError(t, err)
		response, err := httpClient.Do(depConfigureReq)
		require.NoError(t, err)
		defer func() {
			if err := response.Body.Close(); err != nil {
				t.Logf("Could not close response body %s", err)
			}
		}()
		require.Equal(t, http.StatusOK, response.StatusCode)

		apiPath := fmt.Sprintf("/saas-manager/v1/application/tenants/%s/subscriptions", subscriptionConsumerTenantID)
		subscribeReq, err := http.NewRequest(http.MethodPost, conf.SubscriptionConfig.URL+apiPath, bytes.NewBuffer([]byte("{\"subscriptionParams\": {}}")))
		require.NoError(t, err)
		subscriptionToken := token.GetClientCredentialsToken(t, ctx, conf.SubscriptionConfig.TokenURL+conf.TokenPath, conf.SubscriptionConfig.ClientID, conf.SubscriptionConfig.ClientSecret, "tenantFetcherClaims")
		subscribeReq.Header.Add(subscription.AuthorizationHeader, fmt.Sprintf("Bearer %s", subscriptionToken))
		subscribeReq.Header.Add(subscription.ContentTypeHeader, subscription.ContentTypeApplicationJson)
		subscribeReq.Header.Add(conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionProviderSubaccountID)

		// unsubscribe request execution to ensure no resources/subscriptions are left unintentionally due to old unsubscribe failures or broken tests in the middle.
		// In case there isn't subscription it will fail-safe without error
		subscription.BuildAndExecuteUnsubscribeRequest(t, runtime.ID, runtime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID)

		t.Logf("Creating a subscription between consumer with subaccount id: %q and tenant id: %q, and provider with name: %q, id: %q and subaccount id: %q", subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, runtime.Name, runtime.ID, subscriptionProviderSubaccountID)
		resp, err := httpClient.Do(subscribeReq)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Could not close response body %s", err)
			}
		}()
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusAccepted, resp.StatusCode, fmt.Sprintf("actual status code %d is different from the expected one: %d. Reason: %v", resp.StatusCode, http.StatusAccepted, string(body)))

		defer subscription.BuildAndExecuteUnsubscribeRequest(t, runtime.ID, runtime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID)

		subJobStatusPath := resp.Header.Get(subscription.LocationHeader)
		require.NotEmpty(t, subJobStatusPath)
		subJobStatusURL := conf.SubscriptionConfig.URL + subJobStatusPath
		require.Eventually(t, func() bool {
			return subscription.GetSubscriptionJobStatus(t, httpClient, subJobStatusURL, subscriptionToken) == subscription.JobSucceededStatus
		}, subscription.EventuallyTimeout, subscription.EventuallyTick)
		t.Logf("Successfully created subscription between consumer with subaccount id: %q and tenant id: %q, and provider with name: %q, id: %q and subaccount id: %q", subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, runtime.Name, runtime.ID, subscriptionProviderSubaccountID)

		// After successful subscription from above we call the director component with "double authentication(token + certificate)" in order to test claims validation is successful
		consumerToken := token.GetUserToken(t, ctx, conf.ConsumerTokenURL+conf.TokenPath, conf.ProviderClientID, conf.ProviderClientSecret, conf.BasicUsername, conf.BasicPassword, "subscriptionClaims")
		headers := map[string][]string{subscription.AuthorizationHeader: {fmt.Sprintf("Bearer %s", consumerToken)}}

		t.Log("Calling director to verify claims validation is successful...")
		getRtmReq := fixtures.FixGetRuntimeRequest(runtime.ID)
		getRtmReq.Header = headers
		rtmExt := graphql.RuntimeExt{}

		err = testctx.Tc.RunOperationWithCustomTenant(ctx, directorCertSecuredClient, subscriptionProviderSubaccountID, getRtmReq, &rtmExt)
		require.NoError(t, err)
		require.Equal(t, runtime.ID, rtmExt.ID)
		require.Equal(t, runtimeInput.Name, rtmExt.Name)
		t.Log("Director claims validation was successful")

		// After successful subscription from above, the part of the code below prepare and execute a request to the ord service

		// HTTP client configured with certificate with patched subject, issued from cert-rotation job
		certHttpClient := CreateHttpClientWithCert(providerClientKey, providerRawCertChain, conf.SkipSSLValidation)

		// Make a request to the ORD service with http client containing certificate with provider information and token with the consumer data.
		t.Log("Getting consumer application using both provider and consumer credentials...")
		respBody := makeRequestWithHeaders(t, certHttpClient, conf.ORDExternalCertSecuredServiceURL+"/systemInstances?$format=json", headers)
		require.Equal(t, 1, len(gjson.Get(respBody, "value").Array()))
		require.Equal(t, consumerApp.Name, gjson.Get(respBody, "value.0.title").String())
		t.Log("Successfully fetched consumer application using both provider and consumer credentials")

		subscription.BuildAndExecuteUnsubscribeRequest(t, runtime.ID, runtime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID)

		t.Log("Validating no application is returned after successful unsubscription request...")
		respBody = makeRequestWithHeaders(t, certHttpClient, conf.ORDExternalCertSecuredServiceURL+"/systemInstances?$format=json", headers)
		require.Equal(t, 0, len(gjson.Get(respBody, "value").Array()))
		t.Log("Successfully validated no application is returned after successful unsubscription request")

		t.Log("Validating director returns error during claims validation after unsubscribe request is successfully executed...")
		err = testctx.Tc.RunOperationWithCustomTenant(ctx, directorCertSecuredClient, subscriptionProviderSubaccountID, getRtmReq, &rtmExt)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("Consumer's external tenant %s was not found as subscription record in the runtime context table for the runtime in the provider tenant", subscriptionConsumerSubaccountID))
		t.Log("Successfully validated an error is returned during claims validation after unsubscribe request")
	})
}

func assignToFormation(t *testing.T, ctx context.Context, objectID, objectType, formationName, tenantID string) {
	assignReq := fixtures.FixAssignFormationRequest(objectID, objectType, formationName)
	executeGQLRequest(t, ctx, assignReq, formationName, tenantID)
}

func unassignFromFormation(t *testing.T, ctx context.Context, objectID, objectType, formationName, tenantID string) {
	unassignReq := fixtures.FixUnassignFormationRequest(objectID, objectType, formationName)
	executeGQLRequest(t, ctx, unassignReq, formationName, tenantID)
}

func executeGQLRequest(t *testing.T, ctx context.Context, gqlRequest *gcli.Request, formationName, tenantID string) {
	var formation graphql.Formation
	err := testctx.Tc.RunOperationWithCustomTenant(ctx, certSecuredGraphQLClient, tenantID, gqlRequest, &formation)
	require.NoError(t, err)
	require.Equal(t, formationName, formation.Name)
}

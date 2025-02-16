package model_test

import (
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestEventAPIDefinitionInput_ToEventAPIDefinition(t *testing.T) {
	// GIVEN
	id := "foo"
	bndlID := "bar"
	appID := "baz"
	desc := "Sample"
	name := "sample"
	group := "sampleGroup"

	testCases := []struct {
		Name     string
		Input    *model.EventDefinitionInput
		Expected *model.EventDefinition
	}{
		{
			Name: "All properties given",
			Input: &model.EventDefinitionInput{
				Name:        name,
				Description: &desc,
				Group:       &group,
			},
			Expected: &model.EventDefinition{
				ApplicationID: appID,
				Name:          name,
				Description:   &desc,
				Group:         &group,
				BaseEntity: &model.BaseEntity{
					ID:    id,
					Ready: true,
				},
			},
		},
		{
			Name:     "Nil",
			Input:    nil,
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// WHEN
			result := testCase.Input.ToEventDefinitionWithinBundle(id, appID, bndlID, 0)

			// then
			assert.Equal(t, testCase.Expected, result)
		})
	}
}

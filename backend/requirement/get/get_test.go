package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/jackstenglein/chess-dojo-scheduler/backend/api"
	"github.com/jackstenglein/chess-dojo-scheduler/backend/database"
)

const testRequirementId = "test-requirement-id"

var goldenRequirement = &database.Requirement{
	Id:          testRequirementId,
	Status:      database.Active,
	Category:    "Test Category",
	Name:        "Test Name",
	Description: "Test Description",
	Counts: map[database.DojoCohort]int{
		"2300-2400": 25,
		"2400+":     50,
	},
	NumberOfCohorts:   2,
	UnitScore:         1,
	VideoUrls:         []string{"Test Video Urls"},
	PositionUrls:      []string{"Test Position"},
	ScoreboardDisplay: database.Hidden,
	UpdatedAt:         "2023-03-01T12:00:00Z",
	SortPriority:      "0.0.0",
}

func getEvent(testName, id string) api.Request {
	return api.Request{
		RequestContext: events.APIGatewayProxyRequestContext{
			RequestID: testName,
		},
		PathParameters: map[string]string{
			"id": id,
		},
	}
}

func TestGetRequirement(t *testing.T) {
	ctx := context.Background()

	table := []struct {
		name            string
		requirementId   string
		wantCode        int
		wantErr         bool
		wantRequirement *database.Requirement
	}{
		{
			name:     "MissingId",
			wantCode: 400,
			wantErr:  true,
		},
		{
			name:            "SuccessfulRequest",
			requirementId:   testRequirementId,
			wantCode:        200,
			wantRequirement: goldenRequirement,
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			event := getEvent(tc.name, tc.requirementId)
			got, err := Handler(ctx, event)

			if err != nil {
				t.Errorf("GetRequirement(%v) got err: %v", event, err)
			}

			if got.StatusCode != tc.wantCode {
				t.Errorf("GetRequirement(%v) response: %v", event, got)
				t.Fatalf("GetRequirement(%v) got status: %d; want status: %d", event, got.StatusCode, tc.wantCode)
			}

			if !tc.wantErr {
				gotReq := &database.Requirement{}
				json.Unmarshal([]byte(got.Body), gotReq)

				if diff := cmp.Diff(tc.wantRequirement, gotReq); diff != "" {
					t.Errorf("GetRequirement(%v) diff (-want +got):\n%s", event, diff)
				}
			}
		})
	}
}

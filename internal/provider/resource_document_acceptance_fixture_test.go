package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
)

const (
	documentAcceptanceTeamQuery = `query getDocumentAcceptanceTeam {
  viewer {
    teams(first: 1) {
      nodes {
        id
      }
    }
  }
}
`

	documentAcceptanceProjectCreateMutation = `mutation createDocumentAcceptanceProject(
  $name: String!,
  $teamId: String!
) {
  projectCreate(input: {name: $name, teamIds: [$teamId]}) {
    success
    project {
      id
    }
  }
}
`

	documentAcceptanceIssueCreateMutation = `mutation createDocumentAcceptanceIssue(
  $title: String!,
  $teamId: String!,
  $projectId: String!
) {
  issueCreate(input: {title: $title, teamId: $teamId, projectId: $projectId}) {
    success
    issue {
      id
    }
  }
}
`

	documentAcceptanceIssueDeleteMutation = `mutation deleteDocumentAcceptanceIssue($id: String!) {
  issueDelete(id: $id) {
    success
  }
}
`

	documentAcceptanceProjectDeleteMutation = `mutation deleteDocumentAcceptanceProject($id: String!) {
  projectDelete(id: $id) {
    success
  }
}
`
)

func TestDocumentAcceptanceFixtureSetup(t *testing.T) {
	t.Parallel()

	client := &documentAcceptanceFixtureClient{
		responses: map[string]string{
			"getDocumentAcceptanceTeam": `{
				"viewer": {
					"teams": {
						"nodes": [{"id": "team-id"}]
					}
				}
			}`,
			"createDocumentAcceptanceProject": `{
				"projectCreate": {
					"success": true,
					"project": {"id": "project-id"}
				}
			}`,
			"createDocumentAcceptanceIssue": `{
				"issueCreate": {
					"success": true,
					"issue": {"id": "issue-id"}
				}
			}`,
		},
	}
	fixture := &documentAcceptanceFixture{}

	err := fixture.setup(context.Background(), client, "unique fixture")
	if err != nil {
		t.Fatalf("unexpected fixture setup error: %v", err)
	}

	if fixture.teamID != "team-id" {
		t.Errorf("expected team ID team-id, got %q", fixture.teamID)
	}
	if fixture.projectID != "project-id" {
		t.Errorf("expected project ID project-id, got %q", fixture.projectID)
	}
	if fixture.issueID != "issue-id" {
		t.Errorf("expected issue ID issue-id, got %q", fixture.issueID)
	}

	expectedOperations := []string{
		"getDocumentAcceptanceTeam",
		"createDocumentAcceptanceProject",
		"createDocumentAcceptanceIssue",
	}
	if !reflect.DeepEqual(client.operations, expectedOperations) {
		t.Fatalf("expected operations %v, got %v", expectedOperations, client.operations)
	}

	expectedVariables := map[string]map[string]any{
		"getDocumentAcceptanceTeam": nil,
		"createDocumentAcceptanceProject": {
			"name":   "unique fixture project",
			"teamId": "team-id",
		},
		"createDocumentAcceptanceIssue": {
			"title":     "unique fixture issue",
			"teamId":    "team-id",
			"projectId": "project-id",
		},
	}
	if !reflect.DeepEqual(client.variables, expectedVariables) {
		t.Fatalf("expected variables %#v, got %#v", expectedVariables, client.variables)
	}
}

func TestDocumentAcceptanceFixtureSetupRequiresTeam(t *testing.T) {
	t.Parallel()

	client := &documentAcceptanceFixtureClient{
		responses: map[string]string{
			"getDocumentAcceptanceTeam": `{
				"viewer": {"teams": {"nodes": []}}
			}`,
		},
	}
	fixture := &documentAcceptanceFixture{}

	err := fixture.setup(context.Background(), client, "unique fixture")
	if err == nil || !strings.Contains(err.Error(), "team") {
		t.Fatalf("expected missing-team error, got %v", err)
	}

	expectedOperations := []string{"getDocumentAcceptanceTeam"}
	if !reflect.DeepEqual(client.operations, expectedOperations) {
		t.Fatalf("expected operations %v, got %v", expectedOperations, client.operations)
	}
}

func TestDocumentAcceptanceFixtureSetupRetainsReturnedIDsOnFailure(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		responses         map[string]string
		expectedProjectID string
		expectedIssueID   string
		expectedError     string
	}{
		"project mutation reports failure": {
			responses: map[string]string{
				"getDocumentAcceptanceTeam": `{
					"viewer": {"teams": {"nodes": [{"id": "team-id"}]}}
				}`,
				"createDocumentAcceptanceProject": `{
					"projectCreate": {
						"success": false,
						"project": {"id": "project-id"}
					}
				}`,
			},
			expectedProjectID: "project-id",
			expectedError:     "project",
		},
		"issue mutation reports failure": {
			responses: map[string]string{
				"getDocumentAcceptanceTeam": `{
					"viewer": {"teams": {"nodes": [{"id": "team-id"}]}}
				}`,
				"createDocumentAcceptanceProject": `{
					"projectCreate": {
						"success": true,
						"project": {"id": "project-id"}
					}
				}`,
				"createDocumentAcceptanceIssue": `{
					"issueCreate": {
						"success": false,
						"issue": {"id": "issue-id"}
					}
				}`,
			},
			expectedProjectID: "project-id",
			expectedIssueID:   "issue-id",
			expectedError:     "issue",
		},
		"project response omits entity": {
			responses: map[string]string{
				"getDocumentAcceptanceTeam": `{
					"viewer": {"teams": {"nodes": [{"id": "team-id"}]}}
				}`,
				"createDocumentAcceptanceProject": `{
					"projectCreate": {"success": true}
				}`,
			},
			expectedError: "project",
		},
		"issue response omits entity": {
			responses: map[string]string{
				"getDocumentAcceptanceTeam": `{
					"viewer": {"teams": {"nodes": [{"id": "team-id"}]}}
				}`,
				"createDocumentAcceptanceProject": `{
					"projectCreate": {
						"success": true,
						"project": {"id": "project-id"}
					}
				}`,
				"createDocumentAcceptanceIssue": `{
					"issueCreate": {"success": true}
				}`,
			},
			expectedProjectID: "project-id",
			expectedError:     "issue",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &documentAcceptanceFixtureClient{responses: testCase.responses}
			fixture := &documentAcceptanceFixture{}

			err := fixture.setup(context.Background(), client, "unique fixture")
			if err == nil || !strings.Contains(err.Error(), testCase.expectedError) {
				t.Fatalf("expected %s setup error, got %v", testCase.expectedError, err)
			}
			if fixture.projectID != testCase.expectedProjectID {
				t.Errorf("expected project ID %q, got %q", testCase.expectedProjectID, fixture.projectID)
			}
			if fixture.issueID != testCase.expectedIssueID {
				t.Errorf("expected issue ID %q, got %q", testCase.expectedIssueID, fixture.issueID)
			}
		})
	}
}

func TestDocumentAcceptanceFixtureCleanupAttemptsIssueThenProject(t *testing.T) {
	t.Parallel()

	client := &documentAcceptanceFixtureClient{
		responses: map[string]string{
			"deleteDocumentAcceptanceProject": `{
				"projectDelete": {"success": true}
			}`,
		},
		errors: map[string]error{
			"deleteDocumentAcceptanceIssue": errors.New("issue cleanup failed"),
		},
	}
	fixture := &documentAcceptanceFixture{
		projectID: "project-id",
		issueID:   "issue-id",
	}

	err := fixture.cleanup(context.Background(), client)
	if err == nil || !strings.Contains(err.Error(), "issue cleanup failed") {
		t.Fatalf("expected issue cleanup error, got %v", err)
	}

	expectedOperations := []string{
		"deleteDocumentAcceptanceIssue",
		"deleteDocumentAcceptanceProject",
	}
	if !reflect.DeepEqual(client.operations, expectedOperations) {
		t.Fatalf("expected cleanup operations %v, got %v", expectedOperations, client.operations)
	}
}

func TestDocumentAcceptanceFixtureCleanupSkipsUncreatedEntities(t *testing.T) {
	t.Parallel()

	client := &documentAcceptanceFixtureClient{
		responses: map[string]string{
			"deleteDocumentAcceptanceProject": `{
				"projectDelete": {"success": true}
			}`,
		},
	}
	fixture := &documentAcceptanceFixture{projectID: "project-id"}

	if err := fixture.cleanup(context.Background(), client); err != nil {
		t.Fatalf("unexpected cleanup error: %v", err)
	}

	expectedOperations := []string{"deleteDocumentAcceptanceProject"}
	if !reflect.DeepEqual(client.operations, expectedOperations) {
		t.Fatalf("expected cleanup operations %v, got %v", expectedOperations, client.operations)
	}
}

type documentAcceptanceFixture struct {
	teamID    string
	projectID string
	issueID   string
}

// setup creates the project and issue parents used by the document acceptance test.
func (f *documentAcceptanceFixture) setup(
	ctx context.Context,
	client graphql.Client,
	name string,
) error {
	var teamResponse struct {
		Viewer struct {
			Teams struct {
				Nodes []struct {
					ID string `json:"id"`
				} `json:"nodes"`
			} `json:"teams"`
		} `json:"viewer"`
	}
	err := makeDocumentAcceptanceRequest(
		ctx,
		client,
		"getDocumentAcceptanceTeam",
		documentAcceptanceTeamQuery,
		nil,
		&teamResponse,
	)
	if err != nil {
		return fmt.Errorf("find a team for the document acceptance test: %w", err)
	}
	if len(teamResponse.Viewer.Teams.Nodes) == 0 {
		return errors.New("the Linear token owner must belong to a team for the document acceptance test")
	}
	f.teamID = teamResponse.Viewer.Teams.Nodes[0].ID

	var projectResponse struct {
		ProjectCreate struct {
			Success bool `json:"success"`
			Project struct {
				ID string `json:"id"`
			} `json:"project"`
		} `json:"projectCreate"`
	}
	err = makeDocumentAcceptanceRequest(
		ctx,
		client,
		"createDocumentAcceptanceProject",
		documentAcceptanceProjectCreateMutation,
		map[string]any{
			"name":   name + " project",
			"teamId": f.teamID,
		},
		&projectResponse,
	)
	f.projectID = projectResponse.ProjectCreate.Project.ID
	if err != nil {
		return fmt.Errorf("create the document acceptance test project: %w", err)
	}
	if !projectResponse.ProjectCreate.Success || f.projectID == "" {
		return errors.New("Linear reported that it was unable to create the document acceptance test project")
	}

	var issueResponse struct {
		IssueCreate struct {
			Success bool `json:"success"`
			Issue   struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"issueCreate"`
	}
	err = makeDocumentAcceptanceRequest(
		ctx,
		client,
		"createDocumentAcceptanceIssue",
		documentAcceptanceIssueCreateMutation,
		map[string]any{
			"title":     name + " issue",
			"teamId":    f.teamID,
			"projectId": f.projectID,
		},
		&issueResponse,
	)
	f.issueID = issueResponse.IssueCreate.Issue.ID
	if err != nil {
		return fmt.Errorf("create the document acceptance test issue: %w", err)
	}
	if !issueResponse.IssueCreate.Success || f.issueID == "" {
		return errors.New("Linear reported that it was unable to create the document acceptance test issue")
	}

	return nil
}

// cleanup trashes any fixture entities that were created, preserving dependency order.
func (f *documentAcceptanceFixture) cleanup(ctx context.Context, client graphql.Client) error {
	var cleanupErrors []string

	if f.issueID != "" {
		var response struct {
			IssueDelete struct {
				Success bool `json:"success"`
			} `json:"issueDelete"`
		}
		err := makeDocumentAcceptanceRequest(
			ctx,
			client,
			"deleteDocumentAcceptanceIssue",
			documentAcceptanceIssueDeleteMutation,
			map[string]any{"id": f.issueID},
			&response,
		)
		if err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Sprintf("delete document acceptance test issue: %v", err))
		} else if !response.IssueDelete.Success {
			cleanupErrors = append(cleanupErrors, "Linear reported that it was unable to delete the document acceptance test issue")
		} else {
			f.issueID = ""
		}
	}

	if f.projectID != "" {
		var response struct {
			ProjectDelete struct {
				Success bool `json:"success"`
			} `json:"projectDelete"`
		}
		err := makeDocumentAcceptanceRequest(
			ctx,
			client,
			"deleteDocumentAcceptanceProject",
			documentAcceptanceProjectDeleteMutation,
			map[string]any{"id": f.projectID},
			&response,
		)
		if err != nil {
			cleanupErrors = append(cleanupErrors, fmt.Sprintf("delete document acceptance test project: %v", err))
		} else if !response.ProjectDelete.Success {
			cleanupErrors = append(cleanupErrors, "Linear reported that it was unable to delete the document acceptance test project")
		} else {
			f.projectID = ""
		}
	}

	if len(cleanupErrors) > 0 {
		return errors.New(strings.Join(cleanupErrors, "; "))
	}
	return nil
}

func makeDocumentAcceptanceRequest(
	ctx context.Context,
	client graphql.Client,
	operationName string,
	query string,
	variables map[string]any,
	data any,
) error {
	return client.MakeRequest(
		ctx,
		&graphql.Request{
			OpName:    operationName,
			Query:     query,
			Variables: variables,
		},
		&graphql.Response{Data: data},
	)
}

func newDocumentAcceptanceClient() graphql.Client {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &authedTransport{
			token:   os.Getenv(envVarName),
			wrapped: http.DefaultTransport,
		},
	}
	return graphql.NewClient("https://api.linear.app/graphql", httpClient)
}

type documentAcceptanceFixtureClient struct {
	responses  map[string]string
	errors     map[string]error
	operations []string
	variables  map[string]map[string]any
}

func (c *documentAcceptanceFixtureClient) MakeRequest(
	_ context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	c.operations = append(c.operations, request.OpName)
	if c.variables == nil {
		c.variables = make(map[string]map[string]any)
	}

	encodedVariables, err := json.Marshal(request.Variables)
	if err != nil {
		return err
	}
	var variables map[string]any
	if err := json.Unmarshal(encodedVariables, &variables); err != nil {
		return err
	}
	c.variables[request.OpName] = variables

	if err := c.errors[request.OpName]; err != nil {
		return err
	}
	payload, ok := c.responses[request.OpName]
	if !ok {
		return errors.New("unexpected GraphQL operation: " + request.OpName)
	}
	return json.Unmarshal([]byte(payload), response.Data)
}

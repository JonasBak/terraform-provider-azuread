package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AuthenticationStrengthPoliciesClient performs operations on AuthenticationStrengthPolicies.
type AuthenticationStrengthPoliciesClient struct {
	BaseClient Client
}

// NewAuthenticationStrengthPoliciesClient returns a new AuthenticationStrengthPoliciesClient
func NewAuthenticationStrengthPoliciesClient(tenantId string) *AuthenticationStrengthPoliciesClient {
	return &AuthenticationStrengthPoliciesClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// Create creates a new AuthenticationStrengthPolicies.
func (c *AuthenticationStrengthPoliciesClient) Create(ctx context.Context, authenticationStrengthPolicy AuthenticationStrengthPolicy) (*AuthenticationStrengthPolicy, int, error) {
	var status int
	body, err := json.Marshal(authenticationStrengthPolicy)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/policies/authenticationStrengthPolicies",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationStrengthPoliciesClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newAuthenticationStrengthPolicy AuthenticationStrengthPolicy
	if err := json.Unmarshal(respBody, &newAuthenticationStrengthPolicy); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newAuthenticationStrengthPolicy, status, nil
}

// Get retrieves a AuthenticationStrengthPolicy.
func (c *AuthenticationStrengthPoliciesClient) Get(ctx context.Context, id string) (*AuthenticationStrengthPolicy, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/policies/authenticationStrengthPolicies/%s", id),
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("AuthenticationStrengthPoliciesClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var authenticationStrengthPolicy AuthenticationStrengthPolicy
	if err := json.Unmarshal(respBody, &authenticationStrengthPolicy); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &authenticationStrengthPolicy, status, nil
}

// Update amends an existing AuthenticationStrengthPolicy.
func (c *AuthenticationStrengthPoliciesClient) Update(ctx context.Context, id string, authenticationStrengthPolicy AuthenticationStrengthPolicy) (int, error) {
	var status int

	body, err := json.Marshal(authenticationStrengthPolicy)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/policies/authenticationStrengthPolicies/%s", id),
		},
	})
	if err != nil {
		return status, fmt.Errorf("AuthenticationStrengthPoliciesClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}

// Delete removes a AuthenticationStrengthPolicy.
func (c *AuthenticationStrengthPoliciesClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/policies/authenticationStrengthPolicies/%s/$ref", id),
		},
	})
	if err != nil {
		return status, fmt.Errorf("AuthenticationStrengthPoliciesClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

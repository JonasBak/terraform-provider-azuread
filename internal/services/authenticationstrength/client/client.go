package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	AuthenticationStrengthPoliciesClient *msgraph.AuthenticationStrengthPoliciesClient
}

func NewClient(o *common.ClientOptions) *Client {
	authenticationStrengthPoliciesClient := msgraph.NewAuthenticationStrengthPoliciesClient(o.TenantID)
	o.ConfigureClient(&authenticationStrengthPoliciesClient.BaseClient)

	return &Client{
		AuthenticationStrengthPoliciesClient: authenticationStrengthPoliciesClient,
	}
}

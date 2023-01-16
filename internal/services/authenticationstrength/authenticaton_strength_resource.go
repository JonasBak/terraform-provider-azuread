package authenticationstrength

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/manicminer/hamilton/msgraph"
	// "github.com/manicminer/hamilton/odata"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/helpers"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func authenticationStrengthPolicyResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: authenticationStrengthPolicyResourceCreate,
		ReadContext:   authenticationStrengthPolicyResourceRead,
		UpdateContext: authenticationStrengthPolicyResourceUpdate,
		DeleteContext: authenticationStrengthPolicyResourceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			if _, err := uuid.ParseUUID(id); err != nil {
				return fmt.Errorf("specified ID (%q) is not valid: %s", id, err)
			}
			return nil
		}),

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"description": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"allowed_combinations": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func authenticationStrengthPolicyResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AuthenticationStrength.AuthenticationStrengthPoliciesClient

	display_name := d.Get("display_name").(string)
	description := d.Get("description").(string)
	allowed_combinations := []string{}
	for _, v := range d.Get("allowed_combinations").([]interface{}) {
		allowed_combinations = append(allowed_combinations, v.(string))
	}

	properties := msgraph.AuthenticationStrengthPolicy{
		DisplayName:         &display_name,
		Description:         &description,
		AllowedCombinations: &allowed_combinations,
	}

	policy, _, err := client.Create(ctx, properties)
	if err != nil {
		return tf.ErrorDiagF(err, "Could not create authentication strength policy")
	}

	if policy.ID == nil || *policy.ID == "" {
		return tf.ErrorDiagF(errors.New("Bad API response"), "Object ID returned for authentication strength policy is nil/empty")
	}

	d.SetId(*policy.ID)

	return authenticationStrengthPolicyResourceRead(ctx, d, meta)
}

func authenticationStrengthPolicyResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AuthenticationStrength.AuthenticationStrengthPoliciesClient

	display_name := d.Get("display_name").(string)
	description := d.Get("description").(string)
	allowed_combinations := []string{}
	for _, v := range d.Get("allowed_combinations").([]interface{}) {
		allowed_combinations = append(allowed_combinations, v.(string))
	}

	properties := msgraph.AuthenticationStrengthPolicy{
		DisplayName:         &display_name,
		Description:         &description,
	}

	if _, err := client.Update(ctx, d.Id(), properties); err != nil {
		return tf.ErrorDiagF(err, "Could not update conditional access policy with ID: %q", d.Id())
	}

	return nil
}

func authenticationStrengthPolicyResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AuthenticationStrength.AuthenticationStrengthPoliciesClient

	policy, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Authentication Strength Policy with Object ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return tf.ErrorDiagPathF(err, "id", "Retrieving Authentication Strength Policy with object ID %q", d.Id())
	}

	tf.Set(d, "display_name", policy.DisplayName)
	tf.Set(d, "description", policy.Description)
	tf.Set(d, "allowed_combinations", policy.AllowedCombinations)

	return nil
}

func authenticationStrengthPolicyResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).AuthenticationStrength.AuthenticationStrengthPoliciesClient
	policyId := d.Id()

	_, status, err := client.Get(ctx, policyId)
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Authentication Strength Policy with ID %q already deleted", policyId)
			return nil
		}

		return tf.ErrorDiagPathF(err, "id", "Retrieving authentication strength policy with ID %q", policyId)
	}

	status, err = client.Delete(ctx, policyId)
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Deleting authentication strength policy with ID %q, got status %d", policyId, status)
	}

	if err := helpers.WaitForDeletion(ctx, func(ctx context.Context) (*bool, error) {
		client.BaseClient.DisableRetries = true
		if _, status, err := client.Get(ctx, policyId); err != nil {
			if status == http.StatusNotFound {
				return utils.Bool(false), nil
			}
			return nil, err
		}
		return utils.Bool(true), nil
	}); err != nil {
		return tf.ErrorDiagF(err, "Waiting for deletion of authentication strength policy with ID %q", policyId)
	}

	return nil
}

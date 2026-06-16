// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverdomainlists"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnssecurityrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverSecurityRuleModel struct {
	Name                     string            `tfschema:"name"`
	DnsResolverPolicyId      string            `tfschema:"dns_resolver_policy_id"`
	Location                 string            `tfschema:"location"`
	Action                   string            `tfschema:"action"`
	DnsResolverDomainListIds []string          `tfschema:"dns_resolver_domain_list_ids"`
	Priority                 int64             `tfschema:"priority"`
	Enabled                  bool              `tfschema:"enabled"`
	Tags                     map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverSecurityRuleResource struct{}

var (
	_ sdk.ResourceWithIdentity = PrivateDNSResolverSecurityRuleResource{}
	_ sdk.ResourceWithUpdate   = PrivateDNSResolverSecurityRuleResource{}
)

func (r PrivateDNSResolverSecurityRuleResource) Identity() resourceids.ResourceId {
	return &dnssecurityrules.DnsSecurityRuleId{}
}

func (r PrivateDNSResolverSecurityRuleResource) ResourceType() string {
	return "azurerm_private_dns_resolver_security_rule"
}

func (r PrivateDNSResolverSecurityRuleResource) ModelObject() interface{} {
	return &PrivateDNSResolverSecurityRuleModel{}
}

func (r PrivateDNSResolverSecurityRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnssecurityrules.ValidateDnsSecurityRuleID
}

func (r PrivateDNSResolverSecurityRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dns_resolver_policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: dnsresolverpolicies.ValidateDnsResolverPolicyID,
		},

		"location": commonschema.Location(),

		"action": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				dnssecurityrules.PossibleValuesForActionType(),
				false,
			),
		},

		"dns_resolver_domain_list_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: dnsresolverdomainlists.ValidateDnsResolverDomainListID,
			},
		},

		"priority": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(100, 65000),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverSecurityRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverSecurityRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverSecurityRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.DnsSecurityRulesClient

			policyId, err := dnsresolverpolicies.ParseDnsResolverPolicyID(model.DnsResolverPolicyId)
			if err != nil {
				return err
			}

			id := dnssecurityrules.NewDnsSecurityRuleID(policyId.SubscriptionId, policyId.ResourceGroupName, policyId.DnsResolverPolicyName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			action := dnssecurityrules.ActionType(model.Action)
			state := dnssecurityrules.DnsSecurityRuleStateEnabled
			if !model.Enabled {
				state = dnssecurityrules.DnsSecurityRuleStateDisabled
			}

			properties := dnssecurityrules.DnsSecurityRule{
				Location: location.Normalize(model.Location),
				Properties: dnssecurityrules.DnsSecurityRuleProperties{
					Action: dnssecurityrules.DnsSecurityRuleAction{
						ActionType: &action,
					},
					DnsResolverDomainLists: expandDnsSecurityRuleDomainListReferences(model.DnsResolverDomainListIds),
					DnsSecurityRuleState:   &state,
					Priority:               model.Priority,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, dnssecurityrules.CreateOrUpdateOperationOptions{}, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverSecurityRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsSecurityRulesClient

			id, err := dnssecurityrules.ParseDnsSecurityRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverSecurityRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			patch := dnssecurityrules.DnsSecurityRulePatch{
				Properties: &dnssecurityrules.DnsSecurityRulePatchProperties{},
			}

			if metadata.ResourceData.HasChange("action") {
				action := dnssecurityrules.ActionType(model.Action)
				patch.Properties.Action = &dnssecurityrules.DnsSecurityRuleAction{
					ActionType: &action,
				}
			}

			if metadata.ResourceData.HasChange("dns_resolver_domain_list_ids") {
				lists := expandDnsSecurityRuleDomainListReferences(model.DnsResolverDomainListIds)
				patch.Properties.DnsResolverDomainLists = &lists
			}

			if metadata.ResourceData.HasChange("enabled") {
				state := dnssecurityrules.DnsSecurityRuleStateEnabled
				if !model.Enabled {
					state = dnssecurityrules.DnsSecurityRuleStateDisabled
				}
				patch.Properties.DnsSecurityRuleState = &state
			}

			if metadata.ResourceData.HasChange("priority") {
				priority := model.Priority
				patch.Properties.Priority = &priority
			}

			if metadata.ResourceData.HasChange("tags") {
				patch.Tags = &model.Tags
			}

			if err := client.UpdateThenPoll(ctx, *id, patch, dnssecurityrules.UpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverSecurityRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsSecurityRulesClient

			id, err := dnssecurityrules.ParseDnsSecurityRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			return r.flatten(metadata, id, model)
		},
	}
}

func (r PrivateDNSResolverSecurityRuleResource) flatten(metadata sdk.ResourceMetaData, id *dnssecurityrules.DnsSecurityRuleId, model *dnssecurityrules.DnsSecurityRule) error {
	state := PrivateDNSResolverSecurityRuleModel{
		Name:                id.DnsSecurityRuleName,
		DnsResolverPolicyId: dnsresolverpolicies.NewDnsResolverPolicyID(id.SubscriptionId, id.ResourceGroupName, id.DnsResolverPolicyName).ID(),
		Location:            location.Normalize(model.Location),
		Priority:            model.Properties.Priority,
		DnsResolverDomainListIds: flattenDnsSecurityRuleDomainListReferences(model.Properties.DnsResolverDomainLists),
	}

	if action := model.Properties.Action.ActionType; action != nil {
		state.Action = string(*action)
	}

	state.Enabled = true
	if ruleState := model.Properties.DnsSecurityRuleState; ruleState != nil {
		state.Enabled = *ruleState == dnssecurityrules.DnsSecurityRuleStateEnabled
	}

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r PrivateDNSResolverSecurityRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsSecurityRulesClient

			id, err := dnssecurityrules.ParseDnsSecurityRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, dnssecurityrules.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDnsSecurityRuleDomainListReferences(input []string) []dnssecurityrules.SubResource {
	output := make([]dnssecurityrules.SubResource, 0, len(input))
	for _, id := range input {
		output = append(output, dnssecurityrules.SubResource{
			Id: id,
		})
	}
	return output
}

func flattenDnsSecurityRuleDomainListReferences(input []dnssecurityrules.SubResource) []string {
	output := make([]string, 0, len(input))
	for _, item := range input {
		// API responses can return mixed-case provider segments; normalise so
		// the value round-trips cleanly via Terraform state.
		if parsed, err := dnsresolverdomainlists.ParseDnsResolverDomainListIDInsensitively(item.Id); err == nil {
			output = append(output, parsed.ID())
			continue
		}
		output = append(output, item.Id)
	}
	return output
}

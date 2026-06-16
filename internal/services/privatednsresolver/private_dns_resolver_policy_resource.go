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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverPolicyModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverPolicyResource struct{}

var (
	_ sdk.ResourceWithIdentity = PrivateDNSResolverPolicyResource{}
	_ sdk.ResourceWithUpdate   = PrivateDNSResolverPolicyResource{}
)

func (r PrivateDNSResolverPolicyResource) Identity() resourceids.ResourceId {
	return &dnsresolverpolicies.DnsResolverPolicyId{}
}

func (r PrivateDNSResolverPolicyResource) ResourceType() string {
	return "azurerm_private_dns_resolver_policy"
}

func (r PrivateDNSResolverPolicyResource) ModelObject() interface{} {
	return &PrivateDNSResolverPolicyModel{}
}

func (r PrivateDNSResolverPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnsresolverpolicies.ValidateDnsResolverPolicyID
}

func (r PrivateDNSResolverPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.DnsResolverPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := dnsresolverpolicies.NewDnsResolverPolicyID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := dnsresolverpolicies.DnsResolverPolicy{
				Location:   location.Normalize(model.Location),
				Properties: &dnsresolverpolicies.DnsResolverPolicyProperties{},
				Tags:       &model.Tags,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, dnsresolverpolicies.CreateOrUpdateOperationOptions{}, metadata.SetIDAndIdentityCallback(&id)); err != nil {
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

func (r PrivateDNSResolverPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverPoliciesClient

			id, err := dnsresolverpolicies.ParseDnsResolverPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				patch := dnsresolverpolicies.DnsResolverPolicyPatch{
					Tags: &model.Tags,
				}

				if err := client.UpdateThenPoll(ctx, *id, patch, dnsresolverpolicies.UpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverPoliciesClient

			id, err := dnsresolverpolicies.ParseDnsResolverPolicyID(metadata.ResourceData.Id())
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

func (r PrivateDNSResolverPolicyResource) flatten(metadata sdk.ResourceMetaData, id *dnsresolverpolicies.DnsResolverPolicyId, model *dnsresolverpolicies.DnsResolverPolicy) error {
	state := PrivateDNSResolverPolicyModel{
		Name:              id.DnsResolverPolicyName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.Normalize(model.Location),
	}

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r PrivateDNSResolverPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverPoliciesClient

			id, err := dnsresolverpolicies.ParseDnsResolverPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, dnsresolverpolicies.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

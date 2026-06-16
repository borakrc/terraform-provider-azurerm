// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicyvirtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverPolicyVirtualNetworkLinkModel struct {
	Name                string            `tfschema:"name"`
	DnsResolverPolicyId string            `tfschema:"dns_resolver_policy_id"`
	Location            string            `tfschema:"location"`
	VirtualNetworkId    string            `tfschema:"virtual_network_id"`
	Tags                map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverPolicyVirtualNetworkLinkResource struct{}

var (
	_ sdk.ResourceWithIdentity = PrivateDNSResolverPolicyVirtualNetworkLinkResource{}
	_ sdk.ResourceWithUpdate   = PrivateDNSResolverPolicyVirtualNetworkLinkResource{}
)

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Identity() resourceids.ResourceId {
	return &dnsresolverpolicyvirtualnetworklinks.DnsResolverPolicyVirtualNetworkLinkId{}
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) ResourceType() string {
	return "azurerm_private_dns_resolver_policy_virtual_network_link"
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) ModelObject() interface{} {
	return &PrivateDNSResolverPolicyVirtualNetworkLinkModel{}
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnsresolverpolicyvirtualnetworklinks.ValidateDnsResolverPolicyVirtualNetworkLinkID
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Arguments() map[string]*pluginsdk.Schema {
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

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverPolicyVirtualNetworkLinkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.DnsResolverPolicyVirtualNetworkLinksClient

			policyId, err := dnsresolverpolicies.ParseDnsResolverPolicyID(model.DnsResolverPolicyId)
			if err != nil {
				return err
			}

			id := dnsresolverpolicyvirtualnetworklinks.NewDnsResolverPolicyVirtualNetworkLinkID(policyId.SubscriptionId, policyId.ResourceGroupName, policyId.DnsResolverPolicyName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := dnsresolverpolicyvirtualnetworklinks.DnsResolverPolicyVirtualNetworkLink{
				Location: location.Normalize(model.Location),
				Properties: dnsresolverpolicyvirtualnetworklinks.DnsResolverPolicyVirtualNetworkLinkProperties{
					VirtualNetwork: dnsresolverpolicyvirtualnetworklinks.SubResource{
						Id: model.VirtualNetworkId,
					},
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, dnsresolverpolicyvirtualnetworklinks.CreateOrUpdateOperationOptions{}, metadata.SetIDAndIdentityCallback(&id)); err != nil {
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

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverPolicyVirtualNetworkLinksClient

			id, err := dnsresolverpolicyvirtualnetworklinks.ParseDnsResolverPolicyVirtualNetworkLinkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverPolicyVirtualNetworkLinkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				patch := dnsresolverpolicyvirtualnetworklinks.DnsResolverPolicyVirtualNetworkLinkPatch{
					Tags: &model.Tags,
				}

				if err := client.UpdateThenPoll(ctx, *id, patch, dnsresolverpolicyvirtualnetworklinks.UpdateOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverPolicyVirtualNetworkLinksClient

			id, err := dnsresolverpolicyvirtualnetworklinks.ParseDnsResolverPolicyVirtualNetworkLinkID(metadata.ResourceData.Id())
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

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) flatten(metadata sdk.ResourceMetaData, id *dnsresolverpolicyvirtualnetworklinks.DnsResolverPolicyVirtualNetworkLinkId, model *dnsresolverpolicyvirtualnetworklinks.DnsResolverPolicyVirtualNetworkLink) error {
	state := PrivateDNSResolverPolicyVirtualNetworkLinkModel{
		Name:                id.VirtualNetworkLinkName,
		DnsResolverPolicyId: dnsresolverpolicies.NewDnsResolverPolicyID(id.SubscriptionId, id.ResourceGroupName, id.DnsResolverPolicyName).ID(),
		Location:            location.Normalize(model.Location),
		VirtualNetworkId:    model.Properties.VirtualNetwork.Id,
	}

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverPolicyVirtualNetworkLinksClient

			id, err := dnsresolverpolicyvirtualnetworklinks.ParseDnsResolverPolicyVirtualNetworkLinkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, dnsresolverpolicyvirtualnetworklinks.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

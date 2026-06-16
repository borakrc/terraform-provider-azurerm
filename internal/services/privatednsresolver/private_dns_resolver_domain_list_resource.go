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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverDomainListModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Domains           []string          `tfschema:"domains"`
	Tags              map[string]string `tfschema:"tags"`
}

type PrivateDNSResolverDomainListResource struct{}

var (
	_ sdk.ResourceWithIdentity = PrivateDNSResolverDomainListResource{}
	_ sdk.ResourceWithUpdate   = PrivateDNSResolverDomainListResource{}
)

func (r PrivateDNSResolverDomainListResource) Identity() resourceids.ResourceId {
	return &dnsresolverdomainlists.DnsResolverDomainListId{}
}

func (r PrivateDNSResolverDomainListResource) ResourceType() string {
	return "azurerm_private_dns_resolver_domain_list"
}

func (r PrivateDNSResolverDomainListResource) ModelObject() interface{} {
	return &PrivateDNSResolverDomainListModel{}
}

func (r PrivateDNSResolverDomainListResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnsresolverdomainlists.ValidateDnsResolverDomainListID
}

func (r PrivateDNSResolverDomainListResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"domains": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverDomainListResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverDomainListResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverDomainListModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.DnsResolverDomainListsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := dnsresolverdomainlists.NewDnsResolverDomainListID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			domains := model.Domains

			properties := dnsresolverdomainlists.DnsResolverDomainList{
				Location: location.Normalize(model.Location),
				Properties: &dnsresolverdomainlists.DnsResolverDomainListProperties{
					Domains: &domains,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, dnsresolverdomainlists.CreateOrUpdateOperationOptions{}, metadata.SetIDAndIdentityCallback(&id)); err != nil {
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

func (r PrivateDNSResolverDomainListResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverDomainListsClient

			id, err := dnsresolverdomainlists.ParseDnsResolverDomainListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverDomainListModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			patch := dnsresolverdomainlists.DnsResolverDomainListPatch{}

			if metadata.ResourceData.HasChange("domains") {
				domains := model.Domains
				patch.Properties = &dnsresolverdomainlists.DnsResolverDomainListPatchProperties{
					Domains: &domains,
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				patch.Tags = &model.Tags
			}

			if err := client.UpdateThenPoll(ctx, *id, patch, dnsresolverdomainlists.UpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverDomainListResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverDomainListsClient

			id, err := dnsresolverdomainlists.ParseDnsResolverDomainListID(metadata.ResourceData.Id())
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

func (r PrivateDNSResolverDomainListResource) flatten(metadata sdk.ResourceMetaData, id *dnsresolverdomainlists.DnsResolverDomainListId, model *dnsresolverdomainlists.DnsResolverDomainList) error {
	state := PrivateDNSResolverDomainListModel{
		Name:              id.DnsResolverDomainListName,
		ResourceGroupName: id.ResourceGroupName,
		Location:          location.Normalize(model.Location),
	}

	if props := model.Properties; props != nil {
		if props.Domains != nil {
			state.Domains = *props.Domains
		}
	}

	if model.Tags != nil {
		state.Tags = *model.Tags
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r PrivateDNSResolverDomainListResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolverDomainListsClient

			id, err := dnsresolverdomainlists.ParseDnsResolverDomainListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, dnsresolverdomainlists.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

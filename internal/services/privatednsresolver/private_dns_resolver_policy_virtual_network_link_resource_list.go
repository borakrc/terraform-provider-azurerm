// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicyvirtualnetworklinks"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDNSResolverPolicyVirtualNetworkLinkListResource struct{}

type PrivateDNSResolverPolicyVirtualNetworkLinkListModel struct {
	DnsResolverPolicyId types.String `tfsdk:"dns_resolver_policy_id"`
}

var _ sdk.FrameworkListWrappedResource = new(PrivateDNSResolverPolicyVirtualNetworkLinkListResource)

func (PrivateDNSResolverPolicyVirtualNetworkLinkListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = PrivateDNSResolverPolicyVirtualNetworkLinkResource{}.ResourceType()
}

func (PrivateDNSResolverPolicyVirtualNetworkLinkListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(PrivateDNSResolverPolicyVirtualNetworkLinkResource{})
}

func (PrivateDNSResolverPolicyVirtualNetworkLinkListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dns_resolver_policy_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: dnsresolverpolicies.ValidateDnsResolverPolicyID},
				},
			},
		},
	}
}

func (PrivateDNSResolverPolicyVirtualNetworkLinkListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.PrivateDnsResolver.DnsResolverPolicyVirtualNetworkLinksClient

	var data PrivateDNSResolverPolicyVirtualNetworkLinkListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	resourceType := PrivateDNSResolverPolicyVirtualNetworkLinkResource{}

	policyId, err := dnsresolverpolicies.ParseDnsResolverPolicyID(data.DnsResolverPolicyId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing DNS Resolver Policy ID for `%s`", resourceType.ResourceType()), err)
		return
	}

	resp, err := client.ListComplete(ctx, dnsresolverpolicyvirtualnetworklinks.NewDnsResolverPolicyID(policyId.SubscriptionId, policyId.ResourceGroupName, policyId.DnsResolverPolicyName), dnsresolverpolicyvirtualnetworklinks.DefaultListOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", resourceType.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := dnsresolverpolicyvirtualnetworklinks.ParseDnsResolverPolicyVirtualNetworkLinkIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing DNS Resolver Policy Virtual Network Link ID", err)
				return
			}

			meta := sdk.NewResourceMetaData(metadata.Client, resourceType)
			meta.SetID(id)

			if err := resourceType.flatten(meta, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", resourceType.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}

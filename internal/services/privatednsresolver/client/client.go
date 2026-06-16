// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverdomainlists"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicyvirtualnetworklinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnssecurityrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DnsForwardingRulesetsClient                *dnsforwardingrulesets.DnsForwardingRulesetsClient
	DnsResolverDomainListsClient               *dnsresolverdomainlists.DnsResolverDomainListsClient
	DnsResolverPoliciesClient                  *dnsresolverpolicies.DnsResolverPoliciesClient
	DnsResolverPolicyVirtualNetworkLinksClient *dnsresolverpolicyvirtualnetworklinks.DnsResolverPolicyVirtualNetworkLinksClient
	DnsResolversClient                         *dnsresolvers.DnsResolversClient
	DnsSecurityRulesClient                     *dnssecurityrules.DnsSecurityRulesClient
	ForwardingRulesClient                      *forwardingrules.ForwardingRulesClient
	InboundEndpointsClient                     *inboundendpoints.InboundEndpointsClient
	OutboundEndpointsClient                    *outboundendpoints.OutboundEndpointsClient
	VirtualNetworkLinksClient                  *virtualnetworklinks.VirtualNetworkLinksClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	dnsForwardingRulesetsClient, err := dnsforwardingrulesets.NewDnsForwardingRulesetsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(dnsForwardingRulesetsClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsForwardingRulesetsClient client: %+v", err)
	}

	dnsResolverDomainListsClient, err := dnsresolverdomainlists.NewDnsResolverDomainListsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsResolverDomainListsClient client: %+v", err)
	}
	o.Configure(dnsResolverDomainListsClient.Client, o.Authorizers.ResourceManager)

	dnsResolverPoliciesClient, err := dnsresolverpolicies.NewDnsResolverPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsResolverPoliciesClient client: %+v", err)
	}
	o.Configure(dnsResolverPoliciesClient.Client, o.Authorizers.ResourceManager)

	dnsResolverPolicyVirtualNetworkLinksClient, err := dnsresolverpolicyvirtualnetworklinks.NewDnsResolverPolicyVirtualNetworkLinksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsResolverPolicyVirtualNetworkLinksClient client: %+v", err)
	}
	o.Configure(dnsResolverPolicyVirtualNetworkLinksClient.Client, o.Authorizers.ResourceManager)

	dnsResolversClient, err := dnsresolvers.NewDnsResolversClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(dnsResolversClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsResolversClient client: %+v", err)
	}

	dnsSecurityRulesClient, err := dnssecurityrules.NewDnsSecurityRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsSecurityRulesClient client: %+v", err)
	}
	o.Configure(dnsSecurityRulesClient.Client, o.Authorizers.ResourceManager)

	forwardingRulesClient, err := forwardingrules.NewForwardingRulesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(forwardingRulesClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ForwardingRulesClient client: %+v", err)
	}

	inboundEndpointsClient, err := inboundendpoints.NewInboundEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(inboundEndpointsClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building InboundEndpointsClient client: %+v", err)
	}

	outboundEndpointsClient, err := outboundendpoints.NewOutboundEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(outboundEndpointsClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building OutboundEndpointsClient client: %+v", err)
	}

	virtualNetworkLinksClient, err := virtualnetworklinks.NewVirtualNetworkLinksClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(virtualNetworkLinksClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualNetworkLinksClient client: %+v", err)
	}

	return &Client{
		DnsForwardingRulesetsClient:                dnsForwardingRulesetsClient,
		DnsResolverDomainListsClient:               dnsResolverDomainListsClient,
		DnsResolverPoliciesClient:                  dnsResolverPoliciesClient,
		DnsResolverPolicyVirtualNetworkLinksClient: dnsResolverPolicyVirtualNetworkLinksClient,
		DnsResolversClient:                         dnsResolversClient,
		DnsSecurityRulesClient:                     dnsSecurityRulesClient,
		ForwardingRulesClient:                      forwardingRulesClient,
		InboundEndpointsClient:                     inboundEndpointsClient,
		OutboundEndpointsClient:                    outboundEndpointsClient,
		VirtualNetworkLinksClient:                  virtualNetworkLinksClient,
	}, nil
}

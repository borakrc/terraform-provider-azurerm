// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverpolicyvirtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDNSResolverPolicyVirtualNetworkLinkResource struct{}

func TestAccPrivateDNSResolverPolicyVirtualNetworkLink_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_policy_virtual_network_link", "test")
	r := PrivateDNSResolverPolicyVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateDNSResolverPolicyVirtualNetworkLink_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_policy_virtual_network_link", "test")
	r := PrivateDNSResolverPolicyVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPrivateDNSResolverPolicyVirtualNetworkLink_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_policy_virtual_network_link", "test")
	r := PrivateDNSResolverPolicyVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateDNSResolverPolicyVirtualNetworkLink_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_policy_virtual_network_link", "test")
	r := PrivateDNSResolverPolicyVirtualNetworkLinkResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dnsresolverpolicyvirtualnetworklinks.ParseDnsResolverPolicyVirtualNetworkLinkID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PrivateDnsResolver.DnsResolverPolicyVirtualNetworkLinksClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_private_dns_resolver_policy" "test" {
  name                = "acctest-drp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_policy_virtual_network_link" "test" {
  name                   = "acctest-drpvnl-%d"
  dns_resolver_policy_id = azurerm_private_dns_resolver_policy.test.id
  location               = azurerm_resource_group.test.location
  virtual_network_id     = azurerm_virtual_network.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_policy_virtual_network_link" "import" {
  name                   = azurerm_private_dns_resolver_policy_virtual_network_link.test.name
  dns_resolver_policy_id = azurerm_private_dns_resolver_policy_virtual_network_link.test.dns_resolver_policy_id
  location               = azurerm_private_dns_resolver_policy_virtual_network_link.test.location
  virtual_network_id     = azurerm_private_dns_resolver_policy_virtual_network_link.test.virtual_network_id
}
`, r.basic(data))
}

func (r PrivateDNSResolverPolicyVirtualNetworkLinkResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_policy_virtual_network_link" "test" {
  name                   = "acctest-drpvnl-%d"
  dns_resolver_policy_id = azurerm_private_dns_resolver_policy.test.id
  location               = azurerm_resource_group.test.location
  virtual_network_id     = azurerm_virtual_network.test.id

  tags = {
    environment = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

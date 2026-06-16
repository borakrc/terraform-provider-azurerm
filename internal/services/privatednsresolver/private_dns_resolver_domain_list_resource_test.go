// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnsresolverdomainlists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDNSResolverDomainListResource struct{}

func TestAccPrivateDNSResolverDomainList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_domain_list", "test")
	r := PrivateDNSResolverDomainListResource{}
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

func TestAccPrivateDNSResolverDomainList_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_domain_list", "test")
	r := PrivateDNSResolverDomainListResource{}
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

func TestAccPrivateDNSResolverDomainList_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_domain_list", "test")
	r := PrivateDNSResolverDomainListResource{}
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

func TestAccPrivateDNSResolverDomainList_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_domain_list", "test")
	r := PrivateDNSResolverDomainListResource{}
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

func (r PrivateDNSResolverDomainListResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dnsresolverdomainlists.ParseDnsResolverDomainListID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PrivateDnsResolver.DnsResolverDomainListsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r PrivateDNSResolverDomainListResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PrivateDNSResolverDomainListResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_domain_list" "test" {
  name                = "acctest-drdl-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  domains = ["contoso.com."]
}
`, r.template(data), data.RandomInteger)
}

func (r PrivateDNSResolverDomainListResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_domain_list" "import" {
  name                = azurerm_private_dns_resolver_domain_list.test.name
  resource_group_name = azurerm_private_dns_resolver_domain_list.test.resource_group_name
  location            = azurerm_private_dns_resolver_domain_list.test.location

  domains = ["contoso.com."]
}
`, r.basic(data))
}

func (r PrivateDNSResolverDomainListResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_domain_list" "test" {
  name                = "acctest-drdl-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  domains = [
    "contoso.com.",
    "adatum.com.",
  ]

  tags = {
    environment = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

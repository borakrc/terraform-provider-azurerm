// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2025-05-01/dnssecurityrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateDNSResolverSecurityRuleResource struct{}

func TestAccPrivateDNSResolverSecurityRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_security_rule", "test")
	r := PrivateDNSResolverSecurityRuleResource{}
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

func TestAccPrivateDNSResolverSecurityRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_security_rule", "test")
	r := PrivateDNSResolverSecurityRuleResource{}
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

func TestAccPrivateDNSResolverSecurityRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_security_rule", "test")
	r := PrivateDNSResolverSecurityRuleResource{}
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

func TestAccPrivateDNSResolverSecurityRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_security_rule", "test")
	r := PrivateDNSResolverSecurityRuleResource{}
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

func (r PrivateDNSResolverSecurityRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dnssecurityrules.ParseDnsSecurityRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PrivateDnsResolver.DnsSecurityRulesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r PrivateDNSResolverSecurityRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_private_dns_resolver_policy" "test" {
  name                = "acctest-drp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_private_dns_resolver_domain_list" "test" {
  name                = "acctest-drdl-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  domains = ["contoso.com."]
}

resource "azurerm_private_dns_resolver_domain_list" "test2" {
  name                = "acctest-drdl2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  domains = ["adatum.com."]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PrivateDNSResolverSecurityRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_security_rule" "test" {
  name                         = "acctest-drsr-%d"
  dns_resolver_policy_id       = azurerm_private_dns_resolver_policy.test.id
  location                     = azurerm_resource_group.test.location
  action                       = "Alert"
  priority                     = 100
  dns_resolver_domain_list_ids = [azurerm_private_dns_resolver_domain_list.test.id]
}
`, r.template(data), data.RandomInteger)
}

func (r PrivateDNSResolverSecurityRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_security_rule" "import" {
  name                         = azurerm_private_dns_resolver_security_rule.test.name
  dns_resolver_policy_id       = azurerm_private_dns_resolver_security_rule.test.dns_resolver_policy_id
  location                     = azurerm_private_dns_resolver_security_rule.test.location
  action                       = azurerm_private_dns_resolver_security_rule.test.action
  priority                     = azurerm_private_dns_resolver_security_rule.test.priority
  dns_resolver_domain_list_ids = azurerm_private_dns_resolver_security_rule.test.dns_resolver_domain_list_ids
}
`, r.basic(data))
}

func (r PrivateDNSResolverSecurityRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver_security_rule" "test" {
  name                   = "acctest-drsr-%d"
  dns_resolver_policy_id = azurerm_private_dns_resolver_policy.test.id
  location               = azurerm_resource_group.test.location
  action                 = "Block"
  priority               = 200
  enabled                = false
  dns_resolver_domain_list_ids = [
    azurerm_private_dns_resolver_domain_list.test.id,
    azurerm_private_dns_resolver_domain_list.test2.id,
  ]

  tags = {
    environment = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

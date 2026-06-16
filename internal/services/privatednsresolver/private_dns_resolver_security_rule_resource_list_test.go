// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccPrivateDNSResolverSecurityRule_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_security_rule", "testlist")
	r := PrivateDNSResolverSecurityRuleResource{}
	resourceName := fmt.Sprintf("acctest-drsr-%d", data.RandomInteger)
	policyName := fmt.Sprintf("acctest-drp-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctest-rg-%d", data.RandomInteger)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.basic(data)},
			{
				Query:  true,
				Config: r.listQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_private_dns_resolver_security_rule.list", 1),
					querycheck.ExpectIdentity("azurerm_private_dns_resolver_security_rule.list", map[string]knownvalue.Check{
						"resource_group_name":      knownvalue.StringExact(resourceGroupName),
						"dns_resolver_policy_name": knownvalue.StringExact(policyName),
						"dns_security_rule_name":   knownvalue.StringExact(resourceName),
						"subscription_id":          knownvalue.StringExact(data.Subscriptions.Primary),
					}),
				},
			},
		},
	})
}

func (r PrivateDNSResolverSecurityRuleResource) listQuery() string {
	return `
list "azurerm_private_dns_resolver_security_rule" "list" {
  provider = azurerm
  config {
    dns_resolver_policy_id = azurerm_private_dns_resolver_policy.test.id
  }
}
`
}

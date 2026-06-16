// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccPrivateDNSResolverSecurityRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_resolver_security_rule", "test")
	r := PrivateDNSResolverSecurityRuleResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":          {},
		"resource_group_name":      {},
		"dns_resolver_policy_name": {},
		"dns_security_rule_name":   {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_private_dns_resolver_security_rule.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_private_dns_resolver_security_rule.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_private_dns_resolver_security_rule.test", tfjsonpath.New("dns_security_rule_name"), tfjsonpath.New("name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}

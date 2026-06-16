---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_security_rule"
description: |-
  Lists DNS Resolver Security Rule resources.
---

# List resource: azurerm_private_dns_resolver_security_rule

Lists DNS Resolver Security Rule resources.

## Example Usage

### List DNS Resolver Security Rules for a DNS Resolver Policy

```hcl
list "azurerm_private_dns_resolver_security_rule" "example" {
  provider = azurerm
  config {
    dns_resolver_policy_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/dnsResolverPolicies/example-policy"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `dns_resolver_policy_id` - (Required) The ID of the DNS Resolver Policy to query.

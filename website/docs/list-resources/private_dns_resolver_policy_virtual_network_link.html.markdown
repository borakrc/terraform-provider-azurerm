---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_policy_virtual_network_link"
description: |-
  Lists DNS Resolver Policy Virtual Network Link resources.
---

# List resource: azurerm_private_dns_resolver_policy_virtual_network_link

Lists DNS Resolver Policy Virtual Network Link resources.

## Example Usage

### List Virtual Network Links for a DNS Resolver Policy

```hcl
list "azurerm_private_dns_resolver_policy_virtual_network_link" "example" {
  provider = azurerm
  config {
    dns_resolver_policy_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/dnsResolverPolicies/example-policy"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `dns_resolver_policy_id` - (Required) The ID of the DNS Resolver Policy to query.

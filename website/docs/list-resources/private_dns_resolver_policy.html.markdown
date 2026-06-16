---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_policy"
description: |-
  Lists DNS Resolver Policy resources.
---

# List resource: azurerm_private_dns_resolver_policy

Lists DNS Resolver Policy resources.

## Example Usage

### List all DNS Resolver Policies in the subscription

```hcl
list "azurerm_private_dns_resolver_policy" "example" {
  provider = azurerm
  config {}
}
```

### List all DNS Resolver Policies in a specific resource group

```hcl
list "azurerm_private_dns_resolver_policy" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.

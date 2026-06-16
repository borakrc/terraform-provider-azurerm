---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_policy"
description: |-
  Manages a DNS Resolver Policy (the parent resource for an Azure DNS security policy).
---

# azurerm_private_dns_resolver_policy

Manages a DNS Resolver Policy (the parent resource for an Azure DNS security policy).

A DNS Resolver Policy is the top-level container in the Azure DNS security policy hierarchy. It is regional and provides the scope under which `azurerm_private_dns_resolver_security_rule` and `azurerm_private_dns_resolver_policy_virtual_network_link` resources are created. For more information see the [Microsoft documentation](https://learn.microsoft.com/azure/dns/dns-security-policy).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_private_dns_resolver_policy" "example" {
  name                = "example-policy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this DNS Resolver Policy. Changing this forces a new DNS Resolver Policy to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the DNS Resolver Policy should exist. Changing this forces a new DNS Resolver Policy to be created.

* `location` - (Required) Specifies the Azure Region where the DNS Resolver Policy should exist. Changing this forces a new DNS Resolver Policy to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the DNS Resolver Policy.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the DNS Resolver Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS Resolver Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the DNS Resolver Policy.
* `update` - (Defaults to 30 minutes) Used when updating the DNS Resolver Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the DNS Resolver Policy.

## Import

DNS Resolver Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsResolverPolicies/policy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-05-01

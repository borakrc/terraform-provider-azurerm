---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_policy_virtual_network_link"
description: |-
  Manages a virtual network link for an Azure DNS security policy.
---

# azurerm_private_dns_resolver_policy_virtual_network_link

Manages a virtual network link for an Azure DNS security policy.

A virtual network link associates an `azurerm_private_dns_resolver_policy` to a Virtual Network so the policy's DNS traffic rules apply to that VNet. A single VNet can be linked to at most one DNS security policy in the same region. For more information see the [Microsoft documentation](https://learn.microsoft.com/azure/dns/dns-security-policy).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_private_dns_resolver_policy" "example" {
  name                = "example-policy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_private_dns_resolver_policy_virtual_network_link" "example" {
  name                   = "example-link"
  dns_resolver_policy_id = azurerm_private_dns_resolver_policy.example.id
  location               = azurerm_resource_group.example.location
  virtual_network_id     = azurerm_virtual_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this link. Changing this forces a new link to be created.

* `dns_resolver_policy_id` - (Required) The ID of the DNS Resolver Policy that the Virtual Network is linked to. Changing this forces a new link to be created.

* `location` - (Required) Specifies the Azure Region where the link should exist. Must be the same region as the parent policy. Changing this forces a new link to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network that is linked to the DNS Resolver Policy. The Virtual Network must be in the same region as the policy. Changing this forces a new link to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the link.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the DNS Resolver Policy Virtual Network Link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the link.
* `read` - (Defaults to 5 minutes) Used when retrieving the link.
* `update` - (Defaults to 30 minutes) Used when updating the link.
* `delete` - (Defaults to 30 minutes) Used when deleting the link.

## Import

Links can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_policy_virtual_network_link.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsResolverPolicies/policy1/virtualNetworkLinks/link1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-05-01

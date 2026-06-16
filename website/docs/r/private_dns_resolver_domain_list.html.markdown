---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_domain_list"
description: |-
  Manages a DNS Resolver Domain List used by Azure DNS security policy rules.
---

# azurerm_private_dns_resolver_domain_list

Manages a DNS Resolver Domain List used by Azure DNS security policy rules.

A domain list is a regional list of DNS domains that is referenced by one or more `azurerm_private_dns_resolver_security_rule` resources to scope the rule's action. For more information see the [Microsoft documentation](https://learn.microsoft.com/azure/dns/dns-security-policy).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_private_dns_resolver_domain_list" "example" {
  name                = "example-domain-list"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  domains = [
    "contoso.com.",
    "adatum.com.",
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this DNS Resolver Domain List. Changing this forces a new DNS Resolver Domain List to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the DNS Resolver Domain List should exist. Changing this forces a new DNS Resolver Domain List to be created.

* `location` - (Required) Specifies the Azure Region where the DNS Resolver Domain List should exist. Changing this forces a new DNS Resolver Domain List to be created.

* `domains` - (Required) A set of DNS domains to include in the list. Must contain at least one entry. Each entry should be a fully qualified domain name (e.g. `contoso.com.`). Wildcard entries are accepted by the Azure API.

* `tags` - (Optional) A mapping of tags which should be assigned to the DNS Resolver Domain List.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the DNS Resolver Domain List.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS Resolver Domain List.
* `read` - (Defaults to 5 minutes) Used when retrieving the DNS Resolver Domain List.
* `update` - (Defaults to 30 minutes) Used when updating the DNS Resolver Domain List.
* `delete` - (Defaults to 30 minutes) Used when deleting the DNS Resolver Domain List.

## Import

DNS Resolver Domain Lists can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_domain_list.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsResolverDomainLists/list1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-05-01

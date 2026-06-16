---
subcategory: "Private DNS Resolver"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_resolver_security_rule"
description: |-
  Manages a DNS traffic rule for an Azure DNS security policy.
---

# azurerm_private_dns_resolver_security_rule

Manages a DNS traffic rule for an Azure DNS security policy.

A DNS security rule decides how DNS queries that match a set of `azurerm_private_dns_resolver_domain_list` resources are handled (`Allow`, `Block` or `Alert`). Rules are evaluated in priority order under a parent `azurerm_private_dns_resolver_policy`. For more information see the [Microsoft documentation](https://learn.microsoft.com/azure/dns/dns-security-policy).

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

resource "azurerm_private_dns_resolver_domain_list" "example" {
  name                = "example-domain-list"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  domains = [
    "contoso.com.",
  ]
}

resource "azurerm_private_dns_resolver_security_rule" "example" {
  name                         = "example-rule"
  dns_resolver_policy_id       = azurerm_private_dns_resolver_policy.example.id
  location                     = azurerm_resource_group.example.location
  action                       = "Block"
  priority                     = 100
  dns_resolver_domain_list_ids = [azurerm_private_dns_resolver_domain_list.example.id]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this DNS Security Rule. Changing this forces a new rule to be created.

* `dns_resolver_policy_id` - (Required) The ID of the DNS Resolver Policy that owns this rule. Changing this forces a new rule to be created.

* `location` - (Required) Specifies the Azure Region where the rule should exist. Must be the same region as the parent policy. Changing this forces a new rule to be created.

* `action` - (Required) The action taken when the rule matches a DNS query. Possible values are `Allow`, `Alert` and `Block`.

* `priority` - (Required) The priority of the rule within the parent policy. Lower numbers are evaluated first. Must be between `100` and `65000`.

* `dns_resolver_domain_list_ids` - (Required) A set of `azurerm_private_dns_resolver_domain_list` IDs that the rule applies to. Must contain at least one entry.

---

* `enabled` - (Optional) Whether the rule is enabled. Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the DNS Security Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the rule.
* `update` - (Defaults to 30 minutes) Used when updating the rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the rule.

## Import

DNS Security Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_resolver_security_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Network/dnsResolverPolicies/policy1/dnsSecurityRules/rule1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-05-01

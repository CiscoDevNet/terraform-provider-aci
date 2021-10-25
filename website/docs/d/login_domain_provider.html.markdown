---
layout: "aci"
page_title: "ACI: aci_login_domain_provider"
sidebar_current: "docs-aci-data-source-login_domain_provider"
description: |-
  Data source for ACI Login Domain Provider
---

# aci_login_domain_provider #
Data source for ACI Login Domain Provider


## API Information ##
* `Class` - aaaProviderRef

## Example Usage ##
```hcl
data "aci_login_domain_provider" "example" {
  parent_dn  = aci_duo_provider_group.example.id
  name  = "example"
}
```

## Argument Reference ##
* `parent_dn` - (Required) Distinguished name of parent.
* `name` - (Required) Name of object Login Domain Provider.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Login Domain Provider.
* `annotation` - (Optional) Annotation of object Login Domain Provider.
* `name_alias` - (Optional) Name Alias of object Login Domain Provider.
* `description` - (Optional) Description of object Login Domain Provider.
* `order` - (Optional) Order in which Providers are Tried. The relative priority in which the AAA provider will be contacted within the provider group. 

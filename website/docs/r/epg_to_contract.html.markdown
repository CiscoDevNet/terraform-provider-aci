---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_to_contract"
sidebar_current: "docs-aci-resource-epg_to_contract"
description: |-
  Manages ACI EPG to contract relationship.
---

# aci_epg_to_contract

Manages ACI EPG to contract relationship.

## Example Usage

```hcl
resource "aci_epg_to_contract" "example" {
  application_epg_dn = aci_application_epg.demo.id
  contract_dn        = aci_contract.demo_contract.id
  contract_type      = "provider"
  annotation         = "terraform"
  match_t            = "AtleastOne"
  prio               = "unspecified"
}
```

## Argument Reference

- `application_epg_dn` - (Required) Distinguished name of Parent EPG.
- `contract_dn` - (Required) Distinguished name of the Contract object to attach.
- `contract_type` - (Required) Type of the EPG to contract relationship object. Allowed values are "consumer" and "provider".
- `annotation` - (Optional) Annotation of the EPG to contract relationship object.
- `match_t` - (Optional) Matching criteria of the EPG to contract relationship object, only supported for `contract_type` "provider". Allowed values: "All", "AtleastOne", "AtmostOne", "None". Default value: "AtleastOne".
- `prio` - (Optional) Priority of the EPG to contract relationship object. Allowed values: "unspecified", "level1", "level2", "level3", "level4", "level5", "level6". Default value: "unspecified".

## Attribute Reference

- `id` - Attribute id set to the Dn of the provider/consumer contract.

## Importing ##

An existing EPG to contract relationship can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_epg_to_contract.example <Dn>
```

Starting in Terraform version 1.5, an existing EPG to contract relationship can be imported using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

 ```
 import {
    id = "<Dn>"
    to = aci_epg_to_contract.example
 }
 ```
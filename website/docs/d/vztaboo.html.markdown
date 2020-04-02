---
layout: "aci"
page_title: "ACI: aci_taboo_contract"
sidebar_current: "docs-aci-data-source-taboo_contract"
description: |-
  Data source for ACI Taboo Contract
---

# aci_taboo_contract #
Data source for ACI Taboo Contract

## Example Usage ##

```hcl
data "aci_taboo_contract" "example" {

  tenant_dn  = "${aci_tenant.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object taboo_contract.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Taboo Contract.
* `annotation` - (Optional) annotation for object taboo_contract.
* `name_alias` - (Optional) name_alias for object taboo_contract.

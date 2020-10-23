---
layout: "aci"
page_title: "ACI: aci_imported_contract"
sidebar_current: "docs-aci-data-source-imported_contract"
description: |-
  Data source for ACI Imported Contract
---

# aci_imported_contract #
Data source for ACI Imported Contract

## Example Usage ##

```hcl
data "aci_imported_contract" "example" {

  tenant_dn  = "${aci_tenant.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object imported_contract.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Imported Contract.
* `annotation` - (Optional) annotation for object imported_contract.
* `name_alias` - (Optional) name_alias for object imported_contract.

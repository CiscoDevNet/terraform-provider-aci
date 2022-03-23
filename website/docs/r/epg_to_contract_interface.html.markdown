---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_contract_interface"
sidebar_current: "docs-aci-resource-contract_interface"
description: |-
  Manages ACI Contract Interface
---

# aci_contract_interface #

Manages ACI Contract Interface

## API Information ##

* `Class` - fvRsConsIf
* `Distinguished Name` - uni/tn-{name}/ap-{name}/epg-{name}/rsconsIf-{contract_interface_name}

## GUI Information ##

* `Location` - Tenant -> Application Profiles -> Application EPGs -> Contracts


## Example Usage ##

```hcl
resource "aci_contract_interface" "example" {
  application_epg_dn  = aci_application_epg.example.id
  contract_interface_dn = aci_imported_contract.contract_interface.id
  prio = "unspecified"

}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of the parent ApplicationEPG object.
* `contract_interface_dn` - (Required) Distinguished name of the object Contract Interface.
* `annotation` - (Optional) Annotation of the object Contract Interface.
* `prio` - (Optional) prio.The contract interface priority. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.


## Importing ##

An existing ContractInterface can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_contract_interface.example <Dn>
```
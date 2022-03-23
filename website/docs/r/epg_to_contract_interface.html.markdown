---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_to_contract_interface"
sidebar_current: "docs-aci-resource-epg_to_contract_interface"
description: |-
  Manages ACI Contract Interface Relationship
---

# aci_epg_to_contract_interface #

Manages ACI Contract Interface Relationship

## API Information ##

* `Class` - fvRsConsIf
* `Distinguished Name` - uni/tn-{name}/ap-{name}/epg-{name}/rsconsIf-{contract_interface_name}

## GUI Information ##

* `Location` - Tenant -> Application Profiles -> Application EPGs -> Contracts


## Example Usage ##

```hcl
resource "aci_epg_to_contract_interface" "example" {
  application_epg_dn  = aci_application_epg.example.id
  contract_interface_dn = aci_imported_contract.contract_interface.id
  prio = "unspecified"

}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of the parent Application EPG object.
* `contract_interface_dn` - (Required) Distinguished name of the object Contract Interface object.
* `annotation` - (Optional) Annotation of the object Contract Interface Relationship object.
* `prio` - (Optional) The contract interface priority. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.


## Importing ##

An existing Contract Interface Relationship object can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_epg_to_contract_interface.example <Dn>
```
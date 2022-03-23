---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_contract_interface"
sidebar_current: "docs-aci-data-source-contract_interface"
description: |-
  Data source for ACI Contract Interface
---

# aci_contract_interface #

Data source for ACI Contract Interface


## API Information ##

* `Class` - fvRsConsIf
* `Distinguished Name` - uni/tn-{name}/ap-{name}/epg-{name}/rsconsIf-{contract_interface_name}

## GUI Information ##

* `Location` - Tenant -> Application Profiles -> Application EPGs -> Contracts



## Example Usage ##

```hcl
data "aci_contract_interface" "example" {
  application_epg_dn  = aci_application_epg.example.id
  contract_interface_dn = aci_imported_contract.contract_interface.id
}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `contract_interface_dn` - (Required) Distinguished name of object Contract Interface.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Contract Interface.
* `annotation` - (Optional) Annotation of object Contract Interface.
* `prio` - (Optional) prio. The contract interface priority.

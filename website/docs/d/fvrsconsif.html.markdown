---
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
* `Distinguished Name` - uni/tn-{name}/ap-{name}/epg-{name}/rsconsIf-{tnVzCPIfName}

## GUI Information ##

* `Location` - Tenant -> Application Profiles -> Application EPGs -> Contracts



## Example Usage ##

```hcl
data "aci_contract_interface" "example" {
  application_epg_dn  = aci_application_epg.example.id
  tnVzCPIfName  = "example"
}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tnVzCPIfName` - (Required) TnVzCPIfName of object Contract Interface.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Contract Interface.
* `annotation` - (Optional) Annotation of object Contract Interface.
* `name_alias` - (Optional) Name Alias of object Contract Interface.
* `prio` - (Optional) prio. The contract interface priority.

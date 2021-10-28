---
layout: "aci"
page_title: "ACI: aci_tacacs_accounting"
sidebar_current: "docs-aci-data-source-tacacs_accounting"
description: |-
  Data source for ACI TACACS Accounting
---

# aci_tacacs_accounting #

Data source for ACI TACACS Accoounting


## API Information ##

* `Class` - tacacsGroup
* `Distinguished Named` - uni/fabric/tacacsgroup-{name}

## GUI Information ##

* `Location` - Admin -> External Data Collectors -> Monitoring Destinations -> TACACS 



## Example Usage ##

```hcl
data "aci_tacacs_accounting" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object TACACS Accounting.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the TACACS Accounting.
* `annotation` - (Optional) Annotation of object TACACS Accounting.
* `name_alias` - (Optional) Name Alias of object TACACS Accounting.
* `description` - (Optional) Description of object TACACS Accounting.
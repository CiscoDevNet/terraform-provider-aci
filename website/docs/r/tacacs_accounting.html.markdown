---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_tacacs_accounting"
sidebar_current: "docs-aci-resource-tacacs_accounting"
description: |-
  Manages ACI TACACS Accounting
---

# aci_tacacs_accounting #

Manages ACI TACACS Accounting

## API Information ##

* `Class` - tacacsGroup
* `Distinguished Named` - uni/fabric/tacacsgroup-{name}

## GUI Information ##

* `Location` - Admin -> External Data Collectors -> Monitoring Destinations -> TACACS 


## Example Usage ##

```hcl
resource "aci_tacacs_accounting" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  name_alias  = "tacacs_accounting_alias"
  description = "From Terraform"
}
```

## Argument Reference ##


* `name` - (Required) Name of object TACACS Accounting. Type: String.
* `annotation` - (Optional) Annotation of object TACACS Accounting. Type: String.
* `name_alias` - (Optional) Name Alias of object TACACS Accounting. Type: String.
* `description` - (Optional) Description of object TACACS Accounting. Type: String.



## Importing ##

An existing TACACSAccounting can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tacacs_accounting.example <Dn>
```
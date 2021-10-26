---
layout: "aci"
page_title: "ACI: aci_tacacs_source"
sidebar_current: "docs-aci-data-source-tacacs_source"
description: |-
  Data source for ACI TACACS Source
---

# aci_tacacs_source #

Data source for ACI TACACS Source


## API Information ##

* `Class` - tacacsSrc


## Example Usage ##

```hcl
data "aci_tacacs_source" "example" {
  parent_dn  = parent_resource.example.id
  name  = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of parent object of TACACS Source.
* `name` - (Required) name of object TACACS Source.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the TACACS Source.
* `annotation` - (Optional) Annotation of object TACACS Source.
* `name_alias` - (Optional) Name Alias of object TACACS Source.
* `description` - (Optional) Description of object TACACS Source.
* `incl` - (Optional) Include Action. The information to include for the call home source.
* `min_sev` - (Optional) minSev. 

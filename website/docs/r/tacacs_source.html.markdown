---
layout: "aci"
page_title: "ACI: aci_tacacs_source"
sidebar_current: "docs-aci-resource-tacacs_source"
description: |-
  Manages ACI TACACS Source
---

# aci_tacacs_source #

Manages ACI TACACS Source

## API Information ##

* `Class` - tacacsSrc

## Example Usage ##

```hcl
resource "aci_tacacs_source" "example" {
  parent_dn   = parent_resource.example.id
  name        = "example"
  annotation  = "orchestrator:terraform"
  incl        = ["audit","session"]
  min_sev     = "info"
  name_alias  = "tacacs_source_alias"
  description = "From Terraform"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of parent object of TACACS Source.
* `name` - (Required) Name of object TACACS Source.
* `annotation` - (Optional) Annotation of object TACACS Source.
* `name_alias` - (Optional) Name Alias of object TACACS Source. Type: String.
* `description` - (Optional) Description of object TACACS Source. Type: String.
* `incl` - (Optional) Include Action. The information to include for the call home source. Allowed values are "audit", "events", "faults" and "session". Default value is ["audit","session"]. Type: List.
* `min_sev` - (Optional) minSev. Allowed values are "cleared", "critical", "info", "major", "minor" and "warning". Default value is "info". Type: String.
* `relation_tacacs_rs_dest_group` - (Optional) Represents the relation to a TACACS Destination Group (class tacacsGroup). Type: String.



## Importing ##

An existing TACACSSource can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tacacs_source.example <Dn>
```
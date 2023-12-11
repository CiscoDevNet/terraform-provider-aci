---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_tacacs_source"
sidebar_current: "docs-aci-resource-aci_tacacs_source"
description: |-
  Manages ACI TACACS Source
---

# aci_tacacs_source #

Manages ACI TACACS Source

## API Information ##

* `Class` - tacacsSrc
* `Supported Distinguished Name` - <br>
[1] uni/tn-{name}/monepg-{name}/tarepg-{scope}/tacacssrc-{name}<br>
[2] uni/infra/moninfra-{name}/tarinfra-{scope}/tacacssrc-{name}<br>
[3] uni/fabric/monfab-{name}/tarfab-{scope}/tacacssrc-{name}<br>
[4] uni/fabric/moncommon/fsevp-{code}/tacacssrc-{name}<br>
[5] uni/tn-{name}/monepg-{name}/tarepg-{scope}/fsevp-{code}/tacacssrc-{name}<br>
[6] uni/infra/moninfra-{name}/tarinfra-{scope}/fsevp-{code}/tacacssrc-{name}<br>
[7] uni/fabric/monfab-{name}/tarfab-{scope}/fsevp-{code}/tacacssrc-{name}<br>
[8] uni/fabric/moncommon/esevp-{code}/tacacssrc-{name}<br>
[9] uni/tn-{name}/monepg-{name}/tarepg-{scope}/esevp-{code}/tacacssrc-{name}<br>
[10] uni/infra/moninfra-{name}/tarinfra-{scope}/esevp-{code}/tacacssrc-{name}<br>
[11] uni/fabric/monfab-{name}/tarfab-{scope}/esevp-{code}/tacacssrc-{name}<br>
[12] uni/fabric/moncommon/tacacssrc-{name}<br>
[13] uni/tn-{name}/monepg-{name}/tacacssrc-{name}<br>
[14] uni/infra/moninfra-{name}/tacacssrc-{name}<br>
[15] uni/fabric/monfab-{name}/tacacssrc-{name}<br>

## Example Usage ##

```hcl
resource "aci_tacacs_source" "example" {
  parent_dn   = "uni/fabric/moncommon"
  name        = "example"
  annotation  = "orchestrator:terraform"
  incl        = ["audit","session"]
  min_sev     = "info"
  name_alias  = "tacacs_source_alias"
  description = "From Terraform"
}
```

## Argument Reference ##

* `name` - (Required) Name of object TACACS Source.
* `parent_dn` - (Optional) Distinguished name of parent object of TACACS Source. Default value is "uni/fabric/moncommon". Type: String.
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
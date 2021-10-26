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
* `Distinguished Named` - <br>
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
data "aci_tacacs_source" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) name of object TACACS Source.

## Attribute Reference ##
* `parent_dn` - (Optional) Distinguished name of parent object of TACACS Source.
* `id` - Attribute id set to the Dn of the TACACS Source.
* `annotation` - (Optional) Annotation of object TACACS Source.
* `name_alias` - (Optional) Name Alias of object TACACS Source.
* `description` - (Optional) Description of object TACACS Source.
* `incl` - (Optional) Include Action. The information to include for the call home source.
* `min_sev` - (Optional) minSev. 

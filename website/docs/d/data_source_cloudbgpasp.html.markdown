---
layout: "aci"
page_title: "ACI: aci_autonomous_system_profile"
sidebar_current: "docs-aci-data-source-autonomous_system_profile"
description: |-
  Data source for ACI Autonomous System Profile
---

# aci_autonomous_system_profile #
Data source for ACI Autonomous System Profile  
<b>Note: This resource is supported in Cloud APIC only.


## Example Usage ##

```hcl
data "aci_autonomous_system_profile" "auto_prof" {
}
```
## Argument Reference ##
This data source don't have any arguments.

## Attribute Reference

* `id` - Attribute id set to the Dn of the Autonomous System Profile.
* `annotation` - (Optional) annotation for object autonomous_system_profile.
* `asn` - (Optional) A number that uniquely identifies an autonomous system. 
* `name_alias` - (Optional) name_alias for object autonomous_system_profile.

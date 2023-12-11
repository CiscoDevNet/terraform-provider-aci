---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_autonomous_system_profile"
sidebar_current: "docs-aci-data-source-aci_autonomous_system_profile"
description: |-
  Data source for ACI Autonomous System Profile
---

# aci_autonomous_system_profile #
Data source for ACI Autonomous System Profile  
<b>Note: This resource is supported in Cloud Network Controller only.


## Example Usage ##

```hcl
data "aci_autonomous_system_profile" "auto_prof" {
}
```
## Argument Reference ##
This data source don't have any arguments.

## Attribute Reference

* `id` - Attribute id set to the Dn of the Autonomous System Profile.
* `annotation` - (Optional) Annotation for object Autonomous System Profile.
* `description` - (Optional) Description for object Autonomous System Profile.
* `asn` - (Optional) A number that uniquely identifies an autonomous system. 
* `name_alias` - (Optional) Name alias for object Autonomous System Profile.

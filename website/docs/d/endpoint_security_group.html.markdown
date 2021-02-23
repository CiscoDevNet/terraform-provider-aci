---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group"
sidebar_current: "docs-aci-data-source-endpoint_security_group"
description: |-
  Data source for ACI Endpoint Security Group
---

# aci_endpoint_security_group #
Data source for ACI Endpoint Security Group

## Example Usage ##

```hcl
data "aci_endpoint_security_group" "example" {

  application_profile_dn  = "${aci_application_profile.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `application_profile_dn` - (Required) Distinguished name of parent ApplicationProfile object.
* `name` - (Required) name of Object endpoint_security_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Endpoint Security Group.
* `annotation` - (Optional) annotation for object endpoint_security_group.
* `exception_tag` - (Optional) exception_tag for object endpoint_security_group.
* `flood_on_encap` - (Optional) flood_on_encap for object endpoint_security_group.
* `match_t` - (Optional) match criteria
* `name_alias` - (Optional) name_alias for object endpoint_security_group.
* `pc_enf_pref` - (Optional) enforcement preference
* `pref_gr_memb` - (Optional) pref_gr_memb for object endpoint_security_group.
* `prio` - (Optional) qos priority class id
* `userdom` - (Optional) userdom for object endpoint_security_group.

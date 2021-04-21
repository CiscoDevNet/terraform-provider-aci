---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group"
sidebar_current: "docs-aci-data-source-endpoint_security_group"
description: |-
  Data source for ACI Endpoint Security Group
---

# aci_endpoint_security_group #
Data source for ACI Endpoint Security Group


## API Information ##
* `Class` - fvESg
* `Distinguished Named` - uni/tn-%s/ap-%s/esg-%s

## GUI Location ##


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

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Endpoint Security Group.
* `annotation` - (Optional) Annotation of object endpoint_security_group.
* `name_alias` - (Optional) Name Alias of object endpoint_security_group.
* `exception_tag` - (Optional) 
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2
                     Multicast/Broadcast and Link Local Layer should
                     be flooded only on ENCAP or based on bridg-domain
                     settings
* `match_t` - (Optional) The provider label match criteria.
* `pc_enf_pref` - (Optional) The preferred policy control.
* `pref_gr_memb` - (Optional) Represents parameter used to determine
                    if EPg is part of a group that does not
                    a contract for communication
* `prio` - (Optional) The QoS priority class identifier.

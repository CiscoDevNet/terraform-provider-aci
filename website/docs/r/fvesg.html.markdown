---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group"
sidebar_current: "docs-aci-resource-endpoint_security_group"
description: |-
  Manages ACI Endpoint Security Group
---

# aci_endpoint_security_group #
Manages ACI Endpoint Security Group

## API Information ##
* `Class` - fvESg
* `Distinguished Named` - uni/tn-%s/ap-%s/esg-%s

## GUI Location ##


## Example Usage ##

```hcl
resource "aci_endpoint_security_group" "example" {
  application_profile_dn  = "${aci_application_profile.example.id}"
  name  = "example"
  exception_tag  = []
  flood_on_encap  = ["defaultValue","disabled","enabled",]
  match_t  = ["All","AtleastOne","AtmostOne","None","defaultValue",]
  pc_enf_pref  = ["defaultValue","enforced","unenforced",]
  pref_gr_memb  = ["defaultValue","exclude","include",]
  prio  = []
}
```

## Argument Reference ##
* `application_profile_dn` - (Required) Distinguished name of parent ApplicationProfile object.
* `name` - (Required) Name of Object endpoint_security_group.
* `annotation` - (Optional) Annotation of object endpoint_security_group.
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

## Importing ##

An existing EndpointSecurityGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import endpoint_security_group.example <Dn>
```
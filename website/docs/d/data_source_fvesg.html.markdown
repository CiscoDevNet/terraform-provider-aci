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
* `Distinguished Named` - uni/tn-{name}/ap-{name}/esg-{name}

## GUI Information ##

* `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups

## Example Usage ##

```hcl
data "aci_endpoint_security_group" "example" {
  application_profile_dn  = aci_application_profile.example.id
  name  = "example"
}
```

## Argument Reference ##

* `application_profile_dn` - (Required) Distinguished name of parent Application Profile object.
* `name` - (Required) name of object Endpoint Security Group.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Endpoint Security Group.
* `annotation` - (Optional) Annotation of object Endpoint Security Group.
* `name_alias` - (Optional) Name Alias of object Endpoint Security Group.
* `flood_on_encap` - (Optional) Handles L2 Multicast/Broadcast and Link-Layer traffic at EPG level. It represents Control at EPG level and decides if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP, or based on bridge-domain settings.
* `match_t` - (Optional) The provider label match criteria.
* `pc_enf_pref` - (Optional) The preferred policy control.
* `pref_gr_memb` - (Optional) Represents parameter used to determine
                    if EPg is part of a group that does not
                    a contract for communication.
* `prio` - (Optional) The QoS priority class identifier.

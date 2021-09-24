---
layout: "aci"
page_title: "ACI: aci_vrf_snmp_context_community"
sidebar_current: "docs-aci-data-source-vrf_snmp_context_community"
description: |-
  Data source for ACI SNMP Community
---

# aci_vrf_snmp_context_community #

Data source for ACI SNMP Community


## API Information ##

* `Class` - snmpCommunityP
* `Distinguished Named` - uni/tn-{name}/ctx-{name}/snmpctx/community-{name}

## GUI Information ##

* `Location` - Tenant -> Networking -> VRFs -> Policy -> SNMP Context



## Example Usage ##

```hcl
data "aci_vrf_snmp_context_community" "example" {
  vrf_snmp_context_dn  = aci_vrf_snmp_context.example.id
  name  = "example"
}
```

## Argument Reference ##

* `vrf_snmp_context_dn` - (Required) Distinguished name of parent VRF SNMP Context object.
* `name` - (Required) name of object SNMP Community.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the SNMP Community.
* `annotation` - (Optional) Annotation of object SNMP Community.
* `name_alias` - (Optional) Name Alias of object SNMP Community.
* `description` - (Optional) Description of object SNMP Community.


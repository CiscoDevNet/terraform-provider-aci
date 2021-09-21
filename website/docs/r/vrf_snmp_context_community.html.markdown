---
layout: "aci"
page_title: "ACI: aci_vrf_snmp_context_community"
sidebar_current: "docs-aci-resource-vrf_snmp_context_community"
description: |-
  Manages ACI SNMP Community
---

# aci_vrf_snmp_context_community #

Manages ACI SNMP Community

## API Information ##

* `Class` - snmpCommunityP
* `Distinguished Named` - uni/tn-{name}/ctx-{name}/snmpctx/community-{name}

## GUI Information ##

* `Location` - Tenant -> Networking -> VRFs -> Policy -> SNMP Context


## Example Usage ##

```hcl
resource "aci_vrf_snmp_context_community" "example" {
	vrf_dn = aci_vrf.test.id
	name = "test"
	description = "From Terraform"
	annotation = "Test_Annotation"
	name_alias = "Test_name_alias"
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of parent VRF object.
* `name` - (Required) Name of object SNMP Community.
* `annotation` - (Optional) Annotation of object SNMP Community.
* `description` - (Optional) Description of object SNMP Community.
* `name_alias` - (Optional) Name alias of object SNMP Community.

## Importing ##

An existing SNMPCommunity can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vrf_snmp_context_community.example <Dn>
```
---
subcategory: ""
layout: "aci"
page_title: "ACI: aci_snmp_community"
sidebar_current: "docs-aci-resource-snmp_community"
description: |-
  Manages ACI SNMP Community
---

# aci_snmp_community #

Manages ACI SNMP Community

## API Information ##

* `Class` - snmpCommunityP
* `Distinguished Name` - {parent_dn}/community-{name}

## GUI Information ##

* `Locations` 
- Fabric > Fabric Policies > Policies > Pod > SNMP > {snmp_policy} > Community Policies
- Tenant > {tenant} > Networking > VRFs > {vrf} > SNMP Context > Community Policies


## Example Usage ##

```hcl
resource "aci_snmp_community" "example" {
  parent_dn  = aci_snmp_policy.example.id
  name  = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent object.
* `name` - (Required) Name of the SNMP Community.
* `annotation` - (Optional) Annotation of the SNMP Community.
* `name_alias` - (Optional) Name Alias of the SNMP Community.

## Importing ##

An existing SNMPCommunity can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_snmp_community.example <Dn>
```
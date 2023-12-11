---
subcategory: ""
layout: "aci"
page_title: "ACI: aci_snmp_community"
sidebar_current: "docs-aci-data-source-aci_snmp_community"
description: |-
  Data source for ACI SNMP Community
---

# aci_snmp_community #

Data source for ACI SNMP Community


## API Information ##

* `Class` - snmpCommunityP
* `Distinguished Name` - {parent_dn}/community-{name}

## GUI Information ##

* `Locations` 
- Fabric > Fabric Policies > Policies > Pod > SNMP > {snmp_policy} > Community Policies
- Tenant > {tenant} > Networking > VRFs > {vrf} > SNMP Context > Community Policies



## Example Usage ##

```hcl
data "aci_snmp_community" "example" {
  parent_dn  = aci_snmp_policy.example.id
  name  = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent object.
* `name` - (Required) Name of object SNMP Community.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the SNMP Community.
* `annotation` - (Optional) Annotation of the SNMP Community.
* `name_alias` - (Optional) Name Alias of the SNMP Community.

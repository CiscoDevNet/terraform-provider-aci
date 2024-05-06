---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_endpoint_security_group"
sidebar_current: "docs-aci-data-source-aci_endpoint_security_group"
description: |-
  Data source for ACI Endpoint Security Group
---

# aci_endpoint_security_group #

Data source for ACI Endpoint Security Group

## API Information ##

* `Class` - fvESg
* `Distinguished Name` - uni/tn-{name}/ap-{name}/esg-{name}

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
* `name` - (Required) Name of the object Endpoint Security Group.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Endpoint Security Group.
* `annotation` - (Read-Only) Annotation of the object Endpoint Security Group.
* `description` - (Read-Only) Description of the object Endpoint Security Group.
* `name_alias` - (Read-Only) Name Alias of the object Endpoint Security Group.
* `match_t` - (Read-Only) The provider label match criteria.
* `pc_enf_pref` - (Read-Only) The preferred policy control.
* `pref_gr_memb` - (Read-Only) Represents parameter used to determine if EPG is part of a group that does not a contract for communication.

---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_epg_selector"
sidebar_current: "docs-aci-data-source-aci_endpoint_security_group_epg_selector"
description: |-
  Data source for ACI Endpoint Security Group EPG Selector
---

# aci_endpoint_security_group_epg_selector #

Data source for ACI Endpoint Security Group EPG Selector


## API Information ##

* `Class` - fvEPgSelector
* `Distinguished Name` - uni/tn-{name}/ap-{name}/esg-{name}/epgselector-[{matchEpgDn}]

## GUI Information ##

* `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups > Selectors > EPG Selectors


## Example Usage ##

```hcl
data "aci_endpoint_security_group_epg_selector" "example" {
  endpoint_security_group_dn  = aci_endpoint_security_group.example.id
  match_epg_dn  = aci_application_epg.example.id 
}
```

## Argument Reference ##

* `endpoint_security_group_dn` - (Required) Distinguished name of parent Endpoint Security Group object.
* `match_epg_dn` - (Required) EPG Dn to be associated. 

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Endpoint Security Group EPG Selector.
* `annotation` - (Optional) Annotation of object Endpoint Security Group EPG Selector.
* `name_alias` - (Optional) Name Alias of object Endpoint Security Group EPG Selector.

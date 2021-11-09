---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_epg_selector"
sidebar_current: "docs-aci-resource-endpoint_security_group_epg_selector"
description: |-
  Manages ACI Endpoint Security Group EPG Selector
---

# aci_endpoint_security_group_epg_selector #

Manages ACI Endpoint Security Group EPG Selector

## API Information ##

* `Class` - fvEPgSelector
* `Distinguished Named` - uni/tn-{name}/ap-{name}/esg-{name}/epgselector-[{match_epg_dn}]

## GUI Information ##

* `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups > Selectors > EPG Selectors


## Example Usage ##

```hcl
resource "aci_endpoint_security_group_epg_selector" "example" {
  endpoint_security_group_dn  = aci_endpoint_security_group.example.id
  match_epg_dn  = aci_application_epg.example.id 
  annotation = "orchestrator:terraform"
}
```

## Argument Reference ##

* `endpoint_security_group_dn` - (Required) Distinguished name of parent Endpoint Security Group object.
* `match_epg_dn` - (Required) EPG Dn to be associated.
* `annotation` - (Optional) Annotation of object Endpoint Security Group EPG Selector.

## Importing ##

An existing Endpoint Security Group EPG Selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_endpoint_security_group_epg_selector.example <Dn>
```
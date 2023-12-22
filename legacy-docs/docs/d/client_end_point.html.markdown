---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_client_end_point"
sidebar_current: "docs-aci-data-source-aci_client_end_point"
description: |-
  Data source for ACI Client End Point
---

# aci_client_end_point

Data source for ACI Client End Point

## API Information ##

* `Class` - fvCEp
* `Distinguished Name` - uni/tn-{tenant_name}/ctx-{ctx_name}/cep-{client_end_point_name}
* `Distinguished Name` - uni/tn-{tenant_name}/ap-{ap_name}/epg-{epg_name}/cep-{client_end_point_name}
* `Distinguished Name` - uni/tn-{tenant_name}/l2out-{l2out_name}/instP-{instance_profile_name}/cep-{client_end_point_name}

## GUI Information ##

* `Location` - Tenant -> Application Profiles -> Application EPGs -> Operational -> Client End-Points
* `Location` - Tenant -> Networking -> VRFs -> Operational -> Client End-Points
## Example Usage

```hcl
data "aci_client_end_point" "check" {
  mac                 = "25:56:68:78:98:74"
  ip                  = "1.2.3.4"
  vlan                = "5"
  allow_empty_result  = true
}
```

## Argument Reference

- `name` - (Optional) Name of the Client End Point object.
- `mac` - (Optional) MAC address of the Client End Point object.
- `ip` - (Optional) IP address of the Client End Point object.
- `vlan` - (Optional) VLAN of the Client End Point object.
- `allow_empty_result` - (Optional) Empty return instead of error when client is not found. Default value is "false". 

## Attribute Reference

- `id` - Attribute id set as all Dns for matching with the Client End Point.
- `fvcep_objects` - List of all Client End Point objects which matched to the given filter attributes.
- `fvcep_objects.name` - Name of the Client End Point object.
- `fvcep_objects.mac` - Mac address of the Client End Point object.
- `fvcep_objects.ip` - IP address of the Client End Point object.
- `fvcep_objects.ips` - List of `fvIp` addresses mapped to the Client End Point object.
- `fvcep_objects.vlan` - VLAN of the Client End Point object.
- `fvcep_objects.tenant_name` - Parent Tenant name of the Client End Point object.
- `fvcep_objects.vrf_name` - Parent VRF name of the Client End Point object.
- `fvcep_objects.application_profile_name` - Parent Application Profile name of the Client End Point object.
- `fvcep_objects.epg_name` - Parent EPG name of the Client End Point object.
- `fvcep_objects.esg_name` - Parent ESG name of the Client End Point object.
- `fvcep_objects.l2out_name` - Parent L2Out name of the Client End Point object.
- `fvcep_objects.instance_profile_name` - Parent Instance Profile name of the Client End Point object.
- `fvcep_objects.endpoint_path` - List of endpoint paths associated with the Client End Point object.
- `fvcep_objects.base_epg.tenant_name` - Name of the Tenant of the base EPG of the Client End Point object when the Client End Point is associated with an ESG.
- `fvcep_objects.base_epg.application_profile_name` - Name of the Application Profile of the base EPG of the Client End Point object when the Client End Point is associated with an ESG.
- `fvcep_objects.base_epg.epg_name` - Name of the base EPG of the Client End Point object when the Client End Point is associated with an ESG.

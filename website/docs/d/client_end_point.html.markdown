---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_client_end_point"
sidebar_current: "docs-aci-data-source-client_end_point"
description: |-
  Data source for ACI Client End Point
---

# aci_client_end_point

Data source for ACI Client End Point

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

- `name` - (Optional) Name of Object client end point.
- `mac` - (Optional) MAC address of the object client end point.
- `ip` - (Optional) IP address of the object client end point.
- `vlan` - (Optional) VLAN for the object client end point.
- `allow_empty_result` - (Optional) Empty return instead of error when client is not found. Default value is "false". 

## Attribute Reference

- `id` - Attribute id set as all Dns for matching the Client End Point.
- `fvcep_objects` - List of all client end point objects which matched to the given filter attributes.
- `fvcep_objects.name` - Name of object client end point.
- `fvcep_objects.mac` - Mac address of object client end point.
- `fvcep_objects.ip` - IP address of object client end point.
- `fvcep_objects.vlan` - vlan of client end point object.
- `fvcep_objects.tenant_name` - Parent Tenant name for matched client end point.
- `fvcep_objects.vrf_name` - Parent vrf name for matched client end point.
- `fvcep_objects.application_profile_name` - Parent application profile name for matched client end point.
- `fvcep_objects.epg_name` - Parent epg name for matched client end point.
- `fvcep_objects.l2out_name` - Parent l2out name for matched client end point.
- `fvcep_objects.instance_profile_name` - Parent instance profile name for matched client end point.
- `fvcep_objects.endpoint_path` - List of endpoint paths associated with client end point.

---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_destination_of_redirected_traffic"
sidebar_current: "docs-aci-data-source-aci_destination_of_redirected_traffic"
description: |-
  Data source for ACI Destination of redirected traffic
---

# aci_destination_of_redirected_traffic

Data source for ACI Destination of redirected traffic

## API Information ##

* `Class` - vnsRedirectDest
* `Distinguished Name` - uni/tn-{tenant_name}/svcCont/svcRedirectPol-{service_redirect_policy_name}/RedirectDest_ip-[{ip}]

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> L4-L7 Policy-Based Redirect -> L3 Destinations

## Example Usage

```hcl

data "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn  = aci_service_redirect_policy.example.id
  ip                          = "1.2.3.4"
}

```

## Argument Reference

- `service_redirect_policy_dn` - (Required) Distinguished name of the parent Service Redirect Policy object.
- `ip` - (Required) The IP address.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Destination of redirected traffic.
- `annotation` - (Optional) Annotation of the destination of redirected traffic object.
- `destination` - (Optional) Specifies the description of a policy component.
- `dest_name` - (Optional) The name of the destination of redirected traffic object. 
- `ip2` - (Optional) IP2 of the destination of redirected traffic object.
- `mac` - (Optional) The MAC address.
- `name_alias` - (Optional) Name alias of the destination of redirected traffic object.
- `pod_id` - (Optional) The Pod identifier.

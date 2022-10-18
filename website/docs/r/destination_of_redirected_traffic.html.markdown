---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_destination_of_redirected_traffic"
sidebar_current: "docs-aci-resource-destination_of_redirected_traffic"
description: |-
  Manages ACI Destination of redirected traffic
---

# aci_destination_of_redirected_traffic

Manages ACI Destination of redirected traffic

## API Information ##

* `Class` - vnsRedirectDest
* `Distinguished Name` - uni/tn-{tenant_name}/svcCont/svcRedirectPol-{service_redirect_policy_name}/RedirectDest_ip-[{ip}]

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> L4-L7 Policy-Based Redirect -> L3 Destinations

## Example Usage

```hcl

resource "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn  = aci_service_redirect_policy.example.id
  ip                          = "1.2.3.4"
  ip2                         = "10.20.30.40"
  dest_name                   = "last"
  pod_id                      = "5"
  annotation                  = "load_traffic_dest"
  description                 = "From Terraform"
  name_alias                  = "load_traffic_dest"
}

```

## Argument Reference

- `service_redirect_policy_dn` - (Required) Distinguished name of the parent Service Redirect Policy object.
- `ip` - (Required) The IP address.
- `mac` - (Optional) The MAC address. This is a required value for APIC prior to Version 5.2 release. This value can be foregone by enabling IPSLA on APIC Version 5.2 and above due to dynamic mac detection feature.
- `annotation` - (Optional) Annotation for the object destination of redirected traffic.
- `description` - (Optional) Description for the object destination of redirected traffic.
- `dest_name` - (Optional) The destination name to which the data was exported. 
- `ip2` - (Optional) IP2 for the object destination of redirected traffic. Default value: "0.0.0.0"
- `name_alias` - (Optional) Name alias for the object destination of redirected traffic.
- `pod_id` - (Optional) The POD identifier. Allowed value range: "1" to "255". Default value: "1"

- `relation_vns_rs_redirect_health_group` - (Optional) Relation to class vns Redirect Health Group. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Destination of redirected traffic.

## Importing

An existing Destination of redirected traffic can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_destination_of_redirected_traffic.example <Dn>
```

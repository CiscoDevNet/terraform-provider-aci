---
layout: "aci"
page_title: "ACI: aci_destination_of_redirected_traffic"
sidebar_current: "docs-aci-resource-destination_of_redirected_traffic"
description: |-
  Manages ACI Destination of redirected traffic
---

# aci_destination_of_redirected_traffic #
Manages ACI Destination of redirected traffic

## Example Usage ##

```hcl

resource "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn  = "${aci_service_redirect_policy.example.id}"
  ip                          = "1.2.3.4"
  mac                         = "12:25:56:98:45:74"
  ip2                         = "10.20.30.40"
  dest_name                   = "last"
  pod_id                      = "5"
}

```


## Argument Reference ##
* `service_redirect_policy_dn` - (Required) Distinguished name of parent Service Redirect Policy object.
* `ip` - (Required) ip of Object destination of redirected traffic.
* `mac` - (Required) mac address.
* `annotation` - (Optional) annotation for object destination of redirected traffic.
* `dest_name` - (Optional) destination name to which data was exported
* `ip` - (Optional) ip address.
* `ip2` - (Optional) ip2 for object destination of redirected traffic.
* `name_alias` - (Optional) name_alias for object destination of redirected traffic.
* `pod_id` - (Optional) pod id.

* `relation_vns_rs_redirect_health_group` - (Optional) Relation to class vns Redirect Health Group. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Destination of redirected traffic.

## Importing ##

An existing Destination of redirected traffic can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_destinationofredirectedtraffic.example <Dn>
```
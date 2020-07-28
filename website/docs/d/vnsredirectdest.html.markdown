---
layout: "aci"
page_title: "ACI: aci_destination_of_redirected_traffic"
sidebar_current: "docs-aci-data-source-destination_of_redirected_traffic"
description: |-
  Data source for ACI Destination of redirected traffic
---

# aci_destination_of_redirected_traffic #
Data source for ACI Destination of redirected traffic

## Example Usage ##

```hcl

data "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn  = "${aci_service_redirect_policy.example.id}"
  ip                          = "1.2.3.4"
}

```


## Argument Reference ##
* `service_redirect_policy_dn` - (Required) Distinguished name of parent Service Redirect Policy object.
* `ip` - (Required) ip of Object destination of redirected traffic.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Destination of redirected traffic.
* `annotation` - (Optional) annotation for object destination of redirected traffic.
* `dest_name` - (Optional) destination name to which data was exported.
* `ip` - (Optional) ip address.
* `ip2` - (Optional) ip2 for object destination of redirected traffic.
* `mac` - (Optional) mac address.
* `name_alias` - (Optional) name_alias for object destination of redirected traffic.
* `pod_id` - (Optional) pod id.

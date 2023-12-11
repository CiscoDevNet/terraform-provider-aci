---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_radius_provider"
sidebar_current: "docs-aci-data-source-aci_radius_provider"
description: |-
  Data source for ACI RADIUS Provider
---

# aci_radius_provider #
Data source for ACI RADIUS Provider


## API Information ##
* `Class` - aaaRadiusProvider
* `Distinguished Name` - uni/userext/duoext/radiusprovider-{name} & uni/userext/radiusext/radiusprovider-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> RADIUS & Admin -> AAA -> Authentication -> DUO -> Radius

## Example Usage ##

```hcl
data "aci_radius_provider" "example" {
  name = "example"
  type = "radius"
}
```

## Argument Reference ##
* `name` - (Required) Host name or IP address of object RADIUS Provider.
* `type` - (Required) Type of object RADIUS Provider. Allowed values are "duo" and "radius".

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the RADIUS Provider.
* `annotation` - (Optional) Annotation of object RADIUS Provider.
* `name_alias` - (Optional) Name Alias of object RADIUS Provider.
* `description` - (Optional) Description of object RADIUS Provider.
* `auth_port` - (Optional) The service port number for RADIUS service.
* `auth_protocol` - (Optional) The RADIUS authentication protocol.
* `monitor_server` - (Optional) Periodic Server Monitoring. 
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. 
* `retries` - (Optional) Number of retries for a for communication with a RADIUS provider server.
* `timeout` - (Optional) The timeout for communication with a RADIUS provider server.

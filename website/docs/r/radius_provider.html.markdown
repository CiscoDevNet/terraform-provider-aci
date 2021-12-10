---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_radius_provider"
sidebar_current: "docs-aci-resource-radius_provider"
description: |-
  Manages ACI RADIUS Provider
---

# aci_radius_provider #
Manages ACI RADIUS Provider

## API Information ##
* `Class` - aaaRadiusProvider
* `Distinguished Named` - uni/userext/duoext/radiusprovider-{name} & uni/userext/radiusext/radiusprovider-{name}

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> RADIUS & Admin -> AAA -> Authentication -> DUO -> Radius

## Example Usage ##
```hcl
resource "aci_radius_provider" "example" {
  name  = "example"
  type = "radius"
  annotation = "orchestrator:terraform"
  auth_port = "1812"
  auth_protocol = "pap"
  key = "example_key_value"
  monitor_server = "disabled"
  monitoring_password = "example_monitoring_password"
  monitoring_user = "default"
  retries = "1"
  timeout = "5"
  description = "from terraform"
  name_alias = "example_name_alias_value"
}
```

## Argument Reference ##
* `name` - (Required) Host name or IP address of object RADIUS Provider.
* `type` - (Required) Type of object RADIUS Provider. Allowed values are "duo" and "radius".
* `annotation` - (Optional) Annotation of object RADIUS Provider.
* `name_alias` - (Optional) Name Alias of object RADIUS Provider.
* `description` - (Optional) Description of object RADIUS Provider.
* `auth_port` - (Optional) The service port number for the RADIUS service. Allowed range: "1" - "65535". Default value is "1812".
* `auth_protocol` - (Optional) Authentication Protocol.The RADIUS authentication protocol. Allowed values are "chap", "mschap" and "pap". Default value is "pap". Type: String.
* `key` - (Optional) A password for the AAA provider database.
* `monitor_server` - (Optional) Periodic Server Monitoring. Allowed values are "disabled" and "enabled". Default value is "disabled". Type: String.
* `monitoring_password` - (Optional) Periodic Server Monitoring Password.
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. Default value is "default".
* `retries` - (Optional) Number of retries for a for communication with a RADIUS provider server. Allowed range is "1" - "5". Default value is "1".
* `timeout` - (Optional) The timeout for communication with a RADIUS provider server. Allowed range: "1" - "60". Default value is "5". (NOTE: For "duo" RADIUS providers, the value of timeout should be greater than or equal to "30".)
* `relation_aaa_rs_prov_to_epp` - (Optional) Represents the relation to a Relation to AProvider Reachability EPP (class fvAREpP).  Type: String.
* `relation_aaa_rs_sec_prov_to_epg` - (Optional) Represents the relation to a Attachable Target Group (class fvATg). A source relation to the endpoint group through which the provider server is reachable. Type: String.



## Importing ##

An existing RADIUSProvider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_radius_provider.example <Dn>
```
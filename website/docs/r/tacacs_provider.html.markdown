---
layout: "aci"
page_title: "ACI: aci_tacacs_provider"
sidebar_current: "docs-aci-resource-tacacs_provider"
description: |-
  Manages ACI TACACS Provider
---

# aci_tacacs_provider #

Manages ACI TACACS Provider

## API Information ##

* `Class` - aaaTacacsPlusProvider
* `Distinguished Named` - uni/userext/tacacsext/tacacsplusprovider-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> TACACS 


## Example Usage ##

```hcl
resource "aci_tacacs_provider" "example" {
  name                = "example"
  annotation          = "orchestrator:terraform"
  name_alias          = "tacacs_provider_alias"
  description         = "From Terraform"
  auth_protocol       = "pap"
  key                 = "example"
  monitor_server      = "disabled"
  monitoring_password = "example" 
  monitoring_user     = "default"
  port                = "49"
  retries             = "1"
  timeout             = "5"
}
```

## Argument Reference ##


* `name` - (Required) Name of object TACACS Provider.
* `annotation` - (Optional) Annotation of object TACACS Provider.
* `auth_protocol` - (Optional) TACACS Authentication Protocol. The TACACS authentication protocol. Allowed values are "chap", "mschap", "pap". Default value is "pap". Type: String.
* `key` - (Optional) Key. A password for the AAA provider database. Type: String.
* `monitor_server` - (Optional) Periodic Server Monitoring. Allowed values are "disabled", "enabled". Default value is "disabled". Type: String.
* `monitoring_password` - (Optional) Periodic Server Monitoring Password. Type: String.
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. Default value is "default". Type: String.
* `port` - (Optional) Port. The service port number for the TACACS service. Allowed range is "1"-"65535". Default value is "49". Type: String.
* `retries` - (Optional) Retries. Allowed range is "1"-"5" and default value is "1". Type: String.
* `timeout` - (Optional) Timeout in Seconds. The timeout for communication with the TACACS provider server. Allowed range is "1"-"60". Default value is "5". Type: String.
* `relation_aaa_rs_prov_to_epp` - (Optional) Represents the relation to a Relation to AProvider Reachability EPP (class fvAREpP).  Type: String.
* `relation_aaa_rs_sec_prov_to_epg` - (Optional) Represents the relation to a Attachable Target Group (class fvATg). A source relation to the endpoint group through which the provider server is reachable. Type: String.



## Importing ##

An existing TACACSProvider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tacacs_provider.example <Dn>
```
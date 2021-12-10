---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_rsa_provider"
sidebar_current: "docs-aci-resource-rsa_provider"
description: |-
  Manages ACI RSA Provider
---

# aci_rsa_provider #

Manages ACI RSA Provider

## API Information ##

* `Class` - aaaRsaProvider
* `Distinguished Named` - uni/userext/rsaext/rsaprovider-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> RSA  


## Example Usage ##

```hcl
resource "aci_rsa_provider" "example" {
  name                   = "example"
  name_alias             = "rsa_provider_alias"
  description            = "From Terraform"
  annotation             = "orchestrator:terraform"
  auth_port              = "1812"
  auth_protocol          = "pap"
  key                    = "key_example"
  monitor_server         = "disabled"
  monitoring_password    = "monitoring_password_example"
  monitoring_user        = "default"
  retries                = "1"
  timeout                = "5"
}
```

## Argument Reference ##


* `name` - (Required) Name of object RSA Provider.
* `annotation` - (Optional) Annotation of object RSA Provider.
* `name_alias` - (Optional) Name Alias of object RSA Provider.
* `description` - (Optional) Description of object RSA Provider.
* `auth_port` - (Optional) Port. Allowed range is "1"-"65535". Default value is "1812". Type: String.
* `auth_protocol` - (Optional) Authentication Protocol. Allowed values are "chap", "mschap", "pap". Default value is "pap". Type: String.
* `key` - (Optional) Key. A password for the AAA provider database.
* `monitor_server` - (Optional) Periodic Server Monitoring. Allowed values are "disabled", "enabled". Default value is "disabled". Type: String.
* `monitoring_password` - (Optional) Periodic Server Monitoring Password. Type: String.
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. Default value is "default". Type: String.
* `retries` - (Optional) Retries. Allowed range is "1"-"5". Default value is "1". Type: String.
* `timeout` - (Optional) Timeout in Seconds. The amount of time between authentication attempts. Allowed range is "1"-"60" and default value is "5". Type: String.
* `relation_aaa_rs_prov_to_epp` - (Optional) Represents the relation to a Relation to AProvider Reachability EPP (class fvAREpP). Type: String.
* `relation_aaa_rs_sec_prov_to_epg` - (Optional) Represents the relation to a Attachable Target Group (class fvATg). A source relation to the endpoint group through which the provider server is reachable. Type: String.



## Importing ##

An existing RSAProvider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_rsa_provider.example <Dn>
```
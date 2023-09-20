---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_rsa_provider"
sidebar_current: "docs-aci-data-source-rsa_provider"
description: |-
  Data source for ACI RSA Provider
---

# aci_rsa_provider #

Data source for ACI RSA Provider


## API Information ##

* `Class` - aaaRsaProvider
* `Distinguished Name` - uni/userext/rsaext/rsaprovider-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> RSA 



## Example Usage ##

```hcl
data "aci_rsa_provider" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object RSA Provider.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the RSA Provider.
* `annotation` - (Optional) Annotation of object RSA Provider.
* `name_alias` - (Optional) Name Alias of object RSA Provider.
* `description` - (Optional) Description of object RSA Provider.
* `auth_port` - (Optional) Port. 
* `auth_protocol` - (Optional) Authentication Protocol. 
* `key` - (Optional) Key. A password for the AAA provider database.
* `monitor_server` - (Optional) Periodic Server Monitoring. 
* `monitoring_password` - (Optional) Periodic Server Monitoring Password. 
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. 
* `retries` - (Optional) Retries. null
* `timeout` - (Optional) Timeout in Seconds. The amount of time between authentication attempts.

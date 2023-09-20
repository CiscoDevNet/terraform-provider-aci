---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_tacacs_provider"
sidebar_current: "docs-aci-data-source-tacacs_provider"
description: |-
  Data source for ACI TACACS Provider
---

# aci_tacacs_provider #

Data source for ACI TACACS Provider


## API Information ##

* `Class` - aaaTacacsPlusProvider
* `Distinguished Name` - uni/userext/tacacsext/tacacsplusprovider-{name}

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> TACACS



## Example Usage ##

```hcl
data "aci_tacacs_provider" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object TACACS Provider.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the TACACS Provider.
* `annotation` - (Optional) Annotation of object TACACS Provider.
* `name_alias` - (Optional) Name Alias of object TACACS Provider.
* `description` - (Optional) Description of object TACACS Provider.
* `auth_protocol` - (Optional) TACACS Authentication Protocol. The TACACS authentication protocol.
* `monitor_server` - (Optional) Periodic Server Monitoring.  
* `monitoring_user` - (Optional) Periodic Server Monitoring Username. 
* `port` - (Optional) Port. The service port number for the TACACS service.
* `retries` - (Optional) Retries. Null.
* `timeout` - (Optional) Timeout in Seconds. The timeout for communication with the TACACS+ provider server.

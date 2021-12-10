---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_console_authentication"
sidebar_current: "docs-aci-data-source-console_authentication"
description: |-
  Data source for ACI Console Authentication
---

# aci_console_authentication #

Data source for ACI Console Authentication 


## API Information ##

* `Class` - aaaConsoleAuth
* `Distinguished Named` - uni/userext/authrealm/consoleauth

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> Console Authentication 



## Example Usage ##

```hcl
data "aci_console_authentication" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Console Authentication.
* `annotation` - (Optional) Annotation of object Console Authentication.
* `name_alias` - (Optional) Name Alias of object Console Authentication.
* `description` - (Optional) Description of object Console Authentication.
* `provider_group` - (Optional) Provider Group. An AAA configuration provider group is a group of remote servers supporting the same AAA protocol that will be used for authentication and authorization. When a provider group is specified, only the servers within that group will be used for authentication and authorization. If no provider group is specified, all servers supporting the realm of AAA protocols will be used for authentication and authorization.
* `realm` - (Optional) Realm. The security method for processing authentication and authorization requests. The realm allows the protected resources on the associated server to be partitioned into a set of protection spaces, each with its own authentication authorization database. This is an abstract class and cannot be instantiated.
* `realm_sub_type` - (Optional) Realm subtype that can be default or Duo.

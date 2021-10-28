---
layout: "aci"
page_title: "ACI: aci_default_authentication"
sidebar_current: "docs-aci-data-source-default_authentication"
description: |-
  Data source for ACI Default Authentication Method for all Logins
---

# aci_default_authentication #
Data source for ACI Default Authentication Method for all Logins


## API Information ##

* `Class` - aaaDefaultAuth
* `Distinguished Named` - uni/userext/authrealm/defaultauth

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> AAA -> Default Authentication

## Example Usage ##

```hcl
data "aci_default_authentication" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Default Authentication Method for all Logins.
* `annotation` - (Optional) Annotation of object Default Authentication Method for all Logins.
* `description` - (Optional) Description of object Default Authentication Method for all Logins.
* `name_alias` - (Optional) Name alias of object Default Authentication Method for all Logins.
* `fallback_check` - (Optional) The parameter to disable fallback in case there are active servers in the default auth type. 
* `provider_group` - (Optional) The group of remote servers supporting the same AAA protocol that will be used for authentication and authorization.
* `realm` - (Optional) The security method for processing authentication and authorization requests. The realm allows the protected resources on the associated server to be partitioned into a set of protection spaces, each with its own authentication authorization database. 
* `realm_sub_type` - (Optional) Realm subtype.
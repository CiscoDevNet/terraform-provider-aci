---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_default_authentication"
sidebar_current: "docs-aci-resource-default_authentication"
description: |-
  Manages ACI Default Authentication Method for all Logins
---

# aci_default_authentication #
Manages ACI Default Authentication Method for all Logins

## API Information ##
* `Class` - aaaDefaultAuth
* `Distinguished Named` - uni/userext/authrealm/defaultauth

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> AAA -> Default Authentication

## Example Usage ##
```hcl
resource "aci_default_authentication" "example" {
  annotation = "orchestrator:terraform"
  fallback_check = "false"
  realm = "local"
  realm_sub_type = "default"
  name_alias = "example_name_alias"
  description = "from terraform"
}
```

## NOTE ##
Users can use the resource of type `aci_default_authentication` to change the configuration of the object Default Authentication Method for all Logins. Users cannot create more than one instance of object Default Authentication Method for all Logins.

## Argument Reference ##
* `annotation` - (Optional) Annotation of object Default Authentication Method for all Logins.
* `description` - (Optional) Description of object Default Authentication Method for all Logins.
* `name_alias` - (Optional) Name alias of object Default Authentication Method for all Logins.
* `fallback_check` - (Optional) The parameter to disable fallback in case there are active servers in the default auth type. Allowed values are "false" and "true". Type: String.
* `provider_group` - (Optional) The group of remote servers supporting the same AAA protocol that will be used for authentication and authorization.
* `realm` - (Optional) The security method for processing authentication and authorization requests. The realm allows the protected resources on the associated server to be partitioned into a set of protection spaces, each with its own authentication authorization database. Allowed values are "ldap", "local", "radius", "rsa", "saml" and "tacacs". Type: String.
* `realm_sub_type` - (Optional) Realm subtype. Allowed values are "default" and "duo". Type: String.


## Importing ##
An existing Default Authentication Method for all Logins can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_default_authentication.example <Dn>
```
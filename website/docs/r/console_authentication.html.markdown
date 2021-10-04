---
layout: "aci"
page_title: "ACI: aci_console_authentication"
sidebar_current: "docs-aci-resource-console_authentication"
description: |-
  Manages ACI Console Authentication
---

# aci_console_authentication #

Manages ACI Console Authentication

## API Information ##

* `Class` - aaaConsoleAuth
* `Distinguished Named` - uni/userext/authrealm/consoleauth

## GUI Information ##

* `Location` - Admin -> AAA -> Authentication -> Console Authentication


## Example Usage ##

```hcl
resource "aci_console_authentication" "example" {
  annotation     = "orchestrator:terraform"
  provider_group = "60"
  realm          = "ldap"
  realm_sub_type = "default"
  name_alias     = "console_alias"
  description    = "From Terraform"
}
```

## NOTE ##
User can use resource of type aci_console_authentication to change configuration of object Console Authentication. User cannot create more than one instances of object Console Authentication.

## Argument Reference ##

* `annotation` - (Optional) Annotation of object Console Authentication Method.
* `name_alias` - (Optional) Name Alias of object Console Authentication. Type: String.
* `description` - (Optional) Description of object Console Authentication. Type: String.
* `provider_group` - (Optional) Provider Group.An AAA configuration provider group is a group of remote servers supporting the same AAA protocol that will be used for authentication and authorization. When a provider group is specified, only the servers within that group will be used for authentication and authorization. If no provider group is specified, all servers supporting the realm of AAA protocols will be used for authentication and authorization.
* `realm` - (Optional) Realm.The security method for processing authentication and authorization requests. The realm allows the protected resources on the associated server to be partitioned into a set of protection spaces, each with its own authentication authorization database. This is an abstract class and cannot be instantiated. Allowed values are "ldap", "local", "radius", "rsa", "saml", "tacacs". Type: String.
* `realm_sub_type` - (Optional) Realm subtype that can be default or Duo.Realm subtype that can be default or Duo Allowed values are "default", "duo". Type: String.


## Importing ##

An existing ConsoleAuthentication can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_console_authentication.example <Dn>
```
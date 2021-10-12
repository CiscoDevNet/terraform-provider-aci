---
layout: "aci"
page_title: "ACI: aci_authentication_properties"
sidebar_current: "docs-aci-resource-authentication_properties"
description: |-
  Manages ACI AAA Authentication Properties and Default Radius Authentication Settings
---

# aci_authentication_properties #
Manages ACI AAA Authentication Properties and Default Radius Authentication Settings

## API Information ##
* `Class` - aaaAuthRealm & aaaPingEp
* `Distinguished Named` - uni/userext/authrealm & uni/userext/pingext

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> AAA -> Policy


## Example Usage ##
```hcl
resource "aci_authentication_properties" "example" {
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias"
  description = "from terraform"
  def_role_policy = "no-login"
  ping_check = "true"
  retries = "1"
  timeout = "5"
}
```

## NOTE ##
* Users can use the resource of type `aci_authentication_properties` to change the configuration of the object AAA Authentication Properties and Default Radius Authentication Settings. Users cannot create more than one instance of object AAA Authentication Properties and Default Radius Authentication Settings.
* Parameters `ping_check`, `retries` and `timeout` are specific to aaaPingEp class. 

## Argument Reference ##
* `annotation` - (Optional) Annotation of objects AAA Authentication Properties and Default Radius Authentication Settings.
* `name_alias` - (Optional) Name Alias of objects AAA Authentication Properties and Default Radius Authentication Settings.
* `description` - (Optional) Description of objects AAA Authentication Properties and Default Radius Authentication Settings.
* `def_role_policy` - (Optional) The default role policy of remote user. Allowed values are "assign-default-role" and "no-login".
* `ping_check` - (Optional) Heart bit ping checks for RADIUS/TACACS/LDAP/SAML/RSA server reachability. Allowed values are "false" and "true".
* `retries` - (Optional) The number of attempts that the authentication method is tried. Allowed range: "0" - "5".
* `timeout` - (Optional) The amount of time between authentication attempts. Allowed range: "1" - "60".


## Importing ##

An existing AAA Authentication Properties and Default Radius Authentication Settings can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_authentication_properties.example <Dn>
```
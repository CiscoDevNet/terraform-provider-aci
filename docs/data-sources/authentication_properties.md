---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_authentication_properties"
sidebar_current: "docs-aci-data-source-aci_authentication_properties"
description: |-
  Data source for ACI AAA Authentication Properties and Default Radius Authentication Settings
---

# aci_authentication_properties #
Data source for ACI AAA Authentication Properties and Default Radius Authentication Settings


## API Information ##
* `Class` - aaaAuthRealm && aaaPingEp
* `Distinguished Name` - uni/userext/authrealm && uni/userext/pingext

## GUI Information ##
* `Location` - Admin -> AAA -> Authentication -> AAA -> Policy

## Example Usage ##
```hcl
data "aci_authentication_properties" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the AAA Authentication.
* `annotation` - (Optional) Annotation of objects AAA Authentication and Default Radius Authentication Settings.
* `name_alias` - (Optional) Name Alias of objects AAA Authentication and Default Radius Authentication Settings.
* `description` - (Optional) Description of objects AAA Authentication and Default Radius Authentication Settings.
* `def_role_policy` - (Optional) The default role policy of object AAA Authentication.
* `ping_check` - (Optional) Heart bit ping checks for RADIUS/TACACS/LDAP/SAML/RSA server reachability.
* `retries` - (Optional) The number of attempts that the authentication method is tried.
* `timeout` - (Optional) The amount of time between authentication attempts.

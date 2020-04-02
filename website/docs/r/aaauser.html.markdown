---
layout: "aci"
page_title: "ACI: aci_local_user"
sidebar_current: "docs-aci-resource-local_user"
description: |-
  Manages ACI Local User
---

# aci_local_user #
Manages ACI Local User

## Example Usage ##

```hcl
resource "aci_local_user" "example" {


  name  = "example"
  account_status  = "example"
  annotation  = "example"
  cert_attribute  = "example"
  clear_pwd_history  = "example"
  email  = "example"
  expiration  = "example"
  expires  = "example"
  first_name  = "example"
  last_name  = "example"
  name_alias  = "example"
  otpenable  = "example"
  otpkey  = "example"
  phone  = "example"
  pwd  = "example"
  pwd_life_time  = "example"
  pwd_update_required  = "example"
  rbac_string  = "example"
  unix_user_id  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object local_user.
* `account_status` - (Optional) local AAA user account status
* `annotation` - (Optional) annotation for object local_user.
* `cert_attribute` - (Optional) cert_attribute for object local_user.
* `clear_pwd_history` - (Optional) clear password history of local user
* `email` - (Optional) email address of the local user
* `expiration` - (Optional) local user account expiration date
* `expires` - (Optional) enables local user account expiration
* `first_name` - (Optional) first name of the local user
* `last_name` - (Optional) last name of the local user
* `name_alias` - (Optional) name_alias for object local_user.
* `otpenable` - (Optional) otpenable for object local_user.
* `otpkey` - (Optional) otpkey for object local_user.
* `phone` - (Optional) phone number of the local user
* `pwd` - (Optional) system user password
* `pwd_life_time` - (Optional) lifetime of the local user password
* `pwd_update_required` - (Optional) pwd_update_required for object local_user.
* `rbac_string` - (Optional) rbac_string for object local_user.
* `unix_user_id` - (Optional) UNIX identifier of the local user



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Local User.

## Importing ##

An existing Local User can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_local_user.example <Dn>
```
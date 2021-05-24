---
layout: "aci"
page_title: "ACI: aci_local_user"
sidebar_current: "docs-aci-resource-local_user"
description: |-
  Manages ACI Local User
---

# aci_local_user

Manages ACI Local User

## Example Usage

```hcl

resource "aci_local_user" "example" {
    name                = "local_user_one"
    account_status      = "active"
    annotation          = "local_user_tag"
    cert_attribute      = "example"
    clear_pwd_history   = "no"
    description         = "from terraform"
    email               = "example@email.com"
    expiration          = "2030-01-01 00:00:00"
    expires             = "yes"
    first_name          = "fname"
    last_name           = "lname"
    name_alias          = "alias_name"
    otpenable           = "no"
    otpkey              = ""
    phone               = "1234567890"
    pwd                 = "StrongPass@123"
    pwd_life_time       = "20"
    pwd_update_required = "no"
    rbac_string         = "example"
}
```

## Argument Reference

- `name` - (Required) Name of Object locally authenticated user.
- `account_status` - (Optional) The status of the locally-authenticated user account.
  Allowed values: "active", "inactive". Default value: "active".
- `annotation` - (Optional) Annotation for object locally authenticated user.
- `cert_attribute` - (Optional) cert-attribute for object locally authenticated user.
- `clear_pwd_history` - (Optional) Allows the administrator to clear the password history of a locally-authenticated user. This is a trigger type attribute, So the value will reset to "no" once histry is cleared. Allowed values: "no", "yes". Default value: no.
- `description` - (Optional) Specifies a description of the policy definition.
- `email` - (Optional) The email address of the locally-authenticated user.
- `expiration` - (Optional) The expiration date of the locally-authenticated user account. The expires property must be enabled to activate an expiration date in format: YYYY-MM-DD HH:MM:SS. Default value: "never".
- `expires` - (Optional) A property to enable an expiration date for the locally-authenticated user account. Allowed values: "yes", "no". Default value is "no".
- `first_name` - (Optional) The first name of the locally-authenticated user.
- `last_name` - (Optional) The last name of the locally-authenticated user.
- `name_alias` - (Optional) Name alias for the locally-authenticated user.
- `otpenable` - (Optional) flag to enable OTP for the user. Allowed values: "yes", "no". Default value is "no".
- `otpkey` - (Optional) OTP-key for object user. Default value is "DISABLEDDISABLED".
- `phone` - (Optional) Phone number of the local user.
- `pwd` - (Optional) System user password
- `pwd_life_time` - (Optional) The lifetime of the local user password.Allowed values are in range of 0-3650. Default value is "0".
- `pwd_update_required` - (Optional) A boolean value indicating whether this account needs password update. Allowed values: "yes", "no". Default value is "no".
- `rbac_string` - (Optional) RBAC-string of the local user.
- `unix_user_id` - (Optional) The UNIX identifier of the local user.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Local User.

## Importing

An existing Local User can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_local_user.example <Dn>
```

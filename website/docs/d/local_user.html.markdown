---
layout: "aci"
page_title: "ACI: aci_local_user"
sidebar_current: "docs-aci-data-source-local_user"
description: |-
  Data source for ACI Local User
---

# aci_local_user

Data source for ACI Local User

## Example Usage

```hcl
data "aci_local_user" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) Name of the local user.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Local User.
- `account_status` - (Optional) The status of the locally-authenticated user account.
- `annotation` - (Optional) Annotation for object local user.
- `cert_attribute` - (Optional) certificate-attribute for object local user.
- `clear_pwd_history` - (Optional) Allows the administrator to clear the password history of a locally-authenticated user. 
- `description` - (Optional) Specifies a description of the policy definition.
- `email` - (Optional) The email address of the locally-authenticated user.
- `expiration` - (Optional) The expiration date of the locally-authenticated user account. The expires property must be enabled to activate an expiration date in format: YYYY-MM-DD HH:MM:SS.
- `expires` - (Optional) A property to enable an expiration date for the locally-authenticated user account.
- `first_name` - (Optional) The first name of the locally-authenticated user.
- `last_name` - (Optional) The last name of the locally-authenticated user.
- `name_alias` - (Optional) Name alias for the locally-authenticated user.
- `otpenable` - (Optional) flag to enable OTP for the user.
- `otpkey` - (Optional) OTP-key for object user.
- `phone` - (Optional) Phone number of the local user.
- `pwd` - (Optional) System user password.
- `pwd_life_time` - (Optional) The lifetime of the local user password.
- `pwd_update_required` - (Optional) A boolean value indicating whether this account needs password update.
- `rbac_string` - (Optional) RBAC-string of the local user.
- `unix_user_id` - (Optional) The UNIX identifier of the local user.

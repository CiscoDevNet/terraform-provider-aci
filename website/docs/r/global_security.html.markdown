---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_global_security"
sidebar_current: "docs-aci-resource-user_management"
description: |-
  Manages ACI Global Security
---

# aci_global_security #

Manages ACI Global Security

## API Information ##

* `Class` - aaaUserEp | aaaPwdProfile | aaaBlockLoginProfile | pkiWebTokenData
* `Distinguished Named` - uni/userext | uni/userext/pwdprofile | uni/userext/blockloginp | uni/userext/pkiext/webtokendata

## GUI Information ##

* `Location` - Admin -> AAA -> Security


## Example Usage ##

```hcl
resource "aci_global_security" "example" {
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "global_security"
  pwd_strength_check = "yes"
  change_count = "2"
  change_during_interval = "enable"
  change_interval = "48"
  expiration_warn_time = "15"
  history_count = "5"
  no_change_interval = "24"
  block_duration = "60"
  enable_login_block = "disable"
  max_failed_attempts = "5"
  max_failed_attempts_window = "5"
  maximum_validity_period = "24"
  session_record_flags = ["login", "logout", "refresh"]
  ui_idle_timeout_seconds = "1200"
  webtoken_timeout_seconds = "600"
  relation_aaa_rs_to_user_ep = aci_global_security.example2.id
}
```

## Argument Reference ##

* `annotation` - (Optional) Annotation of object Global Security.
* `description` - (Optional) Description of object Global Security.
* `name_alias` - (Optional) Name alias of object Global Security.
* `pwd_strength_check` - (Optional) Password Strength Check.The password strength check specifies if the system enforces the strength of the user password. Allowed values are "no" and "yes". Type: String.
* `change_count` - (Optional) Number of Password Changes in Interval.The number of password changes allowed within the change interval. Allowed range is 0-10.
* `change_during_interval` - (Optional) Password Policy.The change count/change interval policy selector. This property enables you to select an option for enforcing password change. Allowed values are "disable" and "enable". Type: String.
* `change_interval` - (Optional) Change Interval in Hours.A time interval for limiting the number of password changes. Allowed range is 0-745.
* `expiration_warn_time` - (Optional) Password Expiration Warn Time in Days.A warning period before password expiration.
A warning will be displayed when a user logs in within this number of days of an impending password expiration. Allowed range is 0-30.
* `history_count` - (Optional) Password History Count.How many retired passwords are stored in a user's password history. Allowed range is 0-15.
* `no_change_interval` - (Optional) No Password Change Interval in Hours.A minimum period after a password change before the user can change the password again. Allowed range is 0-745.
* `block_duration` - (Optional) Duration in minutes for which login should be blocked.Duration in minutes for which future logins should be blocked Allowed range is 1-1440.
* `enable_login_block` - (Optional) Enable blocking of user logins after failed attempts. Allowed values are "disable" and "enable". Type: String.
* `max_failed_attempts` - (Optional) Maximum continuous failed logins before blocking user.max failed login attempts before blocking user login Allowed range is 1-15.
* `max_failed_attempts_window` - (Optional) Time period for maximum continuous failed logins.times in minutes for max login failures to occur before blocking the user Allowed range is 1-720.
* `maximum_validity_period` - (Optional) Maximum Validity Period in hours.The maximum validity period for a webt oken. Allowed range is 4-24.
* `session_record_flags` - (Optional) Session Recording Options.Enables or disables a refresh in the session records. Allowed values are "login", "logout", "refresh". Type: List.
* `ui_idle_timeout_seconds` - (Optional) GUI Idle Timeout in Seconds.The maximum interval time the GUI remains idle before login needs to be refreshed. Allowed range is 60-65525.
* `webtoken_timeout_seconds` - (Optional) Timeout in Seconds.The web token timeout interval. Allowed range is 300-9600.

* `relation_aaa_rs_to_user_ep` - (Optional) Represents the relation to a Global Security (class aaaUserEp).  Type: String.

## Importing ##

An existing UserManagement can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_global_security.example <Dn>
```
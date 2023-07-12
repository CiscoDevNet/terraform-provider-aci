---
subcategory: -
layout: "aci"
page_title: "ACI: aci_user_profile"
sidebar_current: "docs-aci-data-source-user_profile"
description: |-
  Data source for ACI User Profile
---

# aci_user_profile #

Data source for ACI User Profile


## API Information ##

* `Class` - snmpUserP
* `Distinguished Name` - uni/fabric/snmppol-{name}/user-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
data "aci_user_profile" "example" {
  snmppolicy_dn  = aci_snmppolicy.example.id
  name  = "example"
}
```

## Argument Reference ##

* `snmppolicy_dn` - (Required) Distinguished name of the parent SNMPPolicy object.
* `name` - (Required) Name of the User Profile object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the User Profile.
* `annotation` - (Optional) Annotation of the User Profile object.
* `name_alias` - (Optional) Name Alias of the User Profile object.
* `auth_key` - (Optional) Authentication Key. The authentication key for the user profile. The key can be any case-sensitive alphanumeric string up to 64 chars.
* `auth_type` - (Optional) Authentication Type. The authentication type for the user profile. The authentication type is a message authentication code (MAC) that is used between two parties sharing a secret key to validate information transmitted between them. HMAC (Hash MAC) is based on cryptographic hash functions. It can be used in combination with any iterated cryptographic hash function. HMAC MD5 and HMAC SHA1 are two constructs of the HMAC using the MD5 hash function and the SHA1 hash function. HMAC also uses a secret key for calculation and verification of the message authentication values.
* `priv_key` - (Optional) Privacy Key. The privacy key for the user profile.
* `priv_type` - (Optional) Privacy. The encryption type for the user profile.

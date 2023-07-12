---
subcategory: "Fabric Policies"
layout: "aci"
page_title: "ACI: aci_snmp_user"
sidebar_current: "docs-aci-resource-snmp_user"
description: |-
  Manages ACI SNMP User
---

# aci_user_profile #

Manages ACI User Profile

## API Information ##

* `Class` - snmpUserP
* `Distinguished Name` - uni/fabric/snmppol-{name}/user-{name}

## GUI Information ##

* `Location` - Fabric > Fabric Policies > Policies > Pod > SNMP > {snmp_policy} > SNMP V3 Users


## Example Usage ##

```hcl
resource "aci_snmp_user" "example" {
  snmp_policy_dn  = aci_snmppolicy.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  authorization_key = 
  authorization_type = "hmac-sha1-96"

  name_alias = 
  privacy_key = 
  privacy_type = "none"
}
```

## Argument Reference ##

* `snmp_policy_dn` - (Required) Distinguished name of the parent SNMPPolicy object.
* `name` - (Required) Name of the User Profile object.
* `annotation` - (Optional) Annotation of the User Profile object.
* `name_alias` - (Optional) Name Alias of the User Profile object.
* `authorization_key` - (Optional) Authentication Key.The authentication key for the user profile. The key can be any case-sensitive alphanumeric string up to 64 chars.
* `authorization_type` - (Optional) Authentication Type.The authentication type for the user profile. The authentication type is a message authentication code (MAC) that is used between two parties sharing a secret key to validate information transmitted between them. HMAC (Hash MAC) is based on cryptographic hash functions. It can be used in combination with any iterated cryptographic hash function. HMAC MD5 and HMAC SHA1 are two constructs of the HMAC using the MD5 hash function and the SHA1 hash function. HMAC also uses a secret key for calculation and verification of the message authentication values. Allowed values are "hmac-md5-96", "hmac-sha1-96", "hmac-sha2-224", "hmac-sha2-256", "hmac-sha2-384", "hmac-sha2-512", and default value is "hmac-sha1-96". Type: String.

* `privacy_key` - (Optional) Privacy Key.The privacy key for the user profile.
* `privacy_type` - (Optional) Privacy.The encryption type for the user profile. Allowed values are "aes-128", "des", "none", and default value is "none". Type: String.


## Importing ##

An existing UserProfile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_user_profile.example <Dn>
```
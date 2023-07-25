---
subcategory: "Fabric Policies"
layout: "aci"
page_title: "ACI: aci_snmp_user"
sidebar_current: "docs-aci-data-source-snmp_user"
description: |-
  Data source for ACI SNMP User
---

# aci_user_profile #

Data source for ACI SNMP User


## API Information ##

* `Class` - snmpUserP
* `Distinguished Name` - uni/fabric/snmppol-{snmp_policy_name}/user-{name}

## GUI Information ##

* `Location` - Fabric > Fabric Policies > Policies > Pod > SNMP > {snmp_policy} > SNMP V3 Users


## Example Usage ##

```hcl
data "aci_snmp_user" "example" {
  snmp_policy_dn  = "uni/fabric/snmppol-default"
  name            = "example"
}
```

## Argument Reference ##

* `snmp_policy_dn` - (Required) Distinguished name of the parent SNMP Policy object. Type: String.
* `name` - (Required) Name of the SNMP User object. Type: String.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the SNMP User. Type: String.
* `annotation` - (Read-Only) Annotation of the SNMP User object.
* `name_alias` - (Read-Only) Name Alias of the SNMP User object.
* `authorization_type` - (Read-Only) Authorization Type. The authorization type for the SNMP user. The authorization type is a message authentication code (MAC) that is used between two parties sharing a secret key to validate information transmitted between them. HMAC (Hash MAC) is based on cryptographic hash functions. It can be used in combination with any iterated cryptographic hash function. HMAC MD5 and HMAC SHA1 are two constructs of the HMAC using the MD5 hash function and the SHA1 hash function. HMAC also uses a secret key for calculation and verification of the message authentication values.
* `privacy_type` - (Read-Only) Privacy Type. The encryption type for the SNMP user.

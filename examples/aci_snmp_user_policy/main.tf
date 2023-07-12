terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_snmp_user_policy" "foo_snmp_user" {
    snmp_policy_dn = "uni/fabric/snmppol-default"
    authorization_key = "testing123"
    authorization_type = "hmac-md5-96"
    name     = "Greg"
    privacy_key = "my_privacy_key"
    privacy_type = "aes-128"
}
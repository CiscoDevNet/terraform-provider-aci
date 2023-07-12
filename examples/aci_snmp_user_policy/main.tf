terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "ansible_github_ci"
  password = "sJ94G92#8dq2hx*K4qh"
  url      = "https://173.36.219.70"
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
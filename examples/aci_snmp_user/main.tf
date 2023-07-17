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
  url      = "https://173.36.219.68"
  insecure = true
}

resource "aci_snmp_user" "foo_snmp_user" {
    snmp_policy_dn      = "uni/fabric/snmppol-default"
    name                = "Greg"
    authorization_key   = "my_authorization_key"
    authorization_type  = "hmac-md5-96"
    privacy_key         = "my_privacy_key"
    privacy_type        = "aes-128"
}

resource "aci_snmp_user" "foo_snmp_user_1" {
    snmp_policy_dn      = "uni/fabric/snmppol-default"
    name                = "George"
    authorization_key   = "my_authorization_key"
    authorization_type  = "hmac-md5-96"
}
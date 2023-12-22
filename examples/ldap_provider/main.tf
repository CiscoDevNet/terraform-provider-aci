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

resource "aci_ldap_provider" "example" {
  name                 = "example"
  type                 = "duo"
  description          = "from terraform"
  annotation           = "example_annotation"
  name_alias           = "example_name_alias"
  ssl_validation_level = "strict"
  attribute            = "CiscoAvPair"
  basedn               = "CN=Users,DC=host,DC=com"
  enable_ssl           = "yes"
  filter               = "sAMAccountName=$userid"
  key                  = "example_key_value"
  monitor_server       = "enabled"
  monitoring_password  = "example_monitoring_password"
  monitoring_user      = "example_monitoring_user_value"
  port                 = "389"
  retries              = "1"
  rootdn               = "CN=admin,CN=Users,DC=host,DC=com"
  timeout              = "30"
}
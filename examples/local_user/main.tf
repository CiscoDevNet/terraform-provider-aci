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
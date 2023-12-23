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

resource "aci_global_security" "example" {
  annotation                 = "orchestrator:terraform"
  description                = "from terraform"
  name_alias                 = "global_security"
  pwd_strength_check         = "yes"
  change_count               = "2"
  change_during_interval     = "enable"
  change_interval            = "48"
  expiration_warn_time       = "15"
  history_count              = "5"
  no_change_interval         = "24"
  block_duration             = "60"
  enable_login_block         = "disable"
  max_failed_attempts        = "5"
  max_failed_attempts_window = "5"
  maximum_validity_period    = "24"
  session_record_flags       = ["login", "logout", "refresh"]
  ui_idle_timeout_seconds    = "1200"
  webtoken_timeout_seconds   = "600"
}
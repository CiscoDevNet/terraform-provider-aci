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

resource "aci_configuration_import_policy" "example" {
  name                   = "import_pol"
  admin_st               = "untriggered"
  annotation             = "example"
  description            = "from terraform"
  fail_on_decrypt_errors = "yes"
  file_name              = "file.tar.gz"
  import_mode            = "best-effort"
  import_type            = "replace"
  name_alias             = "example"
  snapshot               = "no"
}
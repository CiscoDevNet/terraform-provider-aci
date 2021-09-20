terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_file_remote_path" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  auth_type = "usePassword"
  host = "cisco.com"
  protocol = "sftp"
  remote_path = "example_remote_path"
  remote_port = "0"
  user_name = "example_user_name"
  user_passwd = "password"
  name_alias = "example_name_alias"
  description = "from terraform"
}
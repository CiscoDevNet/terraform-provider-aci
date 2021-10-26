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

resource "aci_stp_if_pol" "demo_stp_if_pol" {
  name        = "demo1"
  description = "This was created by terraform"
  ctrl        = ["bpdu-filter"]
}

resource "aci_stp_if_pol" "demo_stp_if_pol2" {
  name        = "demo2"
  description = "This was created by terraform"
  ctrl        = ["bpdu-filter", "bpdu-guard"]
}
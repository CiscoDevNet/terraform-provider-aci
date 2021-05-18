provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_stp_if_pol" "demo_stp_if_pol" {
  name        = "demo1"
  description = "This was created by terraform"
  ctrl = ["bpdu-filter"]
}

resource "aci_stp_if_pol" "demo_stp_if_pol2" {
  name        = "demo2"
  description = "This was created by terraform"
  ctrl = ["bpdu-filter", "bpdu-guard"]
}
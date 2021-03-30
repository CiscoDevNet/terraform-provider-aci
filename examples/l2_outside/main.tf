
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}


resource "aci_l2_outside" "example" {

  tenant_dn  = aci_tenant.example.id
  name  = "demo_l2_outside"
  annotation  = "example"
  name_alias  = "example"
  target_dscp = "AF11"

}
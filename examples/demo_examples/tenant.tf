resource "aci_tenant" "demotenant" {
  name        = "test_tf_tenant"
  description = "This tenant is created by terraform ACI provider"
}

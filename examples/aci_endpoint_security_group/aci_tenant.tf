resource "aci_tenant" "terraform_tenant" {
    name        = "tf_tenant"
    description = "This tenant is created by terraform"
}

resource "aci_tenant" "tenant_export_cons" {
    name        = "tenant_export_cons"
    description = "This tenant exports contracts to tf_tenant"
}
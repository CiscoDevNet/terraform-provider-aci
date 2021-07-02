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
resource "aci_lacp_policy" "test_lacp" {
  name = "tf_test_lacp"
  description = "From Terraform"
	annotation  = "tag_lacp"
	ctrl        = ["susp-individual"]
	max_links   = "16"
	min_links   = "1"
	mode        = "off"
	name_alias  = "alias_lacp"
}

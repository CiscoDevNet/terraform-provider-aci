# In order to add the cloud subnets in the azure cloud networking, hub networking needs to be disabled explicitly
resource "aci_cloud_template_region_detail" "hub_network" {
  parent_dn      = "uni/tn-infra/infranetwork-default/intnetwork-default/provider-azure-region-westus"
  hub_networking = "disable"
}
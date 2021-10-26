module "aci" {
  source = "./aci_resources"

}

module "vcenter" {
  source = "./vmware_resources"

  aci_tenant_name              = module.aci.aci_tenant_name
  aci_application_profile_name = module.aci.aci_application_profile_name
  aci_epg_name                 = module.aci.aci_epg_name

}


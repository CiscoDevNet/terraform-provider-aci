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

resource "aci_tenant" "tf_tenant" {
  name = "tf_tenant"
}

# VRF setup part
resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "vrf-1"
}

# AAA Domain setup part
resource "aci_aaa_domain" "aaa_domain_1" {
  name = "aaa_domain_1"
}

resource "aci_aaa_domain" "aaa_domain_2" {
  name = "aaa_domain_2"
}

# Third-Party Firewall
resource "aci_cloud_l4_l7_third_party_device" "cloud_third_party_fw" {
  tenant_dn        = aci_tenant.tf_tenant.id
  name             = "cloud_third_party_fw"
  active_active    = "no"
  context_aware    = "single-Context"
  device_type      = "CLOUD"
  function_type    = "GoTo"
  instance_count   = "2"
  is_copy          = "no"
  is_instantiation = "no"
  managed          = "yes"
  mode             = "legacy-Mode"
  prom_mode        = "no"
  service_type     = "FW"
  target_mode      = "unspecified"
  trunking         = "no"

  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id,
    aci_aaa_domain.aaa_domain_2.id
  ]
  relation_cloud_rs_ldev_to_ctx = aci_vrf.vrf1.id

  interface_selectors {
    allow_all = "no"
    name      = "Interface_1"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_1_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_1_ep_2"
    }
  }
  interface_selectors {
    allow_all = "no"
    name      = "Interface_2"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_2_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_2_ep_2"
    }
  }
  interface_selectors {
    allow_all = "no"
    name      = "Interface_3"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_3_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_3_ep_2"
    }
  }
  interface_selectors {
    allow_all = "no"
    name      = "Interface_4"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_4_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_4_ep_2"
    }
  }
}


# Third-Party Load Balancer
resource "aci_cloud_l4_l7_third_party_device" "cloud_third_party_lb" {
  tenant_dn    = aci_tenant.tf_tenant.id
  name         = "cloud_third_party_lb"
  service_type = "ADC"

  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id
  ]
  relation_cloud_rs_ldev_to_ctx = aci_vrf.vrf1.id

  interface_selectors {
    allow_all = "no"
    name      = "Interface_1"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_1_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_1_ep_2"
    }
  }
}

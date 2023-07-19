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

resource "aci_tenant" "footenant" {
  description = "sample aci_tenant from terraform"
  name        = "tf_tenant"
}

resource "aci_igmp_interface_policy" "example_igmp" {
  tenant_dn                  = aci_tenant.footenant.id
  name                       = "example_igmp"
  group_timeout              = "260"
  last_member_count          = "2"
  last_member_response_time  = "1"
  querier_timeout            = "255"
  query_interval             = "125"
  robustness_variable        = "2"
  response_interval          = "10"
  startup_query_count        = "2"
  startup_query_interval     = "31"
  version                    = "v2"
  maximum_mulitcast_entries  = "20"
  reserved_mulitcast_entries = "9"
  state_limit_route_map      = "uni/tn-common/rtmap-test-1"
  report_policy_route_map    = "uni/tn-common/rtmap-test-2"
  static_report_route_map    = "uni/tn-common/rtmap-test-3"
}

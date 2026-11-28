package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type stateUpgradeInput struct {
	Version int64
	State   string
}

func testResourceStateUpgrade(t *testing.T, instance resource.Resource, testCase test.TestCase) {
	t.Helper()
	ctx := context.Background()
	input := testCase.Input.(stateUpgradeInput)
	upgrader, exists := instance.(resource.ResourceWithUpgradeState).UpgradeState(ctx)[input.Version]
	if !exists {
		t.Fatalf("no upgrader registered for version %d", input.Version)
	}
	var current resource.SchemaResponse
	instance.Schema(ctx, resource.SchemaRequest{}, &current)
	if current.Diagnostics.HasError() {
		t.Fatal(current.Diagnostics)
	}
	priorValue, err := tftypes.ValueFromJSON([]byte(input.State), upgrader.PriorSchema.Type().TerraformType(ctx))
	if err != nil {
		t.Fatal(err)
	}
	request := resource.UpgradeStateRequest{State: &tfsdk.State{Schema: upgrader.PriorSchema, Raw: priorValue}}
	response := resource.UpgradeStateResponse{State: tfsdk.State{Schema: current.Schema}}
	upgrader.StateUpgrader(ctx, request, &response)
	if response.Diagnostics.HasError() {
		t.Fatal(response.Diagnostics)
	}
	expected, err := tftypes.ValueFromJSON([]byte(testCase.Expected.(string)), current.Schema.Type().TerraformType(ctx))
	if err != nil {
		t.Fatal(err)
	}
	if !response.State.Raw.Equal(expected) {
		t.Fatalf("upgraded state:\n%s\nexpected:\n%s", response.State.Raw, expected)
	}
}

func TestResourceStateUpgradeHistoricalFixtures(t *testing.T) {
	// Handwritten JSON locks the published field names and agreed conversions independently of generator definitions.
	// These are synthetic historical fixtures, not captured customer state. Omitted object fields decode to typed nulls.
	type input struct {
		Resource resource.Resource
		Version  int64
		State    string
	}
	testCases := []test.TestCase{
		{
			Name: "tenant_named_relation_and_frozen_values",
			Input: input{Resource: NewFvTenantResource(), Version: 1, State: `{
  "id": "uni/tn-migration_tenant", "name": "migration_tenant",
  "annotation": "tenant_annotation", "description": "tenant_description", "name_alias": "tenant_alias",
  "relation_fv_rs_tenant_mon_pol": "uni/tn-common/monepg-observe-policy",
  "relation_fv_rs_tn_deny_rule": ["uni/tn-common/taboo-deny_a", "uni/tn-common/taboo-deny_b"]
}`},
			Expected: `{
  "id": "uni/tn-migration_tenant", "name": "migration_tenant",
  "annotation": "tenant_annotation", "description": "tenant_description", "name_alias": "tenant_alias",
  "relation_fv_rs_tenant_mon_pol": "uni/tn-common/monepg-observe-policy",
  "relation_fv_rs_tn_deny_rule": ["uni/tn-common/taboo-deny_a", "uni/tn-common/taboo-deny_b"],
  "relation_to_monitoring_policy": {"monitoring_policy_name": "observe-policy"}
}`,
		},
		{
			Name: "tenant_null_relation",
			Input: input{Resource: NewFvTenantResource(), Version: 1, State: `{
  "name": "null_relation_tenant", "relation_fv_rs_tenant_mon_pol": null
}`},
			Expected: `{
  "name": "null_relation_tenant", "relation_fv_rs_tenant_mon_pol": null,
  "relation_to_monitoring_policy": {"monitoring_policy_name": ""}
}`,
		},
		{
			Name: "tenant_empty_relation",
			Input: input{Resource: NewFvTenantResource(), Version: 1, State: `{
  "name": "empty_relation_tenant", "relation_fv_rs_tenant_mon_pol": ""
}`},
			Expected: `{
  "name": "empty_relation_tenant", "relation_fv_rs_tenant_mon_pol": "",
  "relation_to_monitoring_policy": {"monitoring_policy_name": ""}
}`,
		},
		{
			Name: "tenant_missing_optional_fields",
			Input: input{Resource: NewFvTenantResource(), Version: 1, State: `{
  "id": "uni/tn-minimal_tenant", "name": "minimal_tenant"
}`},
			Expected: `{
  "id": "uni/tn-minimal_tenant", "name": "minimal_tenant",
  "relation_to_monitoring_policy": {"monitoring_policy_name": ""}
}`,
		},
		{
			Name: "application_profile_parent_and_property_aliases",
			Input: input{Resource: NewFvApResource(), Version: 1, State: `{
  "id": "uni/tn-parent_tenant/ap-migration_app", "name": "migration_app",
  "tenant_dn": "uni/tn-parent_tenant", "prio": "level2",
  "annotation": "app_annotation", "description": "app_description", "name_alias": "app_alias"
}`},
			Expected: `{
  "id": "uni/tn-parent_tenant/ap-migration_app", "name": "migration_app",
  "tenant_dn": "uni/tn-parent_tenant", "parent_dn": "uni/tn-parent_tenant",
  "prio": "level2", "priority": "level2",
  "annotation": "app_annotation", "description": "app_description", "name_alias": "app_alias",
  "relation_to_monitoring_policy": {"monitoring_policy_name": ""}
}`,
		},
		{
			Name: "stp_list_order_and_duplicates",
			Input: input{Resource: NewStpIfPolResource(), Version: 1, State: `{
  "id": "uni/infra/ifPol-migration_stp", "name": "migration_stp",
  "ctrl": ["bpdu-guard", "bpdu-filter", "bpdu-guard"], "description": "stp_description"
}`},
			Expected: `{
  "id": "uni/infra/ifPol-migration_stp", "name": "migration_stp", "description": "stp_description",
  "ctrl": ["bpdu-filter", "bpdu-guard", "bpdu-guard"], "interface_controls": ["bpdu-filter", "bpdu-guard"]
}`,
		},
		{
			Name: "epg_distinct_child_blocks_and_dn_collections",
			Input: input{Resource: NewFvAEPgResource(), Version: 1, State: `{
  "id": "uni/tn-tenant/ap-app/epg-migration_epg", "name": "migration_epg",
  "application_profile_dn": "uni/tn-tenant/ap-app",
  "relation_fv_rs_cons": ["uni/tn-common/brc-contract_a", "uni/tn-common/brc-contract_b"],
  "relation_fv_rs_sec_inherited": ["uni/tn-tenant/ap-app/epg-master_a", "uni/tn-tenant/ap-app/epg-master_b"],
  "relation_fv_rs_node_att": [
    {"node_dn": "topology/pod-1/node-101", "encap": "vlan-101", "mode": "regular", "deployment_immediacy": "immediate", "description": "leaf_a"},
    {"node_dn": "topology/pod-1/node-102", "encap": "vlan-102", "mode": "native", "deployment_immediacy": "lazy", "description": "leaf_b"}
  ]
}`},
			Expected: `{
  "id": "uni/tn-tenant/ap-app/epg-migration_epg", "name": "migration_epg",
  "application_profile_dn": "uni/tn-tenant/ap-app", "parent_dn": "uni/tn-tenant/ap-app",
  "relation_fv_rs_cons": ["uni/tn-common/brc-contract_a", "uni/tn-common/brc-contract_b"],
  "relation_to_consumed_contracts": [{"contract_name": "contract_a"}, {"contract_name": "contract_b"}],
  "relation_fv_rs_sec_inherited": ["uni/tn-tenant/ap-app/epg-master_a", "uni/tn-tenant/ap-app/epg-master_b"],
  "relation_to_contract_masters": [{"target_dn": "uni/tn-tenant/ap-app/epg-master_a"}, {"target_dn": "uni/tn-tenant/ap-app/epg-master_b"}],
  "relation_fv_rs_node_att": [
    {"node_dn": "topology/pod-1/node-101", "encap": "vlan-101", "mode": "regular", "deployment_immediacy": "immediate", "description": "leaf_a"},
    {"node_dn": "topology/pod-1/node-102", "encap": "vlan-102", "mode": "native", "deployment_immediacy": "lazy", "description": "leaf_b"}
  ],
  "relation_to_static_leafs": [
    {"target_dn": "topology/pod-1/node-101", "encapsulation": "vlan-101", "mode": "regular", "deployment_immediacy": "immediate", "description": "leaf_a"},
    {"target_dn": "topology/pod-1/node-102", "encapsulation": "vlan-102", "mode": "native", "deployment_immediacy": "lazy", "description": "leaf_b"}
  ],
  "relation_to_monitoring_policy": {"monitoring_policy_name": ""},
  "relation_to_bridge_domain": {"bridge_domain_name": ""},
  "relation_to_custom_qos_policy": {"custom_qos_policy_name": ""},
  "relation_to_data_plane_policing_policy": {"data_plane_policing_policy_name": ""},
  "relation_to_trust_control_policy": {"trust_control_policy_name": ""},
  "relation_to_imported_contracts": [], "relation_to_fibre_channel_paths": [],
  "relation_to_intra_epg_contracts": [], "relation_to_taboo_contracts": [], "relation_to_provided_contracts": []
}`,
		},
		{
			Name: "bridge_domain_named_blocks_and_ip_alias",
			Input: input{Resource: NewFvBDResource(), Version: 1, State: `{
  "id": "uni/tn-tenant/BD-migration_bd", "name": "migration_bd", "tenant_dn": "uni/tn-tenant",
  "ll_addr": "FE80::0001",
  "relation_fv_rs_bd_to_netflow_monitor_pol": [
    {"flt_type": "ipv4", "tn_netflow_monitor_pol_name": "uni/tn-common/monitorpol-flow_a"},
    {"flt_type": "ipv6", "tn_netflow_monitor_pol_name": "uni/tn-common/monitorpol-flow_b"}
  ]
}`},
			Expected: `{
  "id": "uni/tn-tenant/BD-migration_bd", "name": "migration_bd", "tenant_dn": "uni/tn-tenant", "parent_dn": "uni/tn-tenant",
  "ll_addr": "FE80::0001", "link_local_ipv6_address": "FE80::0001",
  "relation_fv_rs_bd_to_netflow_monitor_pol": [
    {"flt_type": "ipv4", "tn_netflow_monitor_pol_name": "uni/tn-common/monitorpol-flow_a"},
    {"flt_type": "ipv6", "tn_netflow_monitor_pol_name": "uni/tn-common/monitorpol-flow_b"}
  ],
  "relation_to_netflow_monitor_policies": [
    {"filter_type": "ipv4", "netflow_monitor_policy_name": "flow_a"},
    {"filter_type": "ipv6", "netflow_monitor_policy_name": "flow_b"}
  ],
  "relation_to_monitoring_policy": {"monitoring_policy_name": ""},
  "relation_to_end_point_retention_policy": {"end_point_retention_policy_name": ""},
  "relation_to_first_hop_security_policy": {"first_hop_security_policy_name": ""},
  "relation_to_neighbor_discovery_interface_policy": {"neighbor_discovery_interface_policy_name": ""},
  "relation_to_dhcp_relay_policy": {"dhcp_relay_policy_name": ""},
  "relation_to_vrf": {"vrf_name": ""},
  "relation_to_igmp_snooping_policy": {"igmp_snooping_policy_name": ""},
  "relation_to_mld_snooping_policy": {"mld_snooping_policy_name": ""},
  "relation_to_l3_outsides": []
}`,
		},
		{
			Name: "netflow_set_to_singleton_and_unchanged_children",
			Input: input{Resource: NewNetflowMonitorPolResource(), Version: 0, State: `{
  "id": "uni/infra/monitorpol-flow_monitor", "name": "flow_monitor", "parent_dn": "uni/infra",
  "annotation": "monitor_annotation",
  "relation_to_netflow_record": [{"annotation": "record_annotation", "netflow_record_policy_name": "record_one"}],
  "relation_to_netflow_exporters": [
    {"annotation": "exporter_annotation_a", "netflow_exporter_policy_name": "exporter_a"},
    {"annotation": "exporter_annotation_b", "netflow_exporter_policy_name": "exporter_b"}
  ]
}`},
			Expected: `{
  "id": "uni/infra/monitorpol-flow_monitor", "name": "flow_monitor", "parent_dn": "uni/infra",
  "annotation": "monitor_annotation",
  "relation_to_netflow_record": {"annotation": "record_annotation", "netflow_record_policy_name": "record_one"},
  "relation_to_netflow_exporters": [
    {"annotation": "exporter_annotation_a", "netflow_exporter_policy_name": "exporter_a"},
    {"annotation": "exporter_annotation_b", "netflow_exporter_policy_name": "exporter_b"}
  ]
}`,
		},
		{
			Name: "vmm_null_is_not_an_empty_string",
			Input: input{Resource: NewVmmDomPResource(), Version: 1, State: `{
  "id": "uni/vmmp-VMware/dom-migration_domain", "name": "migration_domain", "provider_profile_dn": "uni/vmmp-VMware", "arp_learning": null
}`},
			Expected: `{
  "id": "uni/vmmp-VMware/dom-migration_domain", "name": "migration_domain", "provider_profile_dn": "uni/vmmp-VMware",
  "parent_dn": "uni/vmmp-VMware", "arp_learning": null,
  "relation_to_ip_address_pool": {"target_dn": null}, "relation_to_vlan_pool": {"target_dn": null},
  "relation_to_multicast_pool": {"target_dn": null}, "relation_to_lacp_enhanced_lag_policy": {"target_dn": null}
}`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			fixture := testCase.Input.(input)
			testResourceStateUpgrade(t, fixture.Resource, test.TestCase{
				Input: stateUpgradeInput{Version: fixture.Version, State: fixture.State}, Expected: testCase.Expected,
			})
		})
	}
}

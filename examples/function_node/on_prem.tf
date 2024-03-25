
# terraform plan for on-prem APICs

resource "aci_tenant" "tf_tenant" {
  name        = "tf_ansible_test"
  description = "Terraform tenant"
}

resource "aci_vrf" "tf_vrf" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "tf_vrf"
}
resource "aci_physical_domain" "tf_domain" {
  name = "tf_phy_domain"
}

#  Create devices

# Create a Load Balancer Device
resource "aci_l4_l7_device" "tf_device1" {
  tenant_dn                            = aci_tenant.tf_tenant.id
  name                                 = "tf_device_lb"
  service_type                         = "ADC"
  relation_vns_rs_al_dev_to_phys_dom_p = aci_physical_domain.tf_domain.id
}

# Create a Firewall device
resource "aci_l4_l7_device" "tf_device2" {
  tenant_dn                            = aci_tenant.tf_tenant.id
  name                                 = "tf_device_fw"
  service_type                         = "FW"
  relation_vns_rs_al_dev_to_phys_dom_p = aci_physical_domain.tf_domain.id
}
#  Create a Copy devices
resource "aci_l4_l7_device" "tf_device3" {
  tenant_dn                            = aci_tenant.tf_tenant.id
  name                                 = "tf_device_copy_1"
  is_copy                              = "yes"
  service_type                         = "COPY"
  function_type                        = "None"
  relation_vns_rs_al_dev_to_phys_dom_p = aci_physical_domain.tf_domain.id
}

resource "aci_l4_l7_device" "tf_device4" {
  tenant_dn                            = aci_tenant.tf_tenant.id
  name                                 = "tf_device_copy_2"
  is_copy                              = "yes"
  service_type                         = "COPY"
  function_type                        = "None"
  relation_vns_rs_al_dev_to_phys_dom_p = aci_physical_domain.tf_domain.id
}

# 1. Create a L4-L7 Service Graph Template with copy node
resource "aci_l4_l7_service_graph_template" "tf_sg_1" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "tf_sg_1"
}

# Add a copy node to the service graph template 
# Copy device must be unmanaged
resource "aci_function_node" "tf_copy" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_1.id
  name                            = "CP1"
  func_template_type              = "OTHER"
  is_copy                         = "yes"
  managed                         = "no"
  func_type                       = "None"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device3.id
}

# Create L4-L7 Service Graph connection with template and copy node.
resource "aci_connection" "t1-t2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_1.id
  name                            = "C1"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.tf_sg_1.term_cons_dn,
    aci_l4_l7_service_graph_template.tf_sg_1.term_prov_dn,
  ]
  relation_vns_rs_abs_copy_connection = [aci_function_node.tf_copy.conn_copy_dn]
}

# 2. Create a L4-L7 Service Graph Template with one firewall node
resource "aci_l4_l7_service_graph_template" "tf_sg_2" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "tf_sg_2"
}

resource "aci_function_node" "tf_node_fw" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_2.id
  name                            = "N1"
  func_template_type              = "FW_TRANS"
  is_copy                         = "no"
  managed                         = "no"
  func_type                       = "GoTo"
  routing_mode                    = "unspecified"
  sequence_number                 = "3"
  share_encap                     = "yes"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device1.id
}

resource "aci_connection" "sg2-t1-n1" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_2.id
  name                            = "C1"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.tf_sg_2.term_cons_dn,
    aci_function_node.tf_node_fw.conn_consumer_dn,
  ]
}

# Create L4-L7 Service Graph connection with template T1 and the first node N0.
resource "aci_connection" "sg2-n1-t2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_2.id
  name                            = "C2"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_function_node.tf_node_fw.conn_provider_dn,
    aci_l4_l7_service_graph_template.tf_sg_2.term_prov_dn,
  ]
}


# 3. Create a L4-L7 Service Graph Template with two nodes (firewall and load balancer)
resource "aci_l4_l7_service_graph_template" "tf_sg_3" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "tf_sg_3"
}

resource "aci_function_node" "tf_node_lb3" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_3.id
  name                            = "N1"
  func_template_type              = "ADC_ONE_ARM"
  is_copy                         = "no"
  managed                         = "no"
  func_type                       = "GoTo"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device1.id
}
resource "aci_function_node" "tf_node_fw3" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_3.id
  name                            = "N2"
  func_template_type              = "FW_TRANS"
  is_copy                         = "no"
  managed                         = "no"
  func_type                       = "GoTo"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device2.id
}

# Create L4-L7 Service Graph connection with template T1 and the first node N0.
resource "aci_connection" "sg3-t1-n1" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_3.id
  name                            = "C1"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.tf_sg_3.term_cons_dn,
    aci_function_node.tf_node_lb3.conn_consumer_dn,
  ]
}

resource "aci_connection" "sg3-n1-n2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_3.id
  name                            = "C2"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_function_node.tf_node_lb3.conn_provider_dn,
    aci_function_node.tf_node_fw3.conn_consumer_dn,
  ]
}

resource "aci_connection" "sg3-n2-t2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_3.id
  name                            = "C3"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_function_node.tf_node_fw3.conn_provider_dn,
    aci_l4_l7_service_graph_template.tf_sg_3.term_prov_dn,
  ]
}


# 4. Create a L4-L7 Service Graph Template with two nodes (firewall and load balancer) and 2 copy nodes
resource "aci_l4_l7_service_graph_template" "tf_sg_4" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "tf_sg_4"
}

# Create a Load Balancer Node
resource "aci_function_node" "tf_node1_lb" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "N1"
  func_template_type              = "ADC_ONE_ARM"
  is_copy                         = "no"
  managed                         = "no"
  func_type                       = "GoTo"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device1.id
}

# Create a Firewall Node
resource "aci_function_node" "tf_node2_fw" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "N2"
  func_template_type              = "FW_TRANS"
  is_copy                         = "no"
  managed                         = "no"
  func_type                       = "GoTo"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device2.id
}

# Add 3 Copy Nodes
resource "aci_function_node" "tf_copy1" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "CP1"
  func_template_type              = "OTHER"
  is_copy                         = "yes"
  managed                         = "no"
  func_type                       = "None"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device3.id
}

resource "aci_function_node" "tf_copy2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "CP2"
  func_template_type              = "OTHER"
  is_copy                         = "yes"
  managed                         = "no"
  func_type                       = "None"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device3.id
}

resource "aci_function_node" "tf_copy3" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "CP3"
  func_template_type              = "OTHER"
  is_copy                         = "yes"
  managed                         = "no"
  func_type                       = "None"
  relation_vns_rs_node_to_l_dev   = aci_l4_l7_device.tf_device4.id
}

# Create L4-L7 Service Graph connection with template T1 and the first node N1 with copy node CP1 attached between them.
resource "aci_connection" "t1-n1-cp1" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "C1"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.tf_sg_4.term_cons_dn,
    aci_function_node.tf_node1_lb.conn_consumer_dn,
  ]
  relation_vns_rs_abs_copy_connection = [aci_function_node.tf_copy1.conn_copy_dn]
}

resource "aci_connection" "n1-n2-cp2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "C2"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_function_node.tf_node1_lb.conn_provider_dn,
    aci_function_node.tf_node2_fw.conn_consumer_dn,
  ]
  relation_vns_rs_abs_copy_connection = [aci_function_node.tf_copy2.conn_copy_dn]
}

resource "aci_connection" "n2-t2-cp3" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.tf_sg_4.id
  name                            = "C3"
  adj_type                        = "L2"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "no"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_function_node.tf_node2_fw.conn_provider_dn,
    aci_l4_l7_service_graph_template.tf_sg_4.term_prov_dn,
  ]
  relation_vns_rs_abs_copy_connection = [aci_function_node.tf_copy3.conn_copy_dn]
}
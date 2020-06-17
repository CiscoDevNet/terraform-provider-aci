

provider "aci" {
  username = "admin"
  password = "cisco123"
  url      = "https://192.168.10.102"
  insecure = true
}

resource "aci_vlan_pool" "PHYS-VLAN-POOL" {
  name  = "PHYS-VLAN-POOL"
  alloc_mode  = "static"
}
resource "aci_physical_domain" "PhyDom" {
  depends_on = [aci_vlan_pool.PHYS-VLAN-POOL]
  name  = "PhyDom"
  relation_infra_rs_vlan_ns = aci_vlan_pool.PHYS-VLAN-POOL.id
}
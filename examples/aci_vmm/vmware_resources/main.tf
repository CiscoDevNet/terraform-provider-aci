provider "vsphere" {
  user                 = var.vsphere_user
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_server
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = var.vsphere_datacenter
}

data "vsphere_network" "vm1_net" {
  name          = format("%v|%v|%v", var.aci_tenant_name, var.aci_application_profile_name, var.aci_epg_name)
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_network" "vm2_net" {
  name          = format("%v|%v|%v", var.aci_tenant_name, var.aci_application_profile_name, var.aci_epg_name)
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_datastore" "ds" {
  name          = var.vsphere_datastore
  datacenter_id = data.vsphere_datacenter.dc.id
}

# data "vsphere_compute_cluster" "cl" {
#   name          = var.vsphere_compute_cluster
#   datacenter_id = data.vsphere_datacenter.dc.id
# }

data "vsphere_host" "host" {
  name          = var.vsphere_host_name
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_virtual_machine" "template" {
  name          = var.vsphere_template
  datacenter_id = data.vsphere_datacenter.dc.id
}

resource "vsphere_virtual_machine" "aci_vm1" {
  count        = 1
  name         = var.aci_vm1_name
  datastore_id = data.vsphere_datastore.ds.id

  num_cpus = 8
  memory   = 24576
  guest_id = data.vsphere_virtual_machine.template.guest_id

  #   resource_pool_id = data.vsphere_compute_cluster.cl.resource_pool_id

  resource_pool_id = data.vsphere_host.host.resource_pool_id
  scsi_type = data.vsphere_virtual_machine.template.scsi_type
  disk {
    label = "disk0"
    size  = data.vsphere_virtual_machine.template.disks.0.size
  }
  disk {
    unit_number = 1
    label       = "disk1"
    size        = 40
  }
  folder = var.folder
  network_interface {
    network_id   = data.vsphere_network.vm1_net.id
    adapter_type = data.vsphere_virtual_machine.template.network_interface_types[0]
  }
  clone {
    template_uuid = data.vsphere_virtual_machine.template.id

    customize {
      linux_options {
        host_name = var.aci_vm1_name
        domain    = var.domain_name
      }

      network_interface {
        ipv4_address = var.aci_vm1_address
        ipv4_netmask = "22"
      }

      ipv4_gateway    = var.gateway
      dns_server_list = var.dns_list
      dns_suffix_list = var.dns_search
    }
  }
}

resource "vsphere_virtual_machine" "aci_vm2" {
  count = 1
  name  = var.aci_vm2_name

  #   resource_pool_id = data.vsphere_compute_cluster.cl.resource_pool_id
  resource_pool_id = data.vsphere_host.host.resource_pool_id
  datastore_id     = data.vsphere_datastore.ds.id

  num_cpus = 8
  memory   = 24576
  guest_id = data.vsphere_virtual_machine.template.guest_id

  scsi_type = data.vsphere_virtual_machine.template.scsi_type

  disk {
    label = "disk0"
    size  = data.vsphere_virtual_machine.template.disks.0.size
  }

  disk {
    unit_number = 1
    label       = "disk1"
    size        = 40
  }

  folder = var.folder

  network_interface {
    network_id   = data.vsphere_network.vm2_net.id
    adapter_type = data.vsphere_virtual_machine.template.network_interface_types[0]
  }

  clone {
    template_uuid = data.vsphere_virtual_machine.template.id

    customize {
      linux_options {
        host_name = var.aci_vm2_name
        domain    = var.domain_name
      }

      network_interface {
        ipv4_address = var.aci_vm2_address
        ipv4_netmask = "22"
      }

      ipv4_gateway    = var.gateway
      dns_server_list = var.dns_list
      dns_suffix_list = var.dns_search
    }
  }
}

---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_service_epg"
sidebar_current: "docs-aci-data-source-cloud_service_epg"
description: |-
  Data source for ACI Cloud Service EPG
---

# aci_cloud_service_epg #

Data source for ACI Cloud Service EPG
Note: This datasource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudSvcEPg
* `Distinguished Name` - uni/tn-{tenant_name}/cloudapp-{application_name}/cloudsvcepg-{name}

## GUI Information ##

* `Location` - Application Management -> EPGs


## Example Usage ##

```hcl
data "aci_cloud_service_epg" "example" {
  cloud_application_container_dn  = aci_cloud_applicationcontainer.example.id
  name                            = "example"
}
```

## Argument Reference ##

* `cloud_application_container_dn` - (Required) Distinguished name of the parent Cloud Application Container object. Type: String.
* `name` - (Required) Name of the Cloud Service EPG object. Type: String.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Cloud Service EPG. Type: String.
* `annotation` - (Read-Only) Annotation of the Cloud Service EPG object. Type: String.
* `name_alias` - (Read-Only) Name Alias of the Cloud Service EPG object. Type: String.
* `access_type` - (Read-Only) This refers to the type of connectivity to the service. It could be a public or private connectivity. Type: String.
* `azure_private_endpoint` - (Read-Only) Naming for Azure Private Endpoint created from the Service Cloud EPG. The set of variables supported by the naming override is the same as those supported in the global naming policy. However, there is no mandatory variable enforced by validations. Type: String.
* `custom_service_type` - (Read-Only) A custom service type used when this EPG is use as custom service EPG with public or private access. As an e.g. this is used to provide the service tag for an Azure service. Type: string
* `deployment_type` - (Read-Only) cloud service deployment type.deploymentType refers to the type of deployment of the service. It could be a PaaS service, a PaaS service deployed in the managed VNET/VPC, a SaaS service consumed or a SaaS service offered. Type: String.
* `flood_on_encap` - (Read-Only) Handling of L2 Multicast/Broadcast and Link-Layer traffic at EPG level. Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Type: String.
* `label_match_criteria` - (Read-Only) Provider Label Match Criteria. Type: String.
* `preferred_group_member` - (Read-Only)  Represents the parameter used to determine if the Cloud Service EPG is part of a group that does not have a contract for communication. Type: String.
* `prio` - (Read-Only) The QoS priority class identifier. Type: String.
* `cloud_service_epg_type` - (Read-Only) The specific service type of the object or component. Type: String.
* `relation_cloud_rs_cloud_epg_ctx` - (Read-Only) Represents the relation to a VRF (class fvCtx). Type: String.
* `relation_fv_rs_cons` - (Read-Only) A block representing the relation to a Contract Consumer (class vzBrCP). Type: Block.
  * `prio` - (Read-Only) The system class determines the quality of service and priority for the consumer traffic. Type: String.
  * `target_dn` - (Read-Only) The distinguished name of the target. Type: String
* `relation_fv_rs_cons_if` - (Read-Only) A block representing the relation to a Contract Interface (class vzCPIf) for which the EPG will be a consumer. Type: Block.
  * `prio` - (Read-Only) The contract interface priority. Type: String.
  * `target_dn` - (Read-Only) The distinguished name of the target. Type: String
* `relation_fv_rs_cust_qos_pol` - (Read-Only) Represents the relation to a Custom QOS Policy (class qosCustomPol) that enables different levels of service to be assigned to network traffic, including specifications for the Differentiated Services Code Point (DSCP) value(s) and the 802.1p priority. Type: String.
* `relation_fv_rs_graph_def` - (Read-Only) Represents the relation to a Graph Container (class vzGraphCont). Type: Block.
* `relation_fv_rs_intra_epg` - (Read-Only) Represents the relation to an Intra EPG Contract (class vzBrCP). It also represents that the EPG is moving from "allow all within epg" mode to a "deny all within epg" mode. The only type of traffic allowed between EPs in this EPG is the one specified by contracts EPG associated with this relation. Type: Block.
* `relation_fv_rs_prot_by` - (Read-Only) Represents the relation to a Taboo Contract Association (class vzTaboo) where the EPG will be a provider and consumer to the contract. Type: Block.
* `relation_fv_rs_prov` - (Read-Only) A block representing the relation to a Contract Provider (class vzBrCP). This relationship allows the EPG to be the contract's provider. Type: Block.
  * `label_match_criteria` - (Read-Only) The matched EPG type. Type: String.
  * `prio` - (Read-Only) The system class determines the quality of service and priority for the consumer traffic. Type: String.
  * `target_dn` - (Read-Only) The distinguished name of the target. Type: String
* `relation_fv_rs_prov_def` - (Read-Only) Represents the relation to a Contract EPG Container (class vzCtrctEPgCont). Type: Block.
* `relation_fv_rs_sec_inherited` - (Read-Only) Represents the relation to a Security inheritance (class fvEPg) where the EPG is inheriting security configuration from another EPG. Type: Block.
---
subcategory: - "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_service_epg"
sidebar_current: "docs-aci-resource-cloud_service_epg"
description: |-
  Manages ACI Cloud Service EPG
---

# aci_cloud_service_epg #

Manages ACI Cloud Service EPG

## API Information ##

* `Class` - cloudSvcEPg
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_cloud_service_epg" "example" {
  cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.example.id
  name  = "example"
  access_type = "Public"
  annotation = "orchestrator:terraform"
  az_private_endpoint = 
  custom_svc_type = 
  deployment_type = "Unknown"
  flood_on_encap = "disabled"
  match_t = "AtleastOne"

  name_alias = 
  pref_gr_memb = "exclude"
  prio = "unspecified"
  cloud_service_epg_type = "Unknown"

  cloud_rs_cloud_epg_ctx = aci_resource.example.id

  fv_rs_cons {
    prio = "unspecified"
    target_dn = aci_resource.example.id
  }

  fv_rs_cons_if {
    prio = "unspecified"
    target_dn = aci_resource.example.id
  }

  fv_rs_cust_qos_pol = aci_resource.example.id

  fv_rs_graph_def = [aci_resource.example.id]

  fv_rs_intra_epg = [aci_resource.example.id]

  fv_rs_prot_by = [aci_resource.example.id]

  fv_rs_prov {
    match_t = "AtleastOne"
    prio = "unspecified"
    target_dn = aci_resource.example.id
  }

  fv_rs_prov_def = [aci_resource.example.id]

  fv_rs_sec_inherited = [aci_resource.example.id]
}
```

## Argument Reference ##

* `cloud_applicationcontainer_dn` - (Required) Distinguished name of the parent Cloud Application container object.
* `name` - (Required) Name of the Cloud Service EPG object.
* `annotation` - (Optional) Annotation of the Cloud Service EPg object.
* `name_alias` - (Optional) Name Alias of the Cloud Service EPg object.
* `access_type` - (Optional) cloud service access type.accessType refers to the type of connectivity to the
                    service. It could be a public or private connectivity. Allowed values are "Private", "Public", "PublicAndPrivate", "Unknown", and default value is "Public". Type: String.

* `az_private_endpoint` - (Optional) Naming for Azure Private Endpoint created from the SvcEPg.naming override for any Azure Private Endpoint
                     that gets created from this SvcEPg. The set of variable
                     supported by the naming override is the same of
                     those supported in the global naming policy. However, there
                     is no mandatory variable enforced by validations.
* `custom_svc_type` - (Optional) A Custom Service Type String.a custom service type used when this EPg is used
                     as custom svc EPg with public or private access.
                     As an e.g. this is used to provide the service tag for
                     an Azure service.
* `deployment_type` - (Optional) cloud service deployment type.deploymentType refers to the type of deployment of the
                    service. It could be a PaaS service, a PaaS service
                    deployed in the managed VNET/VPC, a SaaS service consumed
                    or a SaaS service offered. Allowed values are "CloudNative", "CloudNativeManaged", "Third-party", "Third-partyManaged", "Unknown", and default value is "Unknown". Type: String.
* `flood_on_encap` - (Optional) Handling of L2 Multicast/Broadcast and Link-Layer traffic at EPG level.Control at EPG level if the traffic L2
                     Multicast/Broadcast and Link Local Layer should
                     be flooded only on ENCAP or based on bridg-domain
                     settings Allowed values are "disabled", "enabled", and default value is "disabled". Type: String.
* `match_t` - (Optional) Provider Label Match Criteria.The provider label match criteria. Allowed values are "All", "AtleastOne", "AtmostOne", "None", and default value is "AtleastOne". Type: String.

* `pref_gr_memb` - (Optional) Preferred Group Member.Represents parameter used to determine
                    if EPg is part of a group that does not
                    a contract for communication Allowed values are "exclude", "include", and default value is "exclude". Type: String.
* `prio` - (Optional) QOS Class.The QoS priority class identifier. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
* `cloud_service_epg_type` - (Optional) cloud service type.The specific type of the object or component. Allowed values are "Azure-ADDS", "Azure-AksCluster", "Azure-ApiManagement", "Azure-ContainerRegistry", "Azure-CosmosDB", "Azure-Databricks", "Azure-KeyVault", "Azure-Redis", "Azure-SqlServer", "Azure-Storage", "Azure-StorageBlob", "Azure-StorageFile", "Azure-StorageQueue", "Azure-StorageTable", "Custom", "Unknown", and default value is "Unknown". Type: String.
* `relation_cloudrs_cloud_epg_ctx` - (Optional) Represents the relation to a Relationship to the fvCtx of belonging (class fvCtx). Relationship to the fvCtx of belonging Type: String.
* `relation_fvrs_cons` - (Optional) A block representing the relation to a Contract Consumer (class vzBrCP). The Consumer contract profile information. Type: Block.
  * `prio` - (Optional) prio. The system class determines the quality of service and priority for the consumer traffic. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  * `target_dn` - (Required) The distinguished name of the target. Type: String
* `relation_fvrs_cons_if` - (Optional) A block representing the relation to a Contract Interface (class vzCPIf). A contract for which the EPG will be a consumer. Type: Block.
  * `prio` - (Optional) prio. The contract interface priority. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  * `target_dn` - (Required) The distinguished name of the target. Type: String
* `relation_fvrs_cust_qos_pol` - (Optional) Represents the relation to a Custom QOS Policy (class qosCustomPol). A source relation to a custom QoS policy that enables different levels of service to be assigned to network traffic, including specifications for the Differentiated Services Code Point (DSCP) value(s) and the 802.1p Dot1p priority. This is an internal object. Type: String.
* `relation_fvrs_graph_def` - (Optional) Represents the relation to a FvRsGraphDef (class vzGraphCont). A source relation to the graph container. Type: List.
* `relation_fvrs_intra_epg` - (Optional) Represents the relation to a Intra EPg Contract (class vzBrCP). Intra EPg contract:
                      Represents that the EPg is moving from "allow all within epg" mode
                      to a "deny all within epg" mode.
                      The only type of traffic allowed between EPs in this EPg is the one
                      specified by contracts EPg associates to with this relation. Type: List.
* `relation_fvrs_prot_by` - (Optional) Represents the relation to a Taboo Contract Association (class vzTaboo). The taboo contract for which the EPG will be a provider and consumer. Type: List.
* `relation_fvrs_prov` - (Optional) A block representing the relation to a Contract Provider (class vzBrCP). A contract for which the EPG will be a provider. Type: Block.
  * `match_t` - (Optional) matchT. The matched EPG type. Allowed values are "All", "AtleastOne", "AtmostOne", "None", and default value is "AtleastOne". Type: String.
  * `prio` - (Optional) prio. The system class determines the quality of service and priority for the consumer traffic. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  * `target_dn` - (Required) The distinguished name of the target. Type: String
* `relation_fvrs_prov_def` - (Optional) Represents the relation to a Contract EPG Container (class vzCtrctEPgCont). A source relation to a binary contract profile. Type: List.
* `relation_fvrs_sec_inherited` - (Optional) Represents the relation to a Security inheritance (class fvEPg). Represents that the EPg is inheriting security configuration from another EPg Type: List.


## Importing ##

An existing CloudServiceEPg can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_service_epg.example <Dn>
```
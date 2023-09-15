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
* `Distinguished Name` - uni/tn-{tenant_name}/cloudapp-{application_name}/cloudsvcepg-{name}

## GUI Information ##

* `Location` - Application Management -> EPGs


## Example Usage ##

```hcl
resource "aci_cloud_service_epg" "example" {
  cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.example.id
  name                           = "example"
  access_type                    = "Public"
  annotation                     = "orchestrator:terraform"
  deployment_type                = "CloudNative"
  flood_on_encap                 = "disabled"
  label_match_criteria           = "AtleastOne"
  cloud_service_epg_type         = "Azure-SqlServer"
  cloud_rs_cloud_epg_ctx         = aci_vrf.example.id
}
```

## Argument Reference ##

* `cloud_applicationcontainer_dn` - (Required) Distinguished name of the parent Cloud Application container object.
* `name` - (Required) Name of the Cloud Service EPG object.
* `annotation` - (Optional) Annotation of the Cloud Service EPG object.
* `name_alias` - (Optional) Name Alias of the Cloud Service EPG object.
* `access_type` - (Optional) This refers to the type of connectivity to the service. It could be a public or private connectivity. Allowed values are "Private", "Public", "PublicAndPrivate", "Unknown", and default value is "Public". Type: String.
* `azure_private_endpoint` - (Optional) Naming for Azure Private Endpoint created from the Service Cloud EPG. Naming override for any Azure Private Endpoint that gets created from this service EPG. The set of variable supported by the naming override is the same of those supported in the global naming policy. However, there is no mandatory variable enforced by validations. Type: String.
* `custom_service_type` - (Optional) Custom Service type. A custom service type used when this EPG is use as custom service EPG with public or private access. As an e.g. this is used to provide the service tag for an Azure service. Type: string
* `deployment_type` - (Optional) cloud service deployment type.deploymentType refers to the type of deployment of the service. It could be a PaaS service, a PaaS service deployed in the managed VNET/VPC, a SaaS service consumed or a SaaS service offered. Allowed values are "CloudNative", "CloudNativeManaged", "Third-party","Third-partyManaged", "Unknown", and default value is "Unknown". Type: String.
* `flood_on_encap` - (Optional) Handling of L2 Multicast/Broadcast and Link-Layer traffic at EPG level. Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled", "enabled", and default value is "disabled". Type: String.
* `label_match_criteria` - (Optional) Provider Label Match Criteria. Allowed values are "All", "AtleastOne", "AtmostOne", "None", and default value is "AtleastOne". Type: String.
* `preferred_group_member` - (Optional) Represents parameter used to determine if the Cloud Service EPG is part of a group that does not a contract for communication. Allowed values are "exclude", "include", and default value is "exclude". Type: String.
* `prio` - (Optional) The QoS priority class identifier. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
* `cloud_service_epg_type` - (Optional) The specific service type of the object or component. Allowed values are "Azure-ADDS", "Azure-AksCluster", "Azure-ApiManagement", "Azure-ContainerRegistry", "Azure-CosmosDB", "Azure-Databricks", "Azure-KeyVault", "Azure-Redis", "Azure-SqlServer", "Azure-Storage", "Azure-StorageBlob", "Azure-StorageFile", "Azure-StorageQueue", "Azure-StorageTable", "Custom", "Unknown", and default value is "Unknown". Type: String.
* `relation_cloud_rs_cloud_epg_ctx` - (Optional) Represents the relation to a Relationship to the VRF of belonging (class fvCtx). Type: String.
* `relation_fv_rs_cons` - (Optional) A block representing the relation to a Contract Consumer (class vzBrCP). The Consumer contract profile information. Type: Block.
  * `prio` - (Optional) The system class determines the quality of service and priority for the consumer traffic. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  * `target_dn` - (Required) The distinguished name of the target. Type: String
* `relation_fv_rs_cons_if` - (Optional) A block representing the relation to a Contract Interface (class vzCPIf). A contract for which the EPG will be a consumer. Type: Block.
  * `prio` - (Optional) The contract interface priority. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  * `target_dn` - (Required) The distinguished name of the target. Type: String
* `relation_fv_rs_cust_qos_pol` - (Optional) Represents the relation to a Custom QOS Policy (class qosCustomPol). A source relation to a custom QoS policy that enables different levels of service to be assigned to network traffic, including specifications for the Differentiated Services Code Point (DSCP) value(s) and the 802.1p Dot1p priority. This is an internal object. Type: String.
* `relation_fv_rs_graph_def` - (Optional) Represents the relation to a Graph Container (class vzGraphCont). Type: List.
* `relation_fv_rs_intra_epg` - (Optional) Represents the relation to an Intra EPG Contract (class vzBrCP). It also represents that the EPG is moving from "allow all within epg" mode to a "deny all within epg" mode. The only type of traffic allowed between EPs in this EPG is the one specified by contracts EPG associated with this relation. Type: List.
* `relation_fv_rs_prot_by` - (Optional) Represents the relation to a Taboo Contract Association (class vzTaboo). The taboo contract for which the EPG will be a provider and consumer. Type: List.
* `relation_fv_rs_prov` - (Optional) A block representing the relation to a Contract Provider (class vzBrCP). This relationship allows the EPG to be the contract's provider. Type: Block.
  * `match_t` - (Optional) The matched EPG type. Allowed values are "All", "AtleastOne", "AtmostOne", "None", and default value is "AtleastOne". Type: String.
  * `prio` - (Optional) The system class determines the quality of service and priority for the consumer traffic. Allowed values are "level1", "level2", "level3", "level4", "level5", "level6", "unspecified", and default value is "unspecified". Type: String.
  * `target_dn` - (Required) The distinguished name of the target. Type: String
* `relation_fv_rs_prov_def` - (Optional) Represents the relation to a Contract EPG Container (class vzCtrctEPgCont). A source relation to a binary contract profile. Type: List.
* `relation_fv_rs_sec_inherited` - (Optional) Represents the relation to a Security inheritance (class fvEPg) where the EPG is inheriting security configuration from another EPG. Type: List.


## Importing ##

An existing Cloud Service EPG can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_service_epg.example <Dn>
```

Starting in Terraform version 1.5, an existing Cloud Service EPG can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "<Dn>"
  to = aci_cloud_service_epg.example
}
```
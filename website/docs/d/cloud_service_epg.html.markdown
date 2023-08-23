---
subcategory: -
layout: "aci"
page_title: "ACI: aci_cloud_service_epg"
sidebar_current: "docs-aci-data-source-cloud_service_epg"
description: |-
  Data source for ACI Cloud Service EPg
---

# aci_cloud_service_epg #

Data source for ACI Cloud Service EPg
Note: This resource is supported in Cloud APIC only.

## API Information ##

* `Class` - cloudSvcEPg
* `Distinguished Name` - uni/tn-{name}/cloudapp-{name}/cloudsvcepg-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
data "aci_cloud_service_epg" "example" {
  cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.example.id
  name  = "example"
}
```

## Argument Reference ##

* `cloud_applicationcontainer_dn` - (Required) Distinguished name of the parent CloudApplicationcontainer object.
* `name` - (Required) Name of the Cloud Service EPg object.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the Cloud Service EPg.
* `annotation` - (Read-Only) Annotation of the Cloud Service EPg object.
* `name_alias` - (Read-Only) Name Alias of the Cloud Service EPg object.
* `access_type` - (Read-Only) cloud service access type. accessType refers to the type of connectivity to the
                    service. It could be a public or private connectivity.
* `az_private_endpoint` - (Read-Only) Naming for Azure Private Endpoint created from the SvcEPg. naming override for any Azure Private Endpoint
                     that gets created from this SvcEPg. The set of variable
                     supported by the naming override is the same of
                     those supported in the global naming policy. However, there
                     is no mandatory variable enforced by validations.
* `custom_svc_type` - (Read-Only) A Custom Service Type String. a custom service type used when this EPg is used
                     as custom svc EPg with public or private access.
                     As an e.g. this is used to provide the service tag for
                     an Azure service.
* `deployment_type` - (Read-Only) cloud service deployment type. deploymentType refers to the type of deployment of the
                    service. It could be a PaaS service, a PaaS service
                    deployed in the managed VNET/VPC, a SaaS service consumed
                    or a SaaS service offered.
* `flood_on_encap` - (Read-Only) Handling of L2 Multicast/Broadcast and Link-Layer traffic at EPG level. Control at EPG level if the traffic L2
                     Multicast/Broadcast and Link Local Layer should
                     be flooded only on ENCAP or based on bridg-domain
                     settings
* `match_t` - (Read-Only) Provider Label Match Criteria. The provider label match criteria.
* `pref_gr_memb` - (Read-Only) Preferred Group Member. Represents parameter used to determine
                    if EPg is part of a group that does not
                    a contract for communication
* `prio` - (Read-Only) QOS Class. The QoS priority class identifier.
* `cloud_service_epg_type` - (Read-Only) cloud service type. The specific type of the object or component.

* `relation_cloudrs_cloud_epg_ctx` - (Read-Only) Represents the relation to a Relationship to the fvCtx of belonging (class fvCtx). Relationship to the fvCtx of belonging Type: String.
* `relation_fvrs_cons` - (Read-Only) A list of maps representing the relation to a Contract Consumer (class vzBrCP). The Consumer contract profile information. Type: Block.
  * `prio` - (Read-Only) prio. The system class determines the quality of service and priority for the consumer traffic.. Type: String.
  * `target_dn` - (Read-Only) The distinguished name of the target. Type: String
* `relation_fvrs_cons_if` - (Read-Only) A list of maps representing the relation to a Contract Interface (class vzCPIf). A contract for which the EPG will be a consumer. Type: Block.
  * `prio` - (Read-Only) prio. The contract interface priority.. Type: String.
  * `target_dn` - (Read-Only) The distinguished name of the target. Type: String
* `relation_fvrs_cust_qos_pol` - (Read-Only) Represents the relation to a Custom QOS Policy (class qosCustomPol). A source relation to a custom QoS policy that enables different levels of service to be assigned to network traffic, including specifications for the Differentiated Services Code Point (DSCP) value(s) and the 802.1p Dot1p priority. This is an internal object. Type: String.
* `relation_fvrs_graph_def` - (Read-Only) Represents the relation to a FvRsGraphDef (class vzGraphCont). A source relation to the graph container. Type: List.
* `relation_fvrs_intra_epg` - (Read-Only) Represents the relation to a Intra EPg Contract (class vzBrCP). Intra EPg contract:
                      Represents that the EPg is moving from "allow all within epg" mode
                      to a "deny all within epg" mode.
                      The only type of traffic allowed between EPs in this EPg is the one
                      specified by contracts EPg associates to with this relation. Type: List.
* `relation_fvrs_prot_by` - (Read-Only) Represents the relation to a Taboo Contract Association (class vzTaboo). The taboo contract for which the EPG will be a provider and consumer. Type: List.
* `relation_fvrs_prov` - (Read-Only) A list of maps representing the relation to a Contract Provider (class vzBrCP). A contract for which the EPG will be a provider. Type: Block.
  * `match_t` - (Read-Only) matchT. The matched EPG type.. Type: String.
  * `prio` - (Read-Only) prio. The system class determines the quality of service and priority for the consumer traffic.. Type: String.
  * `target_dn` - (Read-Only) The distinguished name of the target. Type: String
* `relation_fvrs_prov_def` - (Read-Only) Represents the relation to a Contract EPG Container (class vzCtrctEPgCont). A source relation to a binary contract profile. Type: List.
* `relation_fvrs_sec_inherited` - (Read-Only) Represents the relation to a Security inheritance (class fvEPg). Represents that the EPg is inheriting security configuration from another EPg Type: List.
package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciOverridePCVPCPolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciOverridePCVPCPolicyGroupCreate,
		UpdateContext: resourceAciOverridePCVPCPolicyGroupUpdate,
		ReadContext:   resourceAciOverridePCVPCPolicyGroupRead,
		DeleteContext: resourceAciOverridePCVPCPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciOverridePCVPCPolicyGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"leaf_access_bundle_policy_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_channel_member": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to lacp:IfPol",
			},
		}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
	}
}

func getRemoteOverridePCVPCPolicyGroup(client *client.Client, dn string) (*models.OverridePCVPCPolicyGroup, error) {
	infraAccBndlSubgrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	infraAccBndlSubgrp := models.OverridePCVPCPolicyGroupFromContainer(infraAccBndlSubgrpCont)
	if infraAccBndlSubgrp.DistinguishedName == "" {
		return nil, fmt.Errorf("OverridePCVPCPolicyGroup %s not found", dn)
	}
	return infraAccBndlSubgrp, nil
}

func setOverridePCVPCPolicyGroupAttributes(infraAccBndlSubgrp *models.OverridePCVPCPolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraAccBndlSubgrp.DistinguishedName)
	d.Set("description", infraAccBndlSubgrp.Description)
	infraAccBndlSubgrpMap, err := infraAccBndlSubgrp.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != infraAccBndlSubgrp.DistinguishedName {
		d.Set("leaf_access_bundle_policy_group_dn", "")
	} else {
		d.Set("leaf_access_bundle_policy_group_dn", GetParentDn(infraAccBndlSubgrp.DistinguishedName, fmt.Sprintf("/"+models.RninfraAccBndlSubgrp, infraAccBndlSubgrpMap["name"])))
	}

	if infraAccBndlSubgrpMap["annotation"] == "" {
		d.Set("annotation", "orchestrator:terraform")
	} else {
		d.Set("annotation", infraAccBndlSubgrpMap["annotation"])
	}

	d.Set("name", infraAccBndlSubgrpMap["name"])
	d.Set("name_alias", infraAccBndlSubgrpMap["nameAlias"])
	return d, nil
}

func getInfraRsLacpInterfacePolData(callType, parentDn string, client *client.Client) interface{} {
	log.Printf("[DEBUG] Beginning GET called by %s function for port channel member of %s", callType, parentDn)
	infraRsLacpInterfacePolData, err := client.ReadRelationinfraRsLacpInterfacePol(parentDn)
	if err == nil {
		log.Printf("[DEBUG] Finished GET called by %s successfully with result: %v", callType, infraRsLacpInterfacePolData)
	} else {
		log.Printf("[DEBUG] Error during GET operation for port channel member of %v: %v", parentDn, err)
		infraRsLacpInterfacePolData = nil
	}
	return infraRsLacpInterfacePolData
}

func resourceAciOverridePCVPCPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraAccBndlSubgrp, err := getRemoteOverridePCVPCPolicyGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setOverridePCVPCPolicyGroupAttributes(infraAccBndlSubgrp, d)
	if err != nil {
		return nil, err
	}

	infraRsLacpInterfacePolData := getInfraRsLacpInterfacePolData("Import", infraAccBndlSubgrp.DistinguishedName, aciClient)
	if infraRsLacpInterfacePolData != nil {
		d.Set("port_channel_member", fmt.Sprintf("uni/infra/lacpifp-%s", infraRsLacpInterfacePolData.(string)))
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOverridePCVPCPolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OverridePCVPCPolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LeafAccessBundlePolicyGroupDn := d.Get("leaf_access_bundle_policy_group_dn").(string)

	infraAccBndlSubgrpAttr := models.OverridePCVPCPolicyGroupAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccBndlSubgrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccBndlSubgrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraAccBndlSubgrpAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAccBndlSubgrpAttr.NameAlias = NameAlias.(string)
	}
	infraAccBndlSubgrp := models.NewOverridePCVPCPolicyGroup(fmt.Sprintf(models.RninfraAccBndlSubgrp, name), LeafAccessBundlePolicyGroupDn, desc, infraAccBndlSubgrpAttr)

	err := aciClient.Save(infraAccBndlSubgrp)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationToinfraRsLacpInterfacePol, ok := d.GetOk("port_channel_member"); ok {
		relationParam := relationToinfraRsLacpInterfacePol.(string)
		checkDns = append(checkDns, relationParam)
	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationToinfraRsLacpInterfacePol, ok := d.GetOk("port_channel_member"); ok {
		relationParam := relationToinfraRsLacpInterfacePol.(string)
		err = aciClient.CreateRelationinfraRsLacpInterfacePol(infraAccBndlSubgrp.DistinguishedName, infraAccBndlSubgrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraAccBndlSubgrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciOverridePCVPCPolicyGroupRead(ctx, d, m)
}
func resourceAciOverridePCVPCPolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OverridePCVPCPolicyGroup: Beginning Update")
	deleted := false
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LeafAccessBundlePolicyGroupDn := d.Get("leaf_access_bundle_policy_group_dn").(string)

	infraAccBndlSubgrpAttr := models.OverridePCVPCPolicyGroupAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAccBndlSubgrpAttr.Annotation = Annotation.(string)
	} else {
		infraAccBndlSubgrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraAccBndlSubgrpAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAccBndlSubgrpAttr.NameAlias = NameAlias.(string)
	}

	if d.HasChange("annotation") || d.HasChange("name_alias") || d.HasChange("description") {
		err := aciClient.DeleteOverridePCVPCPolicyGroup(name, GetMOName(LeafAccessBundlePolicyGroupDn))
		if err != nil {
			return diag.FromErr(err)
		}
		deleted = true
	}

	infraAccBndlSubgrp := models.NewOverridePCVPCPolicyGroup(fmt.Sprintf(models.RninfraAccBndlSubgrp, name), LeafAccessBundlePolicyGroupDn, desc, infraAccBndlSubgrpAttr)
	err := aciClient.Save(infraAccBndlSubgrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("port_channel_member") || deleted {
		_, newRelParam := d.GetChange("port_channel_member")
		checkDns = append(checkDns, newRelParam.(string))
	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("port_channel_member") || deleted {
		_, newRelParam := d.GetChange("port_channel_member")

		if !deleted {
			err = aciClient.DeleteRelationinfraRsLacpInterfacePol(infraAccBndlSubgrp.DistinguishedName)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		err = aciClient.CreateRelationinfraRsLacpInterfacePol(infraAccBndlSubgrp.DistinguishedName, infraAccBndlSubgrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraAccBndlSubgrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciOverridePCVPCPolicyGroupRead(ctx, d, m)
}

func resourceAciOverridePCVPCPolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	infraAccBndlSubgrp, err := getRemoteOverridePCVPCPolicyGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setOverridePCVPCPolicyGroupAttributes(infraAccBndlSubgrp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsLacpInterfacePolData := getInfraRsLacpInterfacePolData("Read", infraAccBndlSubgrp.DistinguishedName, aciClient)
	if infraRsLacpInterfacePolData != nil {
		d.Set("port_channel_member", fmt.Sprintf("uni/infra/lacpifp-%s", infraRsLacpInterfacePolData.(string)))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciOverridePCVPCPolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "infraAccBndlSubgrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

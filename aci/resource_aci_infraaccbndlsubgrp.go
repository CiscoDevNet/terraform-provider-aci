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

func resourceAciOverridePolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciOverridePolicyGroupCreate,
		UpdateContext: resourceAciOverridePolicyGroupUpdate,
		ReadContext:   resourceAciOverridePolicyGroupRead,
		DeleteContext: resourceAciOverridePolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciOverridePolicyGroupImport,
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

func getRemoteOverridePolicyGroup(client *client.Client, dn string) (*models.OverridePolicyGroup, error) {
	infraAccBndlSubgrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	infraAccBndlSubgrp := models.OverridePolicyGroupFromContainer(infraAccBndlSubgrpCont)
	if infraAccBndlSubgrp.DistinguishedName == "" {
		return nil, fmt.Errorf("OverridePolicyGroup %s not found", dn)
	}
	return infraAccBndlSubgrp, nil
}

func setOverridePolicyGroupAttributes(infraAccBndlSubgrp *models.OverridePolicyGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
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
	d.Set("annotation", infraAccBndlSubgrpMap["annotation"])
	d.Set("name", infraAccBndlSubgrpMap["name"])
	d.Set("name_alias", infraAccBndlSubgrpMap["nameAlias"])
	return d, nil
}

func resourceAciOverridePolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraAccBndlSubgrp, err := getRemoteOverridePolicyGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setOverridePolicyGroupAttributes(infraAccBndlSubgrp, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOverridePolicyGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OverridePolicyGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LeafAccessBundlePolicyGroupDn := d.Get("leaf_access_bundle_policy_group_dn").(string)

	infraAccBndlSubgrpAttr := models.OverridePolicyGroupAttributes{}

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
	infraAccBndlSubgrp := models.NewOverridePolicyGroup(fmt.Sprintf(models.RninfraAccBndlSubgrp, name), LeafAccessBundlePolicyGroupDn, desc, infraAccBndlSubgrpAttr)

	err := aciClient.Save(infraAccBndlSubgrp)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationToinfraRsLacpInterfacePol, ok := d.GetOk("port_channel_member"); ok {
		relationParam := relationToinfraRsLacpInterfacePol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsLacpInterfacePol, ok := d.GetOk("port_channel_member"); ok {
		relationParam := relationToinfraRsLacpInterfacePol.(string)
		err = aciClient.CreateRelationinfraRsLacpInterfacePol(infraAccBndlSubgrp.DistinguishedName, infraAccBndlSubgrpAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraAccBndlSubgrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciOverridePolicyGroupRead(ctx, d, m)
}
func resourceAciOverridePolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OverridePolicyGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LeafAccessBundlePolicyGroupDn := d.Get("leaf_access_bundle_policy_group_dn").(string)

	infraAccBndlSubgrpAttr := models.OverridePolicyGroupAttributes{}

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

	infraAccBndlSubgrp := models.NewOverridePolicyGroup(fmt.Sprintf(models.RninfraAccBndlSubgrp, name), LeafAccessBundlePolicyGroupDn, desc, infraAccBndlSubgrpAttr)
	infraAccBndlSubgrp.Status = "modified"

	err := aciClient.Save(infraAccBndlSubgrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("port_channel_member") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("port_channel_member")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("port_channel_member") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("port_channel_member")
		err = aciClient.DeleteRelationinfraRsLacpInterfacePol(infraAccBndlSubgrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsLacpInterfacePol(infraAccBndlSubgrp.DistinguishedName, infraAccBndlSubgrpAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraAccBndlSubgrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciOverridePolicyGroupRead(ctx, d, m)
}

func resourceAciOverridePolicyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	infraAccBndlSubgrp, err := getRemoteOverridePolicyGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setOverridePolicyGroupAttributes(infraAccBndlSubgrp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsLacpInterfacePolData, err := aciClient.ReadRelationinfraRsLacpInterfacePol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsLacpInterfacePol %v", err)
		d.Set("port_channel_member", "")
	} else {
		if _, ok := d.GetOk("port_channel_member"); ok {
			tfName := GetMOName(d.Get("port_channel_member").(string))
			if tfName != infraRsLacpInterfacePolData {
				d.Set("port_channel_member", "")
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciOverridePolicyGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

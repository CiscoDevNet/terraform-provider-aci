package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciSPANSourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSPANSourceGroupCreate,
		UpdateContext: resourceAciSPANSourceGroupUpdate,
		ReadContext:   resourceAciSPANSourceGroupRead,
		DeleteContext: resourceAciSPANSourceGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSPANSourceGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"enabled",
					"disabled",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_span_rs_src_grp_to_filter_grp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteSPANSourceGroup(client *client.Client, dn string) (*models.SPANSourceGroup, error) {
	spanSrcGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	spanSrcGrp := models.SPANSourceGroupFromContainer(spanSrcGrpCont)

	if spanSrcGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("SPANSourceGroup %s not found", spanSrcGrp.DistinguishedName)
	}

	return spanSrcGrp, nil
}

func setSPANSourceGroupAttributes(spanSrcGrp *models.SPANSourceGroup, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(spanSrcGrp.DistinguishedName)
	d.Set("description", spanSrcGrp.Description)
	// d.Set("tenant_dn", GetParentDn(spanSrcGrp.DistinguishedName))
	if dn != spanSrcGrp.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	spanSrcGrpMap, _ := spanSrcGrp.ToMap()

	d.Set("name", spanSrcGrpMap["name"])
	d.Set("tenant_dn", GetParentDn(spanSrcGrp.DistinguishedName, fmt.Sprintf("/srcgrp-%s", spanSrcGrpMap["name"])))
	d.Set("admin_st", spanSrcGrpMap["adminSt"])
	d.Set("annotation", spanSrcGrpMap["annotation"])
	d.Set("name_alias", spanSrcGrpMap["nameAlias"])
	return d
}

func resourceAciSPANSourceGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	spanSrcGrp, err := getRemoteSPANSourceGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	spanSrcGrpMap, _ := spanSrcGrp.ToMap()
	name := spanSrcGrpMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/srcgrp-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled := setSPANSourceGroupAttributes(spanSrcGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSPANSourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SPANSourceGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	spanSrcGrpAttr := models.SPANSourceGroupAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		spanSrcGrpAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanSrcGrpAttr.Annotation = Annotation.(string)
	} else {
		spanSrcGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSrcGrpAttr.NameAlias = NameAlias.(string)
	}
	spanSrcGrp := models.NewSPANSourceGroup(fmt.Sprintf("srcgrp-%s", name), TenantDn, desc, spanSrcGrpAttr)

	err := aciClient.Save(spanSrcGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTospanRsSrcGrpToFilterGrp, ok := d.GetOk("relation_span_rs_src_grp_to_filter_grp"); ok {
		relationParam := relationTospanRsSrcGrpToFilterGrp.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTospanRsSrcGrpToFilterGrp, ok := d.GetOk("relation_span_rs_src_grp_to_filter_grp"); ok {
		relationParam := relationTospanRsSrcGrpToFilterGrp.(string)
		err = aciClient.CreateRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(spanSrcGrp.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(spanSrcGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSPANSourceGroupRead(ctx, d, m)
}

func resourceAciSPANSourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SPANSourceGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	spanSrcGrpAttr := models.SPANSourceGroupAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		spanSrcGrpAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanSrcGrpAttr.Annotation = Annotation.(string)
	} else {
		spanSrcGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSrcGrpAttr.NameAlias = NameAlias.(string)
	}
	spanSrcGrp := models.NewSPANSourceGroup(fmt.Sprintf("srcgrp-%s", name), TenantDn, desc, spanSrcGrpAttr)

	spanSrcGrp.Status = "modified"

	err := aciClient.Save(spanSrcGrp)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_span_rs_src_grp_to_filter_grp") {
		_, newRelParam := d.GetChange("relation_span_rs_src_grp_to_filter_grp")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_span_rs_src_grp_to_filter_grp") {
		_, newRelParam := d.GetChange("relation_span_rs_src_grp_to_filter_grp")
		err = aciClient.DeleteRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(spanSrcGrp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(spanSrcGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(spanSrcGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSPANSourceGroupRead(ctx, d, m)

}

func resourceAciSPANSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	spanSrcGrp, err := getRemoteSPANSourceGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSPANSourceGroupAttributes(spanSrcGrp, d)

	spanRsSrcGrpToFilterGrpData, err := aciClient.ReadRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation spanRsSrcGrpToFilterGrp %v", err)
		d.Set("relation_span_rs_src_grp_to_filter_grp", "")

	} else {
		if _, ok := d.GetOk("relation_span_rs_src_grp_to_filter_grp"); ok {
			tfName := d.Get("relation_span_rs_src_grp_to_filter_grp").(string)
			if tfName != spanRsSrcGrpToFilterGrpData {
				d.Set("relation_span_rs_src_grp_to_filter_grp", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSPANSourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "spanSrcGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciSPANSourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciSPANSourceGroupCreate,
		Update: resourceAciSPANSourceGroupUpdate,
		Read:   resourceAciSPANSourceGroupRead,
		Delete: resourceAciSPANSourceGroupDelete,

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
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	schemaFilled := setSPANSourceGroupAttributes(spanSrcGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSPANSourceGroupCreate(d *schema.ResourceData, m interface{}) error {
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
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSrcGrpAttr.NameAlias = NameAlias.(string)
	}
	spanSrcGrp := models.NewSPANSourceGroup(fmt.Sprintf("srcgrp-%s", name), TenantDn, desc, spanSrcGrpAttr)

	err := aciClient.Save(spanSrcGrp)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTospanRsSrcGrpToFilterGrp, ok := d.GetOk("relation_span_rs_src_grp_to_filter_grp"); ok {
		relationParam := relationTospanRsSrcGrpToFilterGrp.(string)
		err = aciClient.CreateRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(spanSrcGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_span_rs_src_grp_to_filter_grp")
		d.Partial(false)

	}

	d.SetId(spanSrcGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSPANSourceGroupRead(d, m)
}

func resourceAciSPANSourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
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
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSrcGrpAttr.NameAlias = NameAlias.(string)
	}
	spanSrcGrp := models.NewSPANSourceGroup(fmt.Sprintf("srcgrp-%s", name), TenantDn, desc, spanSrcGrpAttr)

	spanSrcGrp.Status = "modified"

	err := aciClient.Save(spanSrcGrp)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_span_rs_src_grp_to_filter_grp") {
		_, newRelParam := d.GetChange("relation_span_rs_src_grp_to_filter_grp")
		err = aciClient.DeleteRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(spanSrcGrp.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(spanSrcGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_span_rs_src_grp_to_filter_grp")
		d.Partial(false)

	}

	d.SetId(spanSrcGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSPANSourceGroupRead(d, m)

}

func resourceAciSPANSourceGroupRead(d *schema.ResourceData, m interface{}) error {
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

func resourceAciSPANSourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "spanSrcGrp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

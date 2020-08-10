package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFilterCreate,
		Update: resourceAciFilterUpdate,
		Read:   resourceAciFilterRead,
		Delete: resourceAciFilterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFilterImport,
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vz_rs_filt_graph_att": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vz_rs_fwd_r_flt_p_att": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vz_rs_rev_r_flt_p_att": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteFilter(client *client.Client, dn string) (*models.Filter, error) {
	vzFilterCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzFilter := models.FilterFromContainer(vzFilterCont)

	if vzFilter.DistinguishedName == "" {
		return nil, fmt.Errorf("Filter %s not found", vzFilter.DistinguishedName)
	}

	return vzFilter, nil
}

func setFilterAttributes(vzFilter *models.Filter, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(vzFilter.DistinguishedName)
	d.Set("description", vzFilter.Description)
	// d.Set("tenant_dn", GetParentDn(vzFilter.DistinguishedName))
	if dn != vzFilter.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vzFilterMap, _ := vzFilter.ToMap()

	d.Set("name", vzFilterMap["name"])

	d.Set("annotation", vzFilterMap["annotation"])
	d.Set("name_alias", vzFilterMap["nameAlias"])
	return d
}

func resourceAciFilterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzFilter, err := getRemoteFilter(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFilterAttributes(vzFilter, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFilterCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Filter: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzFilterAttr := models.FilterAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzFilterAttr.Annotation = Annotation.(string)
	} else {
		vzFilterAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzFilterAttr.NameAlias = NameAlias.(string)
	}
	vzFilter := models.NewFilter(fmt.Sprintf("flt-%s", name), TenantDn, desc, vzFilterAttr)

	err := aciClient.Save(vzFilter)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTovzRsFiltGraphAtt, ok := d.GetOk("relation_vz_rs_filt_graph_att"); ok {
		relationParam := relationTovzRsFiltGraphAtt.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvzRsFiltGraphAttFromFilter(vzFilter.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_filt_graph_att")
		d.Partial(false)

	}
	if relationTovzRsFwdRFltPAtt, ok := d.GetOk("relation_vz_rs_fwd_r_flt_p_att"); ok {
		relationParam := relationTovzRsFwdRFltPAtt.(string)
		err = aciClient.CreateRelationvzRsFwdRFltPAttFromFilter(vzFilter.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_fwd_r_flt_p_att")
		d.Partial(false)

	}
	if relationTovzRsRevRFltPAtt, ok := d.GetOk("relation_vz_rs_rev_r_flt_p_att"); ok {
		relationParam := relationTovzRsRevRFltPAtt.(string)
		err = aciClient.CreateRelationvzRsRevRFltPAttFromFilter(vzFilter.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_rev_r_flt_p_att")
		d.Partial(false)

	}

	d.SetId(vzFilter.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFilterRead(d, m)
}

func resourceAciFilterUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Filter: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzFilterAttr := models.FilterAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzFilterAttr.Annotation = Annotation.(string)
	} else {
		vzFilterAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzFilterAttr.NameAlias = NameAlias.(string)
	}
	vzFilter := models.NewFilter(fmt.Sprintf("flt-%s", name), TenantDn, desc, vzFilterAttr)

	vzFilter.Status = "modified"

	err := aciClient.Save(vzFilter)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_vz_rs_filt_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_filt_graph_att")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvzRsFiltGraphAttFromFilter(vzFilter.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_filt_graph_att")
		d.Partial(false)

	}
	if d.HasChange("relation_vz_rs_fwd_r_flt_p_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_fwd_r_flt_p_att")
		err = aciClient.CreateRelationvzRsFwdRFltPAttFromFilter(vzFilter.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_fwd_r_flt_p_att")
		d.Partial(false)

	}
	if d.HasChange("relation_vz_rs_rev_r_flt_p_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_rev_r_flt_p_att")
		err = aciClient.CreateRelationvzRsRevRFltPAttFromFilter(vzFilter.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vz_rs_rev_r_flt_p_att")
		d.Partial(false)

	}

	d.SetId(vzFilter.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFilterRead(d, m)

}

func resourceAciFilterRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzFilter, err := getRemoteFilter(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFilterAttributes(vzFilter, d)

	vzRsFiltGraphAttData, err := aciClient.ReadRelationvzRsFiltGraphAttFromFilter(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsFiltGraphAtt %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_filt_graph_att"); ok {
			tfName := GetMOName(d.Get("relation_vz_rs_filt_graph_att").(string))
			if tfName != vzRsFiltGraphAttData {
				d.Set("relation_vz_rs_filt_graph_att", "")
			}
		}
	}

	vzRsFwdRFltPAttData, err := aciClient.ReadRelationvzRsFwdRFltPAttFromFilter(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsFwdRFltPAtt %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_fwd_r_flt_p_att"); ok {
			tfName := d.Get("relation_vz_rs_fwd_r_flt_p_att").(string)
			if tfName != vzRsFwdRFltPAttData {
				d.Set("relation_vz_rs_fwd_r_flt_p_att", "")
			}
		}
	}

	vzRsRevRFltPAttData, err := aciClient.ReadRelationvzRsRevRFltPAttFromFilter(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsRevRFltPAtt %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_rev_r_flt_p_att"); ok {
			tfName := d.Get("relation_vz_rs_rev_r_flt_p_att").(string)
			if tfName != vzRsRevRFltPAttData {
				d.Set("relation_vz_rs_rev_r_flt_p_att", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFilterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzFilter")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

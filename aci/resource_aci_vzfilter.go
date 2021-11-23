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

func resourceAciFilter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFilterCreate,
		UpdateContext: resourceAciFilterUpdate,
		ReadContext:   resourceAciFilterRead,
		DeleteContext: resourceAciFilterDelete,

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

func setFilterAttributes(vzFilter *models.Filter, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzFilter.DistinguishedName)
	d.Set("description", vzFilter.Description)
	if dn != vzFilter.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vzFilterMap, err := vzFilter.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/flt-%s", vzFilterMap["name"])))

	d.Set("name", vzFilterMap["name"])

	d.Set("annotation", vzFilterMap["annotation"])
	d.Set("name_alias", vzFilterMap["nameAlias"])
	return d, nil
}

func resourceAciFilterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzFilter, err := getRemoteFilter(aciClient, dn)

	if err != nil {
		return nil, err
	}
	vzFilterMap, err := vzFilter.ToMap()
	if err != nil {
		return nil, err
	}
	name := vzFilterMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/flt-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setFilterAttributes(vzFilter, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFilterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTovzRsFiltGraphAtt, ok := d.GetOk("relation_vz_rs_filt_graph_att"); ok {
		relationParam := relationTovzRsFiltGraphAtt.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovzRsFwdRFltPAtt, ok := d.GetOk("relation_vz_rs_fwd_r_flt_p_att"); ok {
		relationParam := relationTovzRsFwdRFltPAtt.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTovzRsRevRFltPAtt, ok := d.GetOk("relation_vz_rs_rev_r_flt_p_att"); ok {
		relationParam := relationTovzRsRevRFltPAtt.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovzRsFiltGraphAtt, ok := d.GetOk("relation_vz_rs_filt_graph_att"); ok {
		relationParam := relationTovzRsFiltGraphAtt.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvzRsFiltGraphAttFromFilter(vzFilter.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTovzRsFwdRFltPAtt, ok := d.GetOk("relation_vz_rs_fwd_r_flt_p_att"); ok {
		relationParam := relationTovzRsFwdRFltPAtt.(string)
		err = aciClient.CreateRelationvzRsFwdRFltPAttFromFilter(vzFilter.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTovzRsRevRFltPAtt, ok := d.GetOk("relation_vz_rs_rev_r_flt_p_att"); ok {
		relationParam := relationTovzRsRevRFltPAtt.(string)
		err = aciClient.CreateRelationvzRsRevRFltPAttFromFilter(vzFilter.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vzFilter.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFilterRead(ctx, d, m)
}

func resourceAciFilterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vz_rs_filt_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_filt_graph_att")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vz_rs_fwd_r_flt_p_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_fwd_r_flt_p_att")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_vz_rs_rev_r_flt_p_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_rev_r_flt_p_att")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vz_rs_filt_graph_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_filt_graph_att")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationvzRsFiltGraphAttFromFilter(vzFilter.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vz_rs_fwd_r_flt_p_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_fwd_r_flt_p_att")
		err = aciClient.CreateRelationvzRsFwdRFltPAttFromFilter(vzFilter.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vz_rs_rev_r_flt_p_att") {
		_, newRelParam := d.GetChange("relation_vz_rs_rev_r_flt_p_att")
		err = aciClient.CreateRelationvzRsRevRFltPAttFromFilter(vzFilter.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vzFilter.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFilterRead(ctx, d, m)

}

func resourceAciFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzFilter, err := getRemoteFilter(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setFilterAttributes(vzFilter, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	vzRsFiltGraphAttData, err := aciClient.ReadRelationvzRsFiltGraphAttFromFilter(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsFiltGraphAtt %v", err)
		d.Set("relation_vz_rs_filt_graph_att", "")

	} else {
		d.Set("relation_vz_rs_filt_graph_att", vzRsFiltGraphAttData.(string))
	}

	vzRsFwdRFltPAttData, err := aciClient.ReadRelationvzRsFwdRFltPAttFromFilter(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsFwdRFltPAtt %v", err)
		d.Set("relation_vz_rs_fwd_r_flt_p_att", "")

	} else {
		d.Set("relation_vz_rs_fwd_r_flt_p_att", vzRsFwdRFltPAttData.(string))
	}

	vzRsRevRFltPAttData, err := aciClient.ReadRelationvzRsRevRFltPAttFromFilter(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsRevRFltPAtt %v", err)
		d.Set("relation_vz_rs_rev_r_flt_p_att", "")

	} else {
		d.Set("relation_vz_rs_rev_r_flt_p_att", vzRsRevRFltPAttData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzFilter")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

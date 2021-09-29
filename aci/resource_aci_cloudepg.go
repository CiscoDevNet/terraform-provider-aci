package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudEPg() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudEPgCreate,
		UpdateContext: resourceAciCloudEPgUpdate,
		ReadContext:   resourceAciCloudEPgRead,
		DeleteContext: resourceAciCloudEPgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudEPgImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_applicationcontainer_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"All",
					"AtleastOne",
					"AtmostOne",
					"None",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"exclude",
					"include",
				}, false),
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"level1",
					"level2",
					"level3",
					"level4",
					"level5",
					"level6",
					"unspecified",
				}, false),
			},

			"relation_fv_rs_sec_inherited": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_cust_qos_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_cloud_rs_cloud_epg_ctx": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_prot_by": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_intra_epg": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteCloudEPg(client *client.Client, dn string) (*models.CloudEPg, error) {
	cloudEPgCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudEPg := models.CloudEPgFromContainer(cloudEPgCont)

	if cloudEPg.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudEPg %s not found", cloudEPg.DistinguishedName)
	}

	return cloudEPg, nil
}

func setCloudEPgAttributes(cloudEPg *models.CloudEPg, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudEPg.DistinguishedName)
	d.Set("description", cloudEPg.Description)

	if dn != cloudEPg.DistinguishedName {
		d.Set("cloud_applicationcontainer_dn", "")
	}
	cloudEPgMap, err := cloudEPg.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("cloud_applicationcontainer_dn", GetParentDn(dn, fmt.Sprintf("/cloudepg-%s", cloudEPgMap["name"])))
	d.Set("name", cloudEPgMap["name"])

	d.Set("annotation", cloudEPgMap["annotation"])
	d.Set("exception_tag", cloudEPgMap["exceptionTag"])
	d.Set("flood_on_encap", cloudEPgMap["floodOnEncap"])
	d.Set("match_t", cloudEPgMap["matchT"])
	d.Set("name_alias", cloudEPgMap["nameAlias"])
	d.Set("pref_gr_memb", cloudEPgMap["prefGrMemb"])
	d.Set("prio", cloudEPgMap["prio"])
	return d, nil
}

func resourceAciCloudEPgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudEPg, err := getRemoteCloudEPg(aciClient, dn)

	if err != nil {
		return nil, err
	}
	cloudEPgMap, err := cloudEPg.ToMap()
	if err != nil {
		return nil, err
	}
	name := cloudEPgMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/cloudepg-%s", name))
	d.Set("cloud_applicationcontainer_dn", pDN)
	schemaFilled, err := setCloudEPgAttributes(cloudEPg, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudEPgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudEPg: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudApplicationcontainerDn := d.Get("cloud_applicationcontainer_dn").(string)

	cloudEPgAttr := models.CloudEPgAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudEPgAttr.Annotation = Annotation.(string)
	} else {
		cloudEPgAttr.Annotation = "{}"
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		cloudEPgAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		cloudEPgAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		cloudEPgAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudEPgAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		cloudEPgAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		cloudEPgAttr.Prio = Prio.(string)
	}
	cloudEPg := models.NewCloudEPg(fmt.Sprintf("cloudepg-%s", name), CloudApplicationcontainerDn, desc, cloudEPgAttr)

	err := aciClient.Save(cloudEPg)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		checkDns = append(checkDns, relationParam)

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTocloudRsCloudEPgCtx, ok := d.GetOk("relation_cloud_rs_cloud_epg_ctx"); ok {
		relationParam := relationTocloudRsCloudEPgCtx.(string)
		checkDns = append(checkDns, relationParam)

	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromCloudEPg(cloudEPg.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProvFromCloudEPg(cloudEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsIfFromCloudEPg(cloudEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsCustQosPolFromCloudEPg(cloudEPg.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsFromCloudEPg(cloudEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTocloudRsCloudEPgCtx, ok := d.GetOk("relation_cloud_rs_cloud_epg_ctx"); ok {
		relationParam := relationTocloudRsCloudEPgCtx.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationcloudRsCloudEPgCtxFromCloudEPg(cloudEPg.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProtByFromCloudEPg(cloudEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsIntraEpgFromCloudEPg(cloudEPg.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}

	d.SetId(cloudEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudEPgRead(ctx, d, m)
}

func resourceAciCloudEPgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudEPg: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudApplicationcontainerDn := d.Get("cloud_applicationcontainer_dn").(string)

	cloudEPgAttr := models.CloudEPgAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudEPgAttr.Annotation = Annotation.(string)
	} else {
		cloudEPgAttr.Annotation = "{}"
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		cloudEPgAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		cloudEPgAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		cloudEPgAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudEPgAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		cloudEPgAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		cloudEPgAttr.Prio = Prio.(string)
	}
	cloudEPg := models.NewCloudEPg(fmt.Sprintf("cloudepg-%s", name), CloudApplicationcontainerDn, desc, cloudEPgAttr)

	cloudEPg.Status = "modified"

	err := aciClient.Save(cloudEPg)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_prov") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fv_rs_cons") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}

	}
	if d.HasChange("relation_cloud_rs_cloud_epg_ctx") {
		_, newRelParam := d.GetChange("relation_cloud_rs_cloud_epg_ctx")
		checkDns = append(checkDns, newRelParam.(string))

	}
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_intra_epg") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInheritedFromCloudEPg(cloudEPg.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromCloudEPg(cloudEPg.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_prov") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsProvFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProvFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsConsIfFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsIfFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsCustQosPolFromCloudEPg(cloudEPg.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_cons") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsConsFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_cloud_rs_cloud_epg_ctx") {
		_, newRelParam := d.GetChange("relation_cloud_rs_cloud_epg_ctx")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationcloudRsCloudEPgCtxFromCloudEPg(cloudEPg.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsProtByFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProtByFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_intra_epg") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsIntraEpgFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationfvRsIntraEpgFromCloudEPg(cloudEPg.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(cloudEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudEPgRead(ctx, d, m)

}

func resourceAciCloudEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudEPg, err := getRemoteCloudEPg(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setCloudEPgAttributes(cloudEPg, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInheritedFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)
		d.Set("relation_fv_rs_sec_inherited", make([]interface{}, 0, 1))

	} else {
		d.Set("relation_fv_rs_sec_inherited", fvRsSecInheritedData)
	}

	fvRsProvData, err := aciClient.ReadRelationfvRsProvFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProv %v", err)
		d.Set("relation_fv_rs_prov", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_fv_rs_prov"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_prov").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsProvDataList := toStringList(fvRsProvData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsProvDataList)

			if !reflect.DeepEqual(tfList, fvRsProvDataList) {
				d.Set("relation_fv_rs_prov", make([]string, 0, 1))
			}
		}
	}

	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIfFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)
		d.Set("relation_fv_rs_cons_if", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_cons_if").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsConsIfDataList := toStringList(fvRsConsIfData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsConsIfDataList)

			if !reflect.DeepEqual(tfList, fvRsConsIfDataList) {
				d.Set("relation_fv_rs_cons_if", make([]string, 0, 1))
			}
		}
	}

	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPolFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)
		d.Set("relation_fv_rs_cust_qos_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_cust_qos_pol").(string))
			if tfName != fvRsCustQosPolData {
				d.Set("relation_fv_rs_cust_qos_pol", "")
			}
		}
	}

	fvRsConsData, err := aciClient.ReadRelationfvRsConsFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCons %v", err)
		d.Set("relation_fv_rs_cons", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_fv_rs_cons"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_cons").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsConsDataList := toStringList(fvRsConsData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsConsDataList)

			if !reflect.DeepEqual(tfList, fvRsConsDataList) {
				d.Set("relation_fv_rs_cons", make([]string, 0, 1))
			}
		}
	}

	cloudRsCloudEPgCtxData, err := aciClient.ReadRelationcloudRsCloudEPgCtxFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsCloudEPgCtx %v", err)
		d.Set("relation_cloud_rs_cloud_epg_ctx", "")

	} else {
		if _, ok := d.GetOk("relation_cloud_rs_cloud_epg_ctx"); ok {
			tfName := GetMOName(d.Get("relation_cloud_rs_cloud_epg_ctx").(string))
			if tfName != cloudRsCloudEPgCtxData {
				d.Set("relation_cloud_rs_cloud_epg_ctx", "")
			}
		}
	}

	fvRsProtByData, err := aciClient.ReadRelationfvRsProtByFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)
		d.Set("relation_fv_rs_prot_by", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_prot_by").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsProtByDataList := toStringList(fvRsProtByData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsProtByDataList)

			if !reflect.DeepEqual(tfList, fvRsProtByDataList) {
				d.Set("relation_fv_rs_prot_by", make([]string, 0, 1))
			}
		}
	}

	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpgFromCloudEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)
		d.Set("relation_fv_rs_intra_epg", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_intra_epg").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			fvRsIntraEpgDataList := toStringList(fvRsIntraEpgData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(fvRsIntraEpgDataList)

			if !reflect.DeepEqual(tfList, fvRsIntraEpgDataList) {
				d.Set("relation_fv_rs_intra_epg", make([]string, 0, 1))
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudEPgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudEPg")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

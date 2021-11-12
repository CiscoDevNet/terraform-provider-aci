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

func resourceAciNodeManagementEPg() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciNodeManagementEPgCreate,
		UpdateContext: resourceAciNodeManagementEPgUpdate,
		ReadContext:   resourceAciNodeManagementEPgRead,
		DeleteContext: resourceAciNodeManagementEPgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciNodeManagementEPgImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			// Common Attributes

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in_band",
					"out_of_band",
				}, false),
			},

			"management_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "uni/tn-mgmt/mgmtp-default",
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

			//Attributes of mgmtInB

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("type").(string) == "in_band" {
						return false
					}
					return true
				},
			},

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("type").(string) == "in_band" {
						return false
					}
					return true
				},
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("type").(string) == "in_band" {
						return false
					}
					return true
				},
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("type").(string) == "in_band" {
						return false
					}
					return true
				},
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"exclude",
					"include",
				}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("type").(string) == "in_band" {
						return false
					}
					return true
				},
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_mgmt_rs_mgmt_bd": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_prot_by": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_mgmt_rs_in_b_st_node": &schema.Schema{
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

			// Attributes of mgmtOob

			"relation_mgmt_rs_oo_b_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_mgmt_rs_oo_b_st_node": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_mgmt_rs_oo_b_ctx": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceAciNodeManagementEPgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	if d.Get("type").(string) == "in_band" {
		return inBandManagementEPgImport(d, m)
	}
	return outOfBandManagementEPgImport(d, m)
}

func resourceAciNodeManagementEPgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("type").(string) == "in_band" {
		log.Printf("[DEBUG] InBandManagementEPg: Beginning Creation")
		return inBandManagementEPgCreate(ctx, d, m)
	}
	log.Printf("[DEBUG] OutOfBandManagementEPg: Beginning Creation")
	return outOfBandManagementEPgCreate(ctx, d, m)
}

func resourceAciNodeManagementEPgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("type").(string) == "in_band" {
		log.Printf("[DEBUG] InBandManagementEPg: Beginning Update")
		return inBandManagementEPgUpdate(ctx, d, m)
	}
	log.Printf("[DEBUG] OutOfBandManagementEPg: Beginning Update")
	return outOfBandManagementEPgUpdate(ctx, d, m)
}

func resourceAciNodeManagementEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	if d.Get("type").(string) == "in_band" {
		return inBandManagementEPgRead(ctx, d, m)
	}
	return outOfBandManagementEPgRead(ctx, d, m)
}

func resourceAciNodeManagementEPgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	if d.Get("type").(string) == "in_band" {
		return inBandManagementEPgDelete(ctx, d, m)
	}
	return outOfBandManagementEPgDelete(ctx, d, m)
}

func getRemoteInBandManagementEPg(client *client.Client, dn string) (*models.InBandManagementEPg, error) {
	mgmtInBCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtInB := models.InBandManagementEPgFromContainer(mgmtInBCont)

	if mgmtInB.DistinguishedName == "" {
		return nil, fmt.Errorf("InBandManagementEPg %s not found", mgmtInB.DistinguishedName)
	}

	return mgmtInB, nil
}

func getRemoteOutOfBandManagementEPg(client *client.Client, dn string) (*models.OutOfBandManagementEPg, error) {
	mgmtOoBCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtOoB := models.OutOfBandManagementEPgFromContainer(mgmtOoBCont)

	if mgmtOoB.DistinguishedName == "" {
		return nil, fmt.Errorf("OutOfBandManagementEPg %s not found", mgmtOoB.DistinguishedName)
	}

	return mgmtOoB, nil
}

func setInBandManagementEPgAttributes(mgmtInB *models.InBandManagementEPg, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mgmtInB.DistinguishedName)
	d.Set("description", mgmtInB.Description)
	dn := d.Id()
	if dn != mgmtInB.DistinguishedName {
		d.Set("management_profile_dn", "")
	}
	mgmtInBMap, err := mgmtInB.ToMap()

	if err != nil {
		return d, err
	}

	d.Set("name", mgmtInBMap["name"])
	d.Set("management_profile_dn", GetParentDn(mgmtInB.DistinguishedName, fmt.Sprintf("/inb-%s", mgmtInBMap["name"])))
	d.Set("annotation", mgmtInBMap["annotation"])
	d.Set("encap", mgmtInBMap["encap"])
	d.Set("exception_tag", mgmtInBMap["exceptionTag"])
	d.Set("flood_on_encap", mgmtInBMap["floodOnEncap"])
	d.Set("match_t", mgmtInBMap["matchT"])
	d.Set("name_alias", mgmtInBMap["nameAlias"])
	d.Set("pref_gr_memb", mgmtInBMap["prefGrMemb"])
	d.Set("prio", mgmtInBMap["prio"])
	return d, nil
}

func setOutOfBandManagementEPgAttributes(mgmtOoB *models.OutOfBandManagementEPg, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mgmtOoB.DistinguishedName)
	d.Set("description", mgmtOoB.Description)
	dn := d.Id()
	if dn != mgmtOoB.DistinguishedName {
		d.Set("management_profile_dn", "")
	}
	mgmtOoBMap, err := mgmtOoB.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("name", mgmtOoBMap["name"])
	d.Set("management_profile_dn", GetParentDn(mgmtOoB.DistinguishedName, fmt.Sprintf("/oob-%s", mgmtOoBMap["name"])))

	d.Set("annotation", mgmtOoBMap["annotation"])
	d.Set("name_alias", mgmtOoBMap["nameAlias"])
	d.Set("prio", mgmtOoBMap["prio"])
	return d, nil
}

func inBandManagementEPgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	mgmtInB, err := getRemoteInBandManagementEPg(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setInBandManagementEPgAttributes(mgmtInB, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func outOfBandManagementEPgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	mgmtOoB, err := getRemoteOutOfBandManagementEPg(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setOutOfBandManagementEPgAttributes(mgmtOoB, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func inBandManagementEPgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ManagementProfileDn := d.Get("management_profile_dn").(string)

	mgmtInBAttr := models.InBandManagementEPgAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtInBAttr.Annotation = Annotation.(string)
	} else {
		mgmtInBAttr.Annotation = "{}"
	}
	if Encap, ok := d.GetOk("encap"); ok {
		mgmtInBAttr.Encap = Encap.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		mgmtInBAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		mgmtInBAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		mgmtInBAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		mgmtInBAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		mgmtInBAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		mgmtInBAttr.Prio = Prio.(string)
	}
	mgmtInB := models.NewInBandManagementEPg(fmt.Sprintf("inb-%s", name), ManagementProfileDn, desc, mgmtInBAttr)

	err := aciClient.Save(mgmtInB)
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
	if relationTomgmtRsMgmtBD, ok := d.GetOk("relation_mgmt_rs_mgmt_bd"); ok {
		relationParam := relationTomgmtRsMgmtBD.(string)
		checkDns = append(checkDns, relationParam)

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTomgmtRsInBStNode, ok := d.GetOk("relation_mgmt_rs_in_b_st_node"); ok {
		relationParamList := toStringList(relationTomgmtRsInBStNode.(*schema.Set).List())
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

			err = aciClient.CreateRelationfvRsSecInheritedFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParam = GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProvFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParam = GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsIfFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := GetMOName(relationTofvRsCustQosPol.(string))
		err = aciClient.CreateRelationfvRsCustQosPolFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTomgmtRsMgmtBD, ok := d.GetOk("relation_mgmt_rs_mgmt_bd"); ok {
		relationParam := GetMOName(relationTomgmtRsMgmtBD.(string))
		err = aciClient.CreateRelationmgmtRsMgmtBDFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParam = GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsConsFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParam = GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsProtByFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTomgmtRsInBStNode, ok := d.GetOk("relation_mgmt_rs_in_b_st_node"); ok {
		relationParamList := toStringList(relationTomgmtRsInBStNode.(*schema.Set).List())
		for _, relationParam := range relationParamList {

			err = aciClient.CreateRelationmgmtRsInBStNodeFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParam = GetMOName(relationParam)
			err = aciClient.CreateRelationfvRsIntraEpgFromInBandManagementEPg(mgmtInB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(mgmtInB.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciNodeManagementEPgRead(ctx, d, m)
}

func outOfBandManagementEPgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ManagementProfileDn := d.Get("management_profile_dn").(string)

	mgmtOoBAttr := models.OutOfBandManagementEPgAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtOoBAttr.Annotation = Annotation.(string)
	} else {
		mgmtOoBAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		mgmtOoBAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		mgmtOoBAttr.Prio = Prio.(string)
	}
	mgmtOoB := models.NewOutOfBandManagementEPg(fmt.Sprintf("oob-%s", name), ManagementProfileDn, desc, mgmtOoBAttr)

	err := aciClient.Save(mgmtOoB)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTomgmtRsOoBProv, ok := d.GetOk("relation_mgmt_rs_oo_b_prov"); ok {
		relationParamList := toStringList(relationTomgmtRsOoBProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTomgmtRsOoBStNode, ok := d.GetOk("relation_mgmt_rs_oo_b_st_node"); ok {
		relationParamList := toStringList(relationTomgmtRsOoBStNode.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTomgmtRsOoBCtx, ok := d.GetOk("relation_mgmt_rs_oo_b_ctx"); ok {
		relationParam := relationTomgmtRsOoBCtx.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTomgmtRsOoBProv, ok := d.GetOk("relation_mgmt_rs_oo_b_prov"); ok {
		relationParamList := toStringList(relationTomgmtRsOoBProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParam = GetMOName(relationParam)
			err = aciClient.CreateRelationmgmtRsOoBProvFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTomgmtRsOoBStNode, ok := d.GetOk("relation_mgmt_rs_oo_b_st_node"); ok {
		relationParamList := toStringList(relationTomgmtRsOoBStNode.(*schema.Set).List())
		for _, relationParam := range relationParamList {

			err = aciClient.CreateRelationmgmtRsOoBStNodeFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTomgmtRsOoBCtx, ok := d.GetOk("relation_mgmt_rs_oo_b_ctx"); ok {
		relationParam := GetMOName(relationTomgmtRsOoBCtx.(string))
		err = aciClient.CreateRelationmgmtRsOoBCtxFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(mgmtOoB.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciNodeManagementEPgRead(ctx, d, m)
}

func inBandManagementEPgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InBandManagementEPg: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ManagementProfileDn := d.Get("management_profile_dn").(string)

	mgmtInBAttr := models.InBandManagementEPgAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtInBAttr.Annotation = Annotation.(string)
	} else {
		mgmtInBAttr.Annotation = "{}"
	}
	if Encap, ok := d.GetOk("encap"); ok {
		mgmtInBAttr.Encap = Encap.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		mgmtInBAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		mgmtInBAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		mgmtInBAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		mgmtInBAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		mgmtInBAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		mgmtInBAttr.Prio = Prio.(string)
	}
	mgmtInB := models.NewInBandManagementEPg(fmt.Sprintf("inb-%s", name), ManagementProfileDn, desc, mgmtInBAttr)

	mgmtInB.Status = "modified"

	err := aciClient.Save(mgmtInB)

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
	if d.HasChange("relation_mgmt_rs_mgmt_bd") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_mgmt_bd")
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
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}

	}
	if d.HasChange("relation_mgmt_rs_in_b_st_node") {
		oldRel, newRel := d.GetChange("relation_mgmt_rs_in_b_st_node")
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

			err = aciClient.DeleteRelationfvRsSecInheritedFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {

			err = aciClient.CreateRelationfvRsSecInheritedFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
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
			relDn = GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsProvFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDn = GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProvFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
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
			relDn = GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsConsIfFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDn = GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsIfFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		err = aciClient.CreateRelationfvRsCustQosPolFromInBandManagementEPg(mgmtInB.DistinguishedName, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_mgmt_rs_mgmt_bd") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_mgmt_bd")
		err = aciClient.CreateRelationmgmtRsMgmtBDFromInBandManagementEPg(mgmtInB.DistinguishedName, GetMOName(newRelParam.(string)))
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
			relDn = GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsConsFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDn = GetMOName(relDn)
			err = aciClient.CreateRelationfvRsConsFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDn = GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsProtByFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDn = GetMOName(relDn)
			err = aciClient.CreateRelationfvRsProtByFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_mgmt_rs_in_b_st_node") {
		oldRel, newRel := d.GetChange("relation_mgmt_rs_in_b_st_node")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {

			err = aciClient.DeleteRelationmgmtRsInBStNodeFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {

			err = aciClient.CreateRelationmgmtRsInBStNodeFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
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
			relDn = GetMOName(relDn)
			err = aciClient.DeleteRelationfvRsIntraEpgFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDn = GetMOName(relDn)
			err = aciClient.CreateRelationfvRsIntraEpgFromInBandManagementEPg(mgmtInB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(mgmtInB.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciNodeManagementEPgRead(ctx, d, m)

}

func outOfBandManagementEPgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OutOfBandManagementEPg: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ManagementProfileDn := d.Get("management_profile_dn").(string)

	mgmtOoBAttr := models.OutOfBandManagementEPgAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtOoBAttr.Annotation = Annotation.(string)
	} else {
		mgmtOoBAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		mgmtOoBAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		mgmtOoBAttr.Prio = Prio.(string)
	}
	mgmtOoB := models.NewOutOfBandManagementEPg(fmt.Sprintf("oob-%s", name), ManagementProfileDn, desc, mgmtOoBAttr)

	mgmtOoB.Status = "modified"

	err := aciClient.Save(mgmtOoB)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_mgmt_rs_oo_b_prov") {
		oldRel, newRel := d.GetChange("relation_mgmt_rs_oo_b_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}

	}
	if d.HasChange("relation_mgmt_rs_oo_b_st_node") {
		oldRel, newRel := d.GetChange("relation_mgmt_rs_oo_b_st_node")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}

	}
	if d.HasChange("relation_mgmt_rs_oo_b_ctx") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_oo_b_ctx")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_mgmt_rs_oo_b_prov") {
		oldRel, newRel := d.GetChange("relation_mgmt_rs_oo_b_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDn = GetMOName(relDn)
			err = aciClient.DeleteRelationmgmtRsOoBProvFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDn = GetMOName(relDn)
			err = aciClient.CreateRelationmgmtRsOoBProvFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_mgmt_rs_oo_b_st_node") {
		oldRel, newRel := d.GetChange("relation_mgmt_rs_oo_b_st_node")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {

			err = aciClient.DeleteRelationmgmtRsOoBStNodeFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {

			err = aciClient.CreateRelationmgmtRsOoBStNodeFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_mgmt_rs_oo_b_ctx") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_oo_b_ctx")
		err = aciClient.CreateRelationmgmtRsOoBCtxFromOutOfBandManagementEPg(mgmtOoB.DistinguishedName, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(mgmtOoB.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciNodeManagementEPgRead(ctx, d, m)

}

func inBandManagementEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	mgmtInB, err := getRemoteInBandManagementEPg(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setInBandManagementEPgAttributes(mgmtInB, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInheritedFromInBandManagementEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)
		d.Set("relation_fv_rs_sec_inherited", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
			relationParamList := toStringList(d.Get("relation_fv_rs_sec_inherited").(*schema.Set).List())
			fvRsSecInheritedDataList := toStringList(fvRsSecInheritedData.(*schema.Set).List())
			sort.Strings(relationParamList)
			sort.Strings(fvRsSecInheritedDataList)

			if !reflect.DeepEqual(relationParamList, fvRsSecInheritedDataList) {
				d.Set("relation_fv_rs_sec_inherited", make([]string, 0, 1))
			}
		}
	}
	fvRsProvData, err := aciClient.ReadRelationfvRsProvFromInBandManagementEPg(dn)
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
	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIfFromInBandManagementEPg(dn)
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
	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPolFromInBandManagementEPg(dn)
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

	mgmtRsMgmtBDData, err := aciClient.ReadRelationmgmtRsMgmtBDFromInBandManagementEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsMgmtBD %v", err)
		d.Set("relation_mgmt_rs_mgmt_bd", "")

	} else {
		if _, ok := d.GetOk("relation_mgmt_rs_mgmt_bd"); ok {
			tfName := GetMOName(d.Get("relation_mgmt_rs_mgmt_bd").(string))
			if tfName != mgmtRsMgmtBDData {
				d.Set("relation_mgmt_rs_mgmt_bd", "")
			}
		}

	}

	fvRsConsData, err := aciClient.ReadRelationfvRsConsFromInBandManagementEPg(dn)
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
	fvRsProtByData, err := aciClient.ReadRelationfvRsProtByFromInBandManagementEPg(dn)
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
	mgmtRsInBStNodeData, err := aciClient.ReadRelationmgmtRsInBStNodeFromInBandManagementEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsInBStNode %v", err)
		d.Set("relation_mgmt_rs_in_b_st_node", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_mgmt_rs_in_b_st_node"); ok {
			relationParamList := toStringList(d.Get("relation_mgmt_rs_in_b_st_node").(*schema.Set).List())
			mgmtRsInBStNodeDataList := toStringList(mgmtRsInBStNodeData.(*schema.Set).List())
			sort.Strings(relationParamList)
			sort.Strings(mgmtRsInBStNodeDataList)

			if !reflect.DeepEqual(relationParamList, mgmtRsInBStNodeDataList) {
				d.Set("relation_mgmt_rs_in_b_st_node", make([]string, 0, 1))
			}
		}
	}
	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpgFromInBandManagementEPg(dn)
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

func outOfBandManagementEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	mgmtOoB, err := getRemoteOutOfBandManagementEPg(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setOutOfBandManagementEPgAttributes(mgmtOoB, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	mgmtRsOoBProvData, err := aciClient.ReadRelationmgmtRsOoBProvFromOutOfBandManagementEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsOoBProv %v", err)
		d.Set("relation_mgmt_rs_oo_b_prov", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_mgmt_rs_oo_b_prov"); ok {
			relationParamList := toStringList(d.Get("relation_mgmt_rs_oo_b_prov").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			mgmtRsOoBProvDataList := toStringList(mgmtRsOoBProvData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(mgmtRsOoBProvDataList)

			if !reflect.DeepEqual(tfList, mgmtRsOoBProvDataList) {
				d.Set("relation_mgmt_rs_oo_b_prov", make([]string, 0, 1))
			}
		}
	}
	mgmtRsOoBStNodeData, err := aciClient.ReadRelationmgmtRsOoBStNodeFromOutOfBandManagementEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsOoBStNode %v", err)
		d.Set("relation_mgmt_rs_oo_b_st_node", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_mgmt_rs_oo_b_st_node"); ok {
			relationParamList := toStringList(d.Get("relation_mgmt_rs_oo_b_st_node").(*schema.Set).List())
			mgmtRsOoBStNodeDataList := toStringList(mgmtRsOoBStNodeData.(*schema.Set).List())
			sort.Strings(relationParamList)
			sort.Strings(mgmtRsOoBStNodeDataList)

			if !reflect.DeepEqual(relationParamList, mgmtRsOoBStNodeDataList) {
				d.Set("relation_mgmt_rs_oo_b_st_node", make([]string, 0, 1))
			}
		}
	}
	mgmtRsOoBCtxData, err := aciClient.ReadRelationmgmtRsOoBCtxFromOutOfBandManagementEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsOoBCtx %v", err)
		d.Set("relation_mgmt_rs_oo_b_ctx", "")

	} else {
		if _, ok := d.GetOk("relation_mgmt_rs_oo_b_ctx"); ok {
			tfName := GetMOName(d.Get("relation_mgmt_rs_oo_b_ctx").(string))
			if tfName != mgmtRsOoBCtxData {
				d.Set("relation_mgmt_rs_oo_b_ctx", "")
			}
		}

	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func inBandManagementEPgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "mgmtInB")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

func outOfBandManagementEPgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "mgmtOoB")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

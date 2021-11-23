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

func resourceAciAny() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAnyCreate,
		UpdateContext: resourceAciAnyUpdate,
		ReadContext:   resourceAciAnyRead,
		DeleteContext: resourceAciAnyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAnyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
					"enabled",
					"disabled",
				}, false),
			},

			"relation_vz_rs_any_to_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_vz_rs_any_to_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_vz_rs_any_to_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteAny(client *client.Client, dn string) (*models.Any, error) {
	vzAnyCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzAny := models.AnyFromContainer(vzAnyCont)

	if vzAny.DistinguishedName == "" {
		return nil, fmt.Errorf("Any %s not found", vzAny.DistinguishedName)
	}

	return vzAny, nil
}

func setAnyAttributes(vzAny *models.Any, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzAny.DistinguishedName)
	d.Set("description", vzAny.Description)
	if dn != vzAny.DistinguishedName {
		d.Set("vrf_dn", "")
	}
	vzAnyMap, err := vzAny.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("vrf_dn", GetParentDn(vzAny.DistinguishedName, fmt.Sprintf("/any")))
	d.Set("annotation", vzAnyMap["annotation"])
	d.Set("match_t", vzAnyMap["matchT"])
	d.Set("name_alias", vzAnyMap["nameAlias"])
	d.Set("pref_gr_memb", vzAnyMap["prefGrMemb"])
	return d, nil
}

func resourceAciAnyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzAny, err := getRemoteAny(aciClient, dn)

	if err != nil {
		return nil, err
	}
	pDN := GetParentDn(dn, fmt.Sprintf("/any"))
	d.Set("vrf_dn", pDN)
	schemaFilled, err := setAnyAttributes(vzAny, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAnyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Any: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	VRFDn := d.Get("vrf_dn").(string)

	vzAnyAttr := models.AnyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzAnyAttr.Annotation = Annotation.(string)
	} else {
		vzAnyAttr.Annotation = "{}"
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		vzAnyAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzAnyAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		vzAnyAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	vzAny := models.NewAny(fmt.Sprintf("any"), VRFDn, desc, vzAnyAttr)
	vzAny.Status = "modified"
	err := aciClient.Save(vzAny)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTovzRsAnyToCons, ok := d.GetOk("relation_vz_rs_any_to_cons"); ok {
		relationParamList := toStringList(relationTovzRsAnyToCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTovzRsAnyToConsIf, ok := d.GetOk("relation_vz_rs_any_to_cons_if"); ok {
		relationParamList := toStringList(relationTovzRsAnyToConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTovzRsAnyToProv, ok := d.GetOk("relation_vz_rs_any_to_prov"); ok {
		relationParamList := toStringList(relationTovzRsAnyToProv.(*schema.Set).List())
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

	if relationTovzRsAnyToCons, ok := d.GetOk("relation_vz_rs_any_to_cons"); ok {
		relationParamList := toStringList(relationTovzRsAnyToCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsAnyToConsFromAny(vzAny.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTovzRsAnyToConsIf, ok := d.GetOk("relation_vz_rs_any_to_cons_if"); ok {
		relationParamList := toStringList(relationTovzRsAnyToConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsAnyToConsIfFromAny(vzAny.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTovzRsAnyToProv, ok := d.GetOk("relation_vz_rs_any_to_prov"); ok {
		relationParamList := toStringList(relationTovzRsAnyToProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsAnyToProvFromAny(vzAny.DistinguishedName, relationParamName)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}

	d.SetId(vzAny.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAnyRead(ctx, d, m)
}

func resourceAciAnyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Any: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	VRFDn := d.Get("vrf_dn").(string)

	vzAnyAttr := models.AnyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzAnyAttr.Annotation = Annotation.(string)
	} else {
		vzAnyAttr.Annotation = "{}"
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		vzAnyAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzAnyAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		vzAnyAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	vzAny := models.NewAny(fmt.Sprintf("any"), VRFDn, desc, vzAnyAttr)

	vzAny.Status = "modified"

	err := aciClient.Save(vzAny)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vz_rs_any_to_cons") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_vz_rs_any_to_cons_if") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_vz_rs_any_to_prov") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_prov")
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

	if d.HasChange("relation_vz_rs_any_to_cons") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationvzRsAnyToConsFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsAnyToConsFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_vz_rs_any_to_cons_if") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationvzRsAnyToConsIfFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsAnyToConsIfFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_vz_rs_any_to_prov") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationvzRsAnyToProvFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsAnyToProvFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(vzAny.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAnyRead(ctx, d, m)

}

func resourceAciAnyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzAny, err := getRemoteAny(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setAnyAttributes(vzAny, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vzRsAnyToConsData, err := aciClient.ReadRelationvzRsAnyToConsFromAny(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsAnyToCons %v", err)
		d.Set("relation_vz_rs_any_to_cons", make([]string, 0, 1))

	} else {
		d.Set("relation_vz_rs_any_to_cons", toStringList(vzRsAnyToConsData.(*schema.Set).List()))
	}

	vzRsAnyToConsIfData, err := aciClient.ReadRelationvzRsAnyToConsIfFromAny(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsAnyToConsIf %v", err)
		d.Set("relation_vz_rs_any_to_cons_if", make([]string, 0, 1))

	} else {
		d.Set("relation_vz_rs_any_to_cons_if", toStringList(vzRsAnyToConsIfData.(*schema.Set).List()))
	}

	vzRsAnyToProvData, err := aciClient.ReadRelationvzRsAnyToProvFromAny(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsAnyToProv %v", err)
		d.Set("relation_vz_rs_any_to_prov", make([]string, 0, 1))
	} else {
		d.Set("relation_vz_rs_any_to_prov", toStringList(vzRsAnyToProvData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAnyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzAny")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

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

func resourceAciApplicationProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciApplicationProfileCreate,
		UpdateContext: resourceAciApplicationProfileUpdate,
		ReadContext:   resourceAciApplicationProfileRead,
		DeleteContext: resourceAciApplicationProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciApplicationProfileImport,
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

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"level3",
					"level2",
					"level1",
					"level4",
					"level5",
					"level6",
				}, false),
			},

			"relation_fv_rs_ap_mon_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteApplicationProfile(client *client.Client, dn string) (*models.ApplicationProfile, error) {
	fvApCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvAp := models.ApplicationProfileFromContainer(fvApCont)

	if fvAp.DistinguishedName == "" {
		return nil, fmt.Errorf("ApplicationProfile %s not found", fvAp.DistinguishedName)
	}

	return fvAp, nil
}

func setApplicationProfileAttributes(fvAp *models.ApplicationProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvAp.DistinguishedName)
	d.Set("description", fvAp.Description)
	if dn != fvAp.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	fvApMap, err := fvAp.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/ap-%s", fvApMap["name"])))

	d.Set("name", fvApMap["name"])

	d.Set("annotation", fvApMap["annotation"])
	d.Set("name_alias", fvApMap["nameAlias"])
	d.Set("prio", fvApMap["prio"])
	return d, nil
}

func resourceAciApplicationProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvAp, err := getRemoteApplicationProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}

	fvApMap, err := fvAp.ToMap()
	if err != nil {
		return nil, err
	}
	name := fvApMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/ap-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setApplicationProfileAttributes(fvAp, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciApplicationProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ApplicationProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvApAttr := models.ApplicationProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvApAttr.Annotation = Annotation.(string)
	} else {
		fvApAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvApAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		fvApAttr.Prio = Prio.(string)
	}
	fvAp := models.NewApplicationProfile(fmt.Sprintf("ap-%s", name), TenantDn, desc, fvApAttr)

	err := aciClient.Save(fvAp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofvRsApMonPol, ok := d.GetOk("relation_fv_rs_ap_mon_pol"); ok {
		relationParam := relationTofvRsApMonPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTofvRsApMonPol, ok := d.GetOk("relation_fv_rs_ap_mon_pol"); ok {
		relationParam := relationTofvRsApMonPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsApMonPolFromApplicationProfile(fvAp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fvAp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciApplicationProfileRead(ctx, d, m)
}

func resourceAciApplicationProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ApplicationProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvApAttr := models.ApplicationProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvApAttr.Annotation = Annotation.(string)
	} else {
		fvApAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvApAttr.NameAlias = NameAlias.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		fvApAttr.Prio = Prio.(string)
	}
	fvAp := models.NewApplicationProfile(fmt.Sprintf("ap-%s", name), TenantDn, desc, fvApAttr)

	fvAp.Status = "modified"

	err := aciClient.Save(fvAp)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_ap_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ap_mon_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_fv_rs_ap_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_ap_mon_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfvRsApMonPolFromApplicationProfile(fvAp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsApMonPolFromApplicationProfile(fvAp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fvAp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciApplicationProfileRead(ctx, d, m)

}

func resourceAciApplicationProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvAp, err := getRemoteApplicationProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setApplicationProfileAttributes(fvAp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsApMonPolData, err := aciClient.ReadRelationfvRsApMonPolFromApplicationProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsApMonPol %v", err)
		d.Set("relation_fv_rs_ap_mon_pol", "")

	} else {
		if _, ok := d.GetOk("relation_fv_rs_ap_mon_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_ap_mon_pol").(string))
			if tfName != fvRsApMonPolData {
				d.Set("relation_fv_rs_ap_mon_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciApplicationProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvAp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

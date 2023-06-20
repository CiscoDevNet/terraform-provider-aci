package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciExternalProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciExternalProfileCreate,
		UpdateContext: resourceAciExternalProfileUpdate,
		ReadContext:   resourceAciExternalProfileRead,
		DeleteContext: resourceAciExternalProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciExternalProfileImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l3outside_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled_af": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"ipv4-mcast",
						"ipv6-mcast",
					}, false),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteExternalProfile(client *client.Client, dn string) (*models.ExternalProfile, error) {
	pimExtPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	pimExtP := models.ExternalProfileFromContainer(pimExtPCont)
	if pimExtP.DistinguishedName == "" {
		return nil, fmt.Errorf("External Profile %s not found", dn)
	}
	return pimExtP, nil
}

func setExternalProfileAttributes(pimExtP *models.ExternalProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(pimExtP.DistinguishedName)
	d.Set("description", pimExtP.Description)
	pimExtPMap, err := pimExtP.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != pimExtP.DistinguishedName {
		d.Set("l3outside_dn", "")
	} else {
		d.Set("l3outside_dn", GetParentDn(pimExtP.DistinguishedName, fmt.Sprintf("/"+models.RnPimExtP)))
	}
	d.Set("annotation", pimExtPMap["annotation"])
	enabledAfGet := make([]string, 0, 1)
	for _, val := range strings.Split(pimExtPMap["enabledAf"], ",") {
		enabledAfGet = append(enabledAfGet, strings.Trim(val, " "))
	}
	sort.Strings(enabledAfGet)
	if enabledAfIntr, ok := d.GetOk("enabled_af"); ok {
		enabledAfAct := make([]string, 0, 1)
		for _, val := range enabledAfIntr.([]interface{}) {
			enabledAfAct = append(enabledAfAct, val.(string))
		}
		sort.Strings(enabledAfAct)
		if reflect.DeepEqual(enabledAfAct, enabledAfGet) {
			d.Set("enabled_af", d.Get("enabled_af").([]interface{}))
		} else {
			d.Set("enabled_af", enabledAfGet)
		}
	} else {
		d.Set("enabled_af", enabledAfGet)
	}
	d.Set("name", pimExtPMap["name"])
	d.Set("name_alias", pimExtPMap["nameAlias"])
	return d, nil
}

func resourceAciExternalProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	pimExtP, err := getRemoteExternalProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setExternalProfileAttributes(pimExtP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciExternalProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ExternalProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	L3OutsideDn := d.Get("l3outside_dn").(string)

	pimExtPAttr := models.ExternalProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimExtPAttr.Annotation = Annotation.(string)
	} else {
		pimExtPAttr.Annotation = "{}"
	}

	if EnabledAf, ok := d.GetOk("enabled_af"); ok {
		enabledAfList := make([]string, 0, 1)
		for _, val := range EnabledAf.([]interface{}) {
			enabledAfList = append(enabledAfList, val.(string))
		}
		EnabledAf := strings.Join(enabledAfList, ",")
		pimExtPAttr.EnabledAf = EnabledAf
	}

	if Name, ok := d.GetOk("name"); ok {
		pimExtPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimExtPAttr.NameAlias = NameAlias.(string)
	}
	pimExtP := models.NewExternalProfile(fmt.Sprintf(models.RnPimExtP), L3OutsideDn, desc, pimExtPAttr)

	err := aciClient.Save(pimExtP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pimExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciExternalProfileRead(ctx, d, m)
}
func resourceAciExternalProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] External Profile: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	L3OutsideDn := d.Get("l3outside_dn").(string)

	pimExtPAttr := models.ExternalProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimExtPAttr.Annotation = Annotation.(string)
	} else {
		pimExtPAttr.Annotation = "{}"
	}

	if EnabledAf, ok := d.GetOk("enabled_af"); ok {
		enabledAfList := make([]string, 0, 1)
		for _, val := range EnabledAf.([]interface{}) {
			enabledAfList = append(enabledAfList, val.(string))
		}
		EnabledAf := strings.Join(enabledAfList, ",")
		pimExtPAttr.EnabledAf = EnabledAf
	}

	if Name, ok := d.GetOk("name"); ok {
		pimExtPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimExtPAttr.NameAlias = NameAlias.(string)
	}
	pimExtP := models.NewExternalProfile(fmt.Sprintf(models.RnPimExtP), L3OutsideDn, desc, pimExtPAttr)

	pimExtP.Status = "modified"

	err := aciClient.Save(pimExtP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pimExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciExternalProfileRead(ctx, d, m)
}

func resourceAciExternalProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	pimExtP, err := getRemoteExternalProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setExternalProfileAttributes(pimExtP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciExternalProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "pimExtP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

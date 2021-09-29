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

func resourceAciRouteControlProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRouteControlProfileCreate,
		UpdateContext: resourceAciRouteControlProfileUpdate,
		ReadContext:   resourceAciRouteControlProfileRead,
		DeleteContext: resourceAciRouteControlProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRouteControlProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
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

			"route_control_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"global",
					"combinable",
				}, false),
			},
		}),
	}
}
func getRemoteRouteControlProfile(client *client.Client, dn string) (*models.RouteControlProfile, error) {
	rtctrlProfileCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlProfile := models.RouteControlProfileFromContainer(rtctrlProfileCont)

	if rtctrlProfile.DistinguishedName == "" {
		return nil, fmt.Errorf("RouteControlProfile %s not found", rtctrlProfile.DistinguishedName)
	}

	return rtctrlProfile, nil
}

func setRouteControlProfileAttributes(rtctrlProfile *models.RouteControlProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(rtctrlProfile.DistinguishedName)
	d.Set("description", rtctrlProfile.Description)
	if dn != rtctrlProfile.DistinguishedName {
		d.Set("parent_dn", "")
	}
	rtctrlProfileMap, err := rtctrlProfile.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("parent_dn", GetParentDn(dn, fmt.Sprintf("/prof-%s", rtctrlProfileMap["name"])))
	d.Set("name", rtctrlProfileMap["name"])
	d.Set("annotation", rtctrlProfileMap["annotation"])
	d.Set("name_alias", rtctrlProfileMap["nameAlias"])
	d.Set("route_control_profile_type", rtctrlProfileMap["type"])

	return d, nil
}

func resourceAciRouteControlProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	rtctrlProfile, err := getRemoteRouteControlProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setRouteControlProfileAttributes(rtctrlProfile, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRouteControlProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RouteControlProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ParentDn := d.Get("parent_dn").(string)

	rtctrlProfileAttr := models.RouteControlProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlProfileAttr.Annotation = Annotation.(string)
	} else {
		rtctrlProfileAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		rtctrlProfileAttr.NameAlias = NameAlias.(string)
	}
	if RouteControlProfileType, ok := d.GetOk("route_control_profile_type"); ok {
		rtctrlProfileAttr.RouteControlProfileType = RouteControlProfileType.(string)
	}
	rtctrlProfile := models.NewRouteControlProfile(fmt.Sprintf("prof-%s", name), ParentDn, desc, rtctrlProfileAttr)

	err := aciClient.Save(rtctrlProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlProfile.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciRouteControlProfileRead(ctx, d, m)
}

func resourceAciRouteControlProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RouteControlProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ParentDn := d.Get("parent_dn").(string)

	rtctrlProfileAttr := models.RouteControlProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlProfileAttr.Annotation = Annotation.(string)
	} else {
		rtctrlProfileAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		rtctrlProfileAttr.NameAlias = NameAlias.(string)
	}
	if RouteControlProfileType, ok := d.GetOk("route_control_profile_type"); ok {
		rtctrlProfileAttr.RouteControlProfileType = RouteControlProfileType.(string)
	}
	rtctrlProfile := models.NewRouteControlProfile(fmt.Sprintf("prof-%s", name), ParentDn, desc, rtctrlProfileAttr)

	rtctrlProfile.Status = "modified"

	err := aciClient.Save(rtctrlProfile)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlProfile.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciRouteControlProfileRead(ctx, d, m)

}

func resourceAciRouteControlProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	rtctrlProfile, err := getRemoteRouteControlProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setRouteControlProfileAttributes(rtctrlProfile, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciRouteControlProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "rtctrlProfile")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

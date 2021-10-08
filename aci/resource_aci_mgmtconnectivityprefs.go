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

func resourceAciMgmtconnectivitypreference() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMgmtconnectivitypreferenceCreate,
		UpdateContext: resourceAciMgmtconnectivitypreferenceUpdate,
		ReadContext:   resourceAciMgmtconnectivitypreferenceRead,
		DeleteContext: resourceAciMgmtconnectivitypreferenceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMgmtconnectivitypreferenceImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"interface_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"inband",
					"ooband",
				}, false),
			},
		})),
	}
}

func getRemoteMgmtconnectivitypreference(client *client.Client, dn string) (*models.Mgmtconnectivitypreference, error) {
	mgmtConnectivityPrefsCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	mgmtConnectivityPrefs := models.MgmtconnectivitypreferenceFromContainer(mgmtConnectivityPrefsCont)
	if mgmtConnectivityPrefs.DistinguishedName == "" {
		return nil, fmt.Errorf("Mgmtconnectivitypreference %s not found", mgmtConnectivityPrefs.DistinguishedName)
	}
	return mgmtConnectivityPrefs, nil
}

func setMgmtconnectivitypreferenceAttributes(mgmtConnectivityPrefs *models.Mgmtconnectivitypreference, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mgmtConnectivityPrefs.DistinguishedName)
	d.Set("description", mgmtConnectivityPrefs.Description)
	mgmtConnectivityPrefsMap, err := mgmtConnectivityPrefs.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", mgmtConnectivityPrefsMap["annotation"])
	d.Set("interface_pref", mgmtConnectivityPrefsMap["interfacePref"])
	d.Set("name_alias", mgmtConnectivityPrefsMap["nameAlias"])
	return d, nil
}

func resourceAciMgmtconnectivitypreferenceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtConnectivityPrefs, err := getRemoteMgmtconnectivitypreference(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMgmtconnectivitypreferenceAttributes(mgmtConnectivityPrefs, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMgmtconnectivitypreferenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Mgmtconnectivitypreference: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	mgmtConnectivityPrefsAttr := models.MgmtconnectivitypreferenceAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtConnectivityPrefsAttr.Annotation = Annotation.(string)
	} else {
		mgmtConnectivityPrefsAttr.Annotation = "{}"
	}

	if InterfacePref, ok := d.GetOk("interface_pref"); ok {
		mgmtConnectivityPrefsAttr.InterfacePref = InterfacePref.(string)
	}

	mgmtConnectivityPrefsAttr.Name = "default"

	mgmtConnectivityPrefs := models.NewMgmtconnectivitypreference(fmt.Sprintf("fabric/connectivityPrefs"), "uni", desc, nameAlias, mgmtConnectivityPrefsAttr)
	mgmtConnectivityPrefs.Status = "modified"
	err := aciClient.Save(mgmtConnectivityPrefs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(mgmtConnectivityPrefs.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMgmtconnectivitypreferenceRead(ctx, d, m)
}

func resourceAciMgmtconnectivitypreferenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Mgmtconnectivitypreference: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	mgmtConnectivityPrefsAttr := models.MgmtconnectivitypreferenceAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtConnectivityPrefsAttr.Annotation = Annotation.(string)
	} else {
		mgmtConnectivityPrefsAttr.Annotation = "{}"
	}

	if InterfacePref, ok := d.GetOk("interface_pref"); ok {
		mgmtConnectivityPrefsAttr.InterfacePref = InterfacePref.(string)
	}

	mgmtConnectivityPrefsAttr.Name = "default"

	mgmtConnectivityPrefs := models.NewMgmtconnectivitypreference(fmt.Sprintf("fabric/connectivityPrefs"), "uni", desc, nameAlias, mgmtConnectivityPrefsAttr)
	mgmtConnectivityPrefs.Status = "modified"
	err := aciClient.Save(mgmtConnectivityPrefs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(mgmtConnectivityPrefs.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMgmtconnectivitypreferenceRead(ctx, d, m)
}

func resourceAciMgmtconnectivitypreferenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtConnectivityPrefs, err := getRemoteMgmtconnectivitypreference(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setMgmtconnectivitypreferenceAttributes(mgmtConnectivityPrefs, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMgmtconnectivitypreferenceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name mgmtConnectivityPrefs cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}

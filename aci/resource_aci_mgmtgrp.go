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

func resourceAciManagedNodeConnectivityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciManagedNodeConnectivityGroupCreate,
		UpdateContext: resourceAciManagedNodeConnectivityGroupUpdate,
		ReadContext:   resourceAciManagedNodeConnectivityGroupRead,
		DeleteContext: resourceAciManagedNodeConnectivityGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciManagedNodeConnectivityGroupImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},
		},
	}
}

func getRemoteManagedNodeConnectivityGroup(client *client.Client, dn string) (*models.ManagedNodeConnectivityGroup, error) {
	mgmtGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	mgmtGrp := models.ManagedNodeConnectivityGroupFromContainer(mgmtGrpCont)
	if mgmtGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("ManagedNodeConnectivityGroup %s not found", mgmtGrp.DistinguishedName)
	}
	return mgmtGrp, nil
}

func setManagedNodeConnectivityGroupAttributes(mgmtGrp *models.ManagedNodeConnectivityGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mgmtGrp.DistinguishedName)
	mgmtGrpMap, err := mgmtGrp.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", mgmtGrpMap["annotation"])
	d.Set("name", mgmtGrpMap["name"])
	return d, nil
}

func resourceAciManagedNodeConnectivityGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtGrp, err := getRemoteManagedNodeConnectivityGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setManagedNodeConnectivityGroupAttributes(mgmtGrp, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciManagedNodeConnectivityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ManagedNodeConnectivityGroup: Beginning Creation")
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)
	mgmtGrpAttr := models.ManagedNodeConnectivityGroupAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtGrpAttr.Annotation = Annotation.(string)
	} else {
		mgmtGrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		mgmtGrpAttr.Name = Name.(string)
	}
	mgmtGrp := models.NewManagedNodeConnectivityGroup(fmt.Sprintf("infra/funcprof/grp-%s", name), "uni", mgmtGrpAttr)
	err := aciClient.Save(mgmtGrp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(mgmtGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciManagedNodeConnectivityGroupRead(ctx, d, m)
}

func resourceAciManagedNodeConnectivityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ManagedNodeConnectivityGroup: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	mgmtGrpAttr := models.ManagedNodeConnectivityGroupAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtGrpAttr.Annotation = Annotation.(string)
	} else {
		mgmtGrpAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		mgmtGrpAttr.Name = Name.(string)
	}
	mgmtGrp := models.NewManagedNodeConnectivityGroup(fmt.Sprintf("infra/funcprof/grp-%s", name), "uni", mgmtGrpAttr)
	mgmtGrp.Status = "modified"
	err := aciClient.Save(mgmtGrp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(mgmtGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciManagedNodeConnectivityGroupRead(ctx, d, m)
}

func resourceAciManagedNodeConnectivityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtGrp, err := getRemoteManagedNodeConnectivityGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	setManagedNodeConnectivityGroupAttributes(mgmtGrp, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciManagedNodeConnectivityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "mgmtGrp")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

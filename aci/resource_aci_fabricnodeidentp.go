package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciFabricNodeMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFabricNodeMemberCreate,
		UpdateContext: resourceAciFabricNodeMemberUpdate,
		ReadContext:   resourceAciFabricNodeMemberRead,
		DeleteContext: resourceAciFabricNodeMemberDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFabricNodeMemberImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"serial": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type: schema.TypeString,
				// Required: true,
				// ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"ext_pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fabric_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"remote-leaf-wan",
				}, false),
			},

			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"leaf",
					"spine",
				}, false),
			},
		}),
	}
}
func getRemoteFabricNodeMember(client *client.Client, dn string) (*models.FabricNodeMember, error) {
	fabricNodeIdentPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNodeIdentP := models.FabricNodeMemberFromContainer(fabricNodeIdentPCont)

	if fabricNodeIdentP.DistinguishedName == "" {
		return nil, fmt.Errorf("Fabric Node Member %s not found", dn)
	}

	return fabricNodeIdentP, nil
}

func setFabricNodeMemberAttributes(fabricNodeIdentP *models.FabricNodeMember, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fabricNodeIdentP.DistinguishedName)
	d.Set("description", fabricNodeIdentP.Description)
	fabricNodeIdentPMap, err := fabricNodeIdentP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("serial", fabricNodeIdentPMap["serial"])
	d.Set("name", fabricNodeIdentPMap["name"])
	d.Set("annotation", fabricNodeIdentPMap["annotation"])
	d.Set("ext_pool_id", fabricNodeIdentPMap["extPoolId"])
	d.Set("fabric_id", fabricNodeIdentPMap["fabricId"])
	d.Set("name_alias", fabricNodeIdentPMap["nameAlias"])
	d.Set("node_id", fabricNodeIdentPMap["nodeId"])
	d.Set("node_type", fabricNodeIdentPMap["nodeType"])
	d.Set("pod_id", fabricNodeIdentPMap["podId"])
	d.Set("role", fabricNodeIdentPMap["role"])
	return d, nil
}

func resourceAciFabricNodeMemberImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fabricNodeIdentP, err := getRemoteFabricNodeMember(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setFabricNodeMemberAttributes(fabricNodeIdentP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFabricNodeMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricNodeMember: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	serial := d.Get("serial").(string)

	name := d.Get("name").(string)

	fabricNodeIdentPAttr := models.FabricNodeMemberAttributes{}
	fabricNodeIdentPAttr.Name = name
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricNodeIdentPAttr.Annotation = Annotation.(string)
	} else {
		fabricNodeIdentPAttr.Annotation = "{}"
	}
	if ExtPoolId, ok := d.GetOk("ext_pool_id"); ok {
		fabricNodeIdentPAttr.ExtPoolId = ExtPoolId.(string)
	}
	if FabricId, ok := d.GetOk("fabric_id"); ok {
		fabricNodeIdentPAttr.FabricId = FabricId.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fabricNodeIdentPAttr.NameAlias = NameAlias.(string)
	}
	if NodeId, ok := d.GetOk("node_id"); ok {
		fabricNodeIdentPAttr.NodeId = NodeId.(string)
	}
	if NodeType, ok := d.GetOk("node_type"); ok {
		fabricNodeIdentPAttr.NodeType = NodeType.(string)
	}
	if PodId, ok := d.GetOk("pod_id"); ok {
		fabricNodeIdentPAttr.PodId = PodId.(string)
	}
	if Role, ok := d.GetOk("role"); ok {
		fabricNodeIdentPAttr.Role = Role.(string)
	}
	fabricNodeIdentP := models.NewFabricNodeMember(fmt.Sprintf("controller/nodeidentpol/nodep-%s", serial), "uni", desc, fabricNodeIdentPAttr)

	err := aciClient.Save(fabricNodeIdentP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricNodeIdentP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFabricNodeMemberRead(ctx, d, m)
}

func resourceAciFabricNodeMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricNodeMember: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	serial := d.Get("serial").(string)
	name := d.Get("name").(string)

	fabricNodeIdentPAttr := models.FabricNodeMemberAttributes{}
	fabricNodeIdentPAttr.Name = name
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricNodeIdentPAttr.Annotation = Annotation.(string)
	} else {
		fabricNodeIdentPAttr.Annotation = "{}"
	}
	if ExtPoolId, ok := d.GetOk("ext_pool_id"); ok {
		fabricNodeIdentPAttr.ExtPoolId = ExtPoolId.(string)
	}
	if FabricId, ok := d.GetOk("fabric_id"); ok {
		fabricNodeIdentPAttr.FabricId = FabricId.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fabricNodeIdentPAttr.NameAlias = NameAlias.(string)
	}
	if NodeId, ok := d.GetOk("node_id"); ok {
		fabricNodeIdentPAttr.NodeId = NodeId.(string)
	}
	if NodeType, ok := d.GetOk("node_type"); ok {
		fabricNodeIdentPAttr.NodeType = NodeType.(string)
	}
	if PodId, ok := d.GetOk("pod_id"); ok {
		fabricNodeIdentPAttr.PodId = PodId.(string)
	}
	if Role, ok := d.GetOk("role"); ok {
		fabricNodeIdentPAttr.Role = Role.(string)
	}
	fabricNodeIdentP := models.NewFabricNodeMember(fmt.Sprintf("controller/nodeidentpol/nodep-%s", serial), "uni", desc, fabricNodeIdentPAttr)

	fabricNodeIdentP.Status = "modified"

	err := aciClient.Save(fabricNodeIdentP)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricNodeIdentP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFabricNodeMemberRead(ctx, d, m)

}

func resourceAciFabricNodeMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fabricNodeIdentP, err := getRemoteFabricNodeMember(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setFabricNodeMemberAttributes(fabricNodeIdentP, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFabricNodeMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fabricNodeIdentP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

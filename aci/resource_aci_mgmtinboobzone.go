package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciManagedNodesZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciManagedNodesZoneCreate,
		UpdateContext: resourceAciManagedNodesZoneUpdate,
		ReadContext:   resourceAciManagedNodesZoneRead,
		DeleteContext: resourceAciManagedNodesZoneDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciManagedNodesZoneImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"managed_node_connectivity_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in_band",
					"out_of_band",
				}, false),
			},

			// relationsships of mgmtInBZone
			"relation_mgmt_rs_addr_inst": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvns:AddrInst",
			},
			"relation_mgmt_rs_in_b": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to mgmt:InB",
			},
			"relation_mgmt_rs_inb_epg": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to mgmt:InB",
			},

			// relationsships of mgmtOoBZone

			"relation_mgmt_rs_oo_b": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to mgmt:OoB",
			},
			"relation_mgmt_rs_oob_epg": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to mgmt:OoB",
			},
		})),
	}
}

func resourceAciManagedNodesZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	splitStr := strings.Split(d.Id(), "/")
	rn := splitStr[len(splitStr)-1]
	if rn == "inbzone" {
		d.Set("type", "in_band")
		d.Set("managed_node_connectivity_group_dn", GetParentDn(d.Id(), "/inbzone"))
		return resourceAciInBManagedNodesZoneImport(d, m)
	} else if rn == "oobzone" {
		d.Set("type", "out_of_band")
		d.Set("managed_node_connectivity_group_dn", GetParentDn(d.Id(), "/oobzone"))
		return resourceAciOOBManagedNodesZoneImport(d, m)
	}

	return nil, fmt.Errorf("Applied DN doesn't belong to mgmtInBZone | mgmtOoBZone")
}

func resourceAciManagedNodesZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Create", d.Id())
	if d.Get("type").(string) == "in_band" {
		return resourceAciInBManagedNodesZoneCreate(ctx, d, m)
	}
	return resourceAciOOBManagedNodesZoneCreate(ctx, d, m)
}

func resourceAciManagedNodesZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Update", d.Id())
	if d.Get("type").(string) == "in_band" {
		return resourceAciInBManagedNodesZoneUpdate(ctx, d, m)
	}
	return resourceAciOOBManagedNodesZoneUpdate(ctx, d, m)
}

func resourceAciManagedNodesZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	if d.Get("type").(string) == "in_band" {
		return resourceAciInBManagedNodesZoneRead(ctx, d, m)
	}
	return resourceAciOOBManagedNodesZoneRead(ctx, d, m)
}

func resourceAciManagedNodesZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	if d.Get("type").(string) == "in_band" {
		return resourceAciInBManagedNodesZoneDelete(ctx, d, m)
	}
	return resourceAciOOBManagedNodesZoneDelete(ctx, d, m)
}

func getRemoteInBManagedNodesZone(client *client.Client, dn string) (*models.InBManagedNodesZone, error) {
	mgmtInBZoneCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	mgmtInBZone := models.InBManagedNodesZoneFromContainer(mgmtInBZoneCont)
	if mgmtInBZone.DistinguishedName == "" {
		return nil, fmt.Errorf("In-Band Managed Nodes Zone %s not found", dn)
	}
	return mgmtInBZone, nil
}

func getRemoteOOBManagedNodesZone(client *client.Client, dn string) (*models.OOBManagedNodesZone, error) {
	mgmtOobZoneCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	mgmtOobZone := models.OOBManagedNodesZoneFromContainer(mgmtOobZoneCont)
	if mgmtOobZone.DistinguishedName == "" {
		return nil, fmt.Errorf("Out-of-Band Managed Nodes Zone %s not found", dn)
	}
	return mgmtOobZone, nil
}

func setInBManagedNodesZoneAttributes(mgmtInBZone *models.InBManagedNodesZone, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mgmtInBZone.DistinguishedName)
	d.Set("description", mgmtInBZone.Description)
	mgmtInBZoneMap, err := mgmtInBZone.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", mgmtInBZoneMap["annotation"])
	d.Set("name", mgmtInBZoneMap["name"])
	d.Set("name_alias", mgmtInBZoneMap["nameAlias"])
	return d, nil
}

func setOOBManagedNodesZoneAttributes(mgmtOobZone *models.OOBManagedNodesZone, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mgmtOobZone.DistinguishedName)
	d.Set("description", mgmtOobZone.Description)
	mgmtOobZoneMap, err := mgmtOobZone.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", mgmtOobZoneMap["annotation"])
	d.Set("name", mgmtOobZoneMap["name"])
	d.Set("name_alias", mgmtOobZoneMap["nameAlias"])
	return d, nil
}
func resourceAciInBManagedNodesZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtInBZone, err := getRemoteInBManagedNodesZone(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setInBManagedNodesZoneAttributes(mgmtInBZone, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOOBManagedNodesZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtOoBZone, err := getRemoteOOBManagedNodesZone(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setOOBManagedNodesZoneAttributes(mgmtOoBZone, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciInBManagedNodesZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InBManagedNodesZone: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ManagedNodeConnectivityGroupDn := d.Get("managed_node_connectivity_group_dn").(string)

	mgmtInBZoneAttr := models.InBManagedNodesZoneAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtInBZoneAttr.Annotation = Annotation.(string)
	} else {
		mgmtInBZoneAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		mgmtInBZoneAttr.Name = Name.(string)
	}
	mgmtInBZone := models.NewInBManagedNodesZone(fmt.Sprintf("inbzone"), ManagedNodeConnectivityGroupDn, desc, nameAlias, mgmtInBZoneAttr)

	err := aciClient.Save(mgmtInBZone)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTomgmtRsAddrInst, ok := d.GetOk("relation_mgmt_rs_addr_inst"); ok {
		relationParam := relationTomgmtRsAddrInst.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTomgmtRsInB, ok := d.GetOk("relation_mgmt_rs_in_b"); ok {
		relationParam := relationTomgmtRsInB.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTomgmtRsInbEpg, ok := d.GetOk("relation_mgmt_rs_inb_epg"); ok {
		relationParam := relationTomgmtRsInbEpg.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTomgmtRsAddrInst, ok := d.GetOk("relation_mgmt_rs_addr_inst"); ok {
		relationParam := relationTomgmtRsAddrInst.(string)
		err = aciClient.CreateRelationmgmtRsAddrInst(mgmtInBZone.DistinguishedName, mgmtInBZoneAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTomgmtRsInB, ok := d.GetOk("relation_mgmt_rs_in_b"); ok {
		relationParam := relationTomgmtRsInB.(string)
		err = aciClient.CreateRelationmgmtRsInB(mgmtInBZone.DistinguishedName, mgmtInBZoneAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTomgmtRsInbEpg, ok := d.GetOk("relation_mgmt_rs_inb_epg"); ok {
		relationParam := relationTomgmtRsInbEpg.(string)
		err = aciClient.CreateRelationmgmtRsInbEpg(mgmtInBZone.DistinguishedName, mgmtInBZoneAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(mgmtInBZone.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciInBManagedNodesZoneRead(ctx, d, m)
}

func resourceAciOOBManagedNodesZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OOBManagedNodesZone: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ManagedNodeConnectivityGroupDn := d.Get("managed_node_connectivity_group_dn").(string)

	mgmtOoBZoneAttr := models.OOBManagedNodesZoneAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtOoBZoneAttr.Annotation = Annotation.(string)
	} else {
		mgmtOoBZoneAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		mgmtOoBZoneAttr.Name = Name.(string)
	}
	mgmtOoBZone := models.NewOOBManagedNodesZone(fmt.Sprintf("oobzone"), ManagedNodeConnectivityGroupDn, desc, nameAlias, mgmtOoBZoneAttr)

	err := aciClient.Save(mgmtOoBZone)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTomgmtRsAddrInst, ok := d.GetOk("relation_mgmt_rs_addr_inst"); ok {
		relationParam := relationTomgmtRsAddrInst.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTomgmtRsOoB, ok := d.GetOk("relation_mgmt_rs_oo_b"); ok {
		relationParam := relationTomgmtRsOoB.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTomgmtRsOobEpg, ok := d.GetOk("relation_mgmt_rs_oob_epg"); ok {
		relationParam := relationTomgmtRsOobEpg.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTomgmtRsAddrInst, ok := d.GetOk("relation_mgmt_rs_addr_inst"); ok {
		relationParam := relationTomgmtRsAddrInst.(string)
		err = aciClient.CreateRelationmgmtRsAddrInst(mgmtOoBZone.DistinguishedName, mgmtOoBZoneAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTomgmtRsOoB, ok := d.GetOk("relation_mgmt_rs_oo_b"); ok {
		relationParam := relationTomgmtRsOoB.(string)
		err = aciClient.CreateRelationmgmtRsOoB(mgmtOoBZone.DistinguishedName, mgmtOoBZoneAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTomgmtRsOobEpg, ok := d.GetOk("relation_mgmt_rs_oob_epg"); ok {
		relationParam := relationTomgmtRsOobEpg.(string)
		err = aciClient.CreateRelationmgmtRsOobEpg(mgmtOoBZone.DistinguishedName, mgmtOoBZoneAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(mgmtOoBZone.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciOOBManagedNodesZoneRead(ctx, d, m)
}

func resourceAciInBManagedNodesZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InBManagedNodesZone: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ManagedNodeConnectivityGroupDn := d.Get("managed_node_connectivity_group_dn").(string)
	mgmtInBZoneAttr := models.InBManagedNodesZoneAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtInBZoneAttr.Annotation = Annotation.(string)
	} else {
		mgmtInBZoneAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		mgmtInBZoneAttr.Name = Name.(string)
	}
	mgmtInBZone := models.NewInBManagedNodesZone(fmt.Sprintf("inbzone"), ManagedNodeConnectivityGroupDn, desc, nameAlias, mgmtInBZoneAttr)

	mgmtInBZone.Status = "modified"
	err := aciClient.Save(mgmtInBZone)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_mgmt_rs_addr_inst") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_addr_inst")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_mgmt_rs_in_b") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_in_b")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_mgmt_rs_inb_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_inb_epg")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_mgmt_rs_addr_inst") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_addr_inst")
		err = aciClient.DeleteRelationmgmtRsAddrInst(mgmtInBZone.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationmgmtRsAddrInst(mgmtInBZone.DistinguishedName, mgmtInBZoneAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_mgmt_rs_in_b") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_in_b")
		err = aciClient.DeleteRelationmgmtRsInB(mgmtInBZone.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationmgmtRsInB(mgmtInBZone.DistinguishedName, mgmtInBZoneAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_mgmt_rs_inb_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_inb_epg")
		err = aciClient.DeleteRelationmgmtRsInbEpg(mgmtInBZone.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationmgmtRsInbEpg(mgmtInBZone.DistinguishedName, mgmtInBZoneAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(mgmtInBZone.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciInBManagedNodesZoneRead(ctx, d, m)
}

func resourceAciOOBManagedNodesZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OOBManagedNodesZone: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ManagedNodeConnectivityGroupDn := d.Get("managed_node_connectivity_group_dn").(string)
	mgmtOoBZoneAttr := models.OOBManagedNodesZoneAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		mgmtOoBZoneAttr.Annotation = Annotation.(string)
	} else {
		mgmtOoBZoneAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		mgmtOoBZoneAttr.Name = Name.(string)
	}
	mgmtOoBZone := models.NewOOBManagedNodesZone(fmt.Sprintf("oobzone"), ManagedNodeConnectivityGroupDn, desc, nameAlias, mgmtOoBZoneAttr)

	mgmtOoBZone.Status = "modified"
	err := aciClient.Save(mgmtOoBZone)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_mgmt_rs_addr_inst") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_addr_inst")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_mgmt_rs_oo_b") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_oo_b")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_mgmt_rs_oob_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_oob_epg")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_mgmt_rs_addr_inst") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_addr_inst")
		err = aciClient.DeleteRelationmgmtRsAddrInst(mgmtOoBZone.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationmgmtRsAddrInst(mgmtOoBZone.DistinguishedName, mgmtOoBZoneAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_mgmt_rs_oo_b") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_oo_b")
		err = aciClient.DeleteRelationmgmtRsOoB(mgmtOoBZone.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationmgmtRsOoB(mgmtOoBZone.DistinguishedName, mgmtOoBZoneAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_mgmt_rs_oob_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_mgmt_rs_oob_epg")
		err = aciClient.DeleteRelationmgmtRsOobEpg(mgmtOoBZone.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationmgmtRsOobEpg(mgmtOoBZone.DistinguishedName, mgmtOoBZoneAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(mgmtOoBZone.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciOOBManagedNodesZoneRead(ctx, d, m)
}

func resourceAciInBManagedNodesZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtInBZone, err := getRemoteInBManagedNodesZone(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setInBManagedNodesZoneAttributes(mgmtInBZone, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	mgmtRsAddrInstData, err := aciClient.ReadRelationmgmtRsAddrInst(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsAddrInst %v", err)
		d.Set("relation_mgmt_rs_addr_inst", "")
	} else {
		setRelationAttribute(d, "relation_mgmt_rs_addr_inst", mgmtRsAddrInstData)
	}

	mgmtRsInBData, err := aciClient.ReadRelationmgmtRsInB(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsInB %v", err)
		d.Set("relation_mgmt_rs_in_b", "")
	} else {
		setRelationAttribute(d, "relation_mgmt_rs_in_b", mgmtRsInBData)
	}

	mgmtRsInbEpgData, err := aciClient.ReadRelationmgmtRsInbEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsInbEpg %v", err)
		d.Set("relation_mgmt_rs_inb_epg", "")
	} else {
		setRelationAttribute(d, "relation_mgmt_rs_inb_epg", mgmtRsInbEpgData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciOOBManagedNodesZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mgmtOoBZone, err := getRemoteOOBManagedNodesZone(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setOOBManagedNodesZoneAttributes(mgmtOoBZone, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	mgmtRsAddrInstData, err := aciClient.ReadRelationmgmtRsAddrInst(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsAddrInst %v", err)
		d.Set("relation_mgmt_rs_addr_inst", "")
	} else {
		setRelationAttribute(d, "relation_mgmt_rs_addr_inst", mgmtRsAddrInstData)
	}

	mgmtRsOoBData, err := aciClient.ReadRelationmgmtRsOoB(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsOoB %v", err)
		d.Set("relation_mgmt_rs_oo_b", "")
	} else {
		setRelationAttribute(d, "relation_mgmt_rs_oo_b", mgmtRsOoBData)
	}

	mgmtRsOobEpgData, err := aciClient.ReadRelationmgmtRsOobEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation mgmtRsOobEpg %v", err)
		d.Set("relation_mgmt_rs_oob_epg", "")
	} else {
		setRelationAttribute(d, "relation_mgmt_rs_oob_epg", mgmtRsOobEpgData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciInBManagedNodesZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "mgmtInBZone")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

func resourceAciOOBManagedNodesZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "mgmtOoBZone")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

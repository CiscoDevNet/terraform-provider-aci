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

func resourceAciLogicalNodeProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLogicalNodeProfileCreate,
		UpdateContext: resourceAciLogicalNodeProfileUpdate,
		ReadContext:   resourceAciLogicalNodeProfileRead,
		DeleteContext: resourceAciLogicalNodeProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLogicalNodeProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"config_issues": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressBitMaskDiffFunc(),
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"none",
					"node-path-misconfig",
					"routerid-not-changable-with-mcast",
					"loopback-ip-missing",
				}, false, "none")),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CS0",
					"CS1",
					"AF11",
					"AF12",
					"AF13",
					"CS2",
					"AF21",
					"AF22",
					"AF23",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"VA",
					"EF",
					"unspecified",
				}, false),
			},
		}),
	}
}
func getRemoteLogicalNodeProfile(client *client.Client, dn string) (*models.LogicalNodeProfile, error) {
	l3extLNodePCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extLNodeP := models.LogicalNodeProfileFromContainer(l3extLNodePCont)

	if l3extLNodeP.DistinguishedName == "" {
		return nil, fmt.Errorf("LogicalNodeProfile %s not found", l3extLNodeP.DistinguishedName)
	}

	return l3extLNodeP, nil
}

func setLogicalNodeProfileAttributes(l3extLNodeP *models.LogicalNodeProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(l3extLNodeP.DistinguishedName)
	d.Set("description", l3extLNodeP.Description)

	if dn != l3extLNodeP.DistinguishedName {
		d.Set("l3_outside_dn", "")
	}
	l3extLNodePMap, err := l3extLNodeP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("l3_outside_dn", GetParentDn(dn, fmt.Sprintf("/lnodep-%s",  l3extLNodePMap["name"])))
	d.Set("name", l3extLNodePMap["name"])

	d.Set("annotation", l3extLNodePMap["annotation"])
	if l3extLNodePMap["configIssues"] == "" {
		d.Set("config_issues", "none")
	} else {
		d.Set("config_issues", l3extLNodePMap["configIssues"])
	}
	d.Set("name_alias", l3extLNodePMap["nameAlias"])
	d.Set("tag", l3extLNodePMap["tag"])
	d.Set("target_dscp", l3extLNodePMap["targetDscp"])
	return d, nil
}

func resourceAciLogicalNodeProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extLNodeP, err := getRemoteLogicalNodeProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	l3extLNodePMap, err := l3extLNodeP.ToMap()
	if err != nil {
		return nil, err
	}
	name := l3extLNodePMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/lnodep-%s", name))
	d.Set("l3_outside_dn", pDN)
	schemaFilled, err := setLogicalNodeProfileAttributes(l3extLNodeP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLogicalNodeProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalNodeProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extLNodePAttr := models.LogicalNodeProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLNodePAttr.Annotation = Annotation.(string)
	} else {
		l3extLNodePAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		l3extLNodePAttr.ConfigIssues = ConfigIssues.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extLNodePAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		l3extLNodePAttr.Tag = Tag.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extLNodePAttr.TargetDscp = TargetDscp.(string)
	}
	l3extLNodeP := models.NewLogicalNodeProfile(fmt.Sprintf("lnodep-%s", name), L3OutsideDn, desc, l3extLNodePAttr)

	err := aciClient.Save(l3extLNodeP)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationTol3extRsNodeL3OutAtt, ok := d.GetOk("relation_l3ext_rs_node_l3_out_att"); ok {
		relationParamList := toStringList(relationTol3extRsNodeL3OutAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(l3extLNodeP.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(l3extLNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLogicalNodeProfileRead(ctx, d, m)
}

func resourceAciLogicalNodeProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalNodeProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extLNodePAttr := models.LogicalNodeProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLNodePAttr.Annotation = Annotation.(string)
	} else {
		l3extLNodePAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		l3extLNodePAttr.ConfigIssues = ConfigIssues.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extLNodePAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		l3extLNodePAttr.Tag = Tag.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extLNodePAttr.TargetDscp = TargetDscp.(string)
	}
	l3extLNodeP := models.NewLogicalNodeProfile(fmt.Sprintf("lnodep-%s", name), L3OutsideDn, desc, l3extLNodePAttr)

	l3extLNodeP.Status = "modified"

	err := aciClient.Save(l3extLNodeP)

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_l3ext_rs_node_l3_out_att") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_node_l3_out_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(l3extLNodeP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(l3extLNodeP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(l3extLNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLogicalNodeProfileRead(ctx, d, m)

}

func resourceAciLogicalNodeProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extLNodeP, err := getRemoteLogicalNodeProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLogicalNodeProfileAttributes(l3extLNodeP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	l3extRsNodeL3OutAttData, err := aciClient.ReadRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsNodeL3OutAtt %v", err)

	} else {
		d.Set("relation_l3ext_rs_node_l3_out_att", l3extRsNodeL3OutAttData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLogicalNodeProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extLNodeP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

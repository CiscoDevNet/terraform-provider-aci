package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciLogicalNodeProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLogicalNodeProfileCreate,
		Update: resourceAciLogicalNodeProfileUpdate,
		Read:   resourceAciLogicalNodeProfileRead,
		Delete: resourceAciLogicalNodeProfileDelete,

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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"config_issues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func setLogicalNodeProfileAttributes(l3extLNodeP *models.LogicalNodeProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l3extLNodeP.DistinguishedName)
	d.Set("description", l3extLNodeP.Description)
	d.Set("l3_outside_dn", GetParentDn(l3extLNodeP.DistinguishedName))
	l3extLNodePMap, _ := l3extLNodeP.ToMap()

	d.Set("name", l3extLNodePMap["name"])

	d.Set("annotation", l3extLNodePMap["annotation"])
	d.Set("config_issues", l3extLNodePMap["configIssues"])
	d.Set("name_alias", l3extLNodePMap["nameAlias"])
	d.Set("tag", l3extLNodePMap["tag"])
	d.Set("target_dscp", l3extLNodePMap["targetDscp"])
	return d
}

func resourceAciLogicalNodeProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extLNodeP, err := getRemoteLogicalNodeProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLogicalNodeProfileAttributes(l3extLNodeP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLogicalNodeProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LogicalNodeProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extLNodePAttr := models.LogicalNodeProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLNodePAttr.Annotation = Annotation.(string)
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTol3extRsNodeL3OutAtt, ok := d.GetOk("relation_l3ext_rs_node_l3_out_att"); ok {
		relationParamList := toStringList(relationTol3extRsNodeL3OutAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(l3extLNodeP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_node_l3_out_att")
			d.Partial(false)
		}
	}

	d.SetId(l3extLNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLogicalNodeProfileRead(d, m)
}

func resourceAciLogicalNodeProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LogicalNodeProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	l3extLNodePAttr := models.LogicalNodeProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLNodePAttr.Annotation = Annotation.(string)
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_l3ext_rs_node_l3_out_att") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_node_l3_out_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(l3extLNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(l3extLNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_node_l3_out_att")
			d.Partial(false)

		}

	}

	d.SetId(l3extLNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLogicalNodeProfileRead(d, m)

}

func resourceAciLogicalNodeProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extLNodeP, err := getRemoteLogicalNodeProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLogicalNodeProfileAttributes(l3extLNodeP, d)

	l3extRsNodeL3OutAttData, err := aciClient.ReadRelationl3extRsNodeL3OutAttFromLogicalNodeProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsNodeL3OutAtt %v", err)

	} else {
		d.Set("relation_l3ext_rs_node_l3_out_att", l3extRsNodeL3OutAttData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLogicalNodeProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extLNodeP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

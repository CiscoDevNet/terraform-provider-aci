package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciL3ExtSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3ExtSubnetCreate,
		Update: resourceAciL3ExtSubnetUpdate,
		Read:   resourceAciL3ExtSubnetRead,
		Delete: resourceAciL3ExtSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3ExtSubnetImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"external_network_instance_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"aggregate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_l3ext_rs_subnet_to_profile": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_l3ext_rs_subnet_to_rt_summ": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteL3ExtSubnet(client *client.Client, dn string) (*models.L3ExtSubnet, error) {
	l3extSubnetCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extSubnet := models.L3ExtSubnetFromContainer(l3extSubnetCont)

	if l3extSubnet.DistinguishedName == "" {
		return nil, fmt.Errorf("Subnet %s not found", l3extSubnet.DistinguishedName)
	}

	return l3extSubnet, nil
}

func setL3ExtSubnetAttributes(l3extSubnet *models.L3ExtSubnet, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l3extSubnet.DistinguishedName)
	d.Set("description", l3extSubnet.Description)
	d.Set("external_network_instance_profile_dn", GetParentDn(l3extSubnet.DistinguishedName))
	l3extSubnetMap, _ := l3extSubnet.ToMap()

	d.Set("ip", l3extSubnetMap["ip"])

	d.Set("aggregate", l3extSubnetMap["aggregate"])
	d.Set("annotation", l3extSubnetMap["annotation"])
	d.Set("ip", l3extSubnetMap["ip"])
	d.Set("name_alias", l3extSubnetMap["nameAlias"])
	d.Set("scope", l3extSubnetMap["scope"])
	return d
}

func resourceAciL3ExtSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extSubnet, err := getRemoteL3ExtSubnet(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3ExtSubnetAttributes(l3extSubnet, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3ExtSubnetCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Subnet: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	ExternalNetworkInstanceProfileDn := d.Get("external_network_instance_profile_dn").(string)

	l3extSubnetAttr := models.L3ExtSubnetAttributes{}
	if Aggregate, ok := d.GetOk("aggregate"); ok {
		l3extSubnetAttr.Aggregate = Aggregate.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extSubnetAttr.Annotation = Annotation.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		l3extSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		l3extSubnetAttr.Scope = Scope.(string)
	}
	l3extSubnet := models.NewL3ExtSubnet(fmt.Sprintf("extsubnet-[%s]", ip), ExternalNetworkInstanceProfileDn, desc, l3extSubnetAttr)

	err := aciClient.Save(l3extSubnet)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("ip")

	d.Partial(false)

	if relationTol3extRsSubnetToProfile, ok := d.GetOk("relation_l3ext_rs_subnet_to_profile"); ok {

		relationParamList := relationTol3extRsSubnetToProfile.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsSubnetToProfileFromL3ExtSubnet(l3extSubnet.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_subnet_to_profile")
			d.Partial(false)
		}

	}
	if relationTol3extRsSubnetToRtSumm, ok := d.GetOk("relation_l3ext_rs_subnet_to_rt_summ"); ok {
		relationParam := relationTol3extRsSubnetToRtSumm.(string)
		err = aciClient.CreateRelationl3extRsSubnetToRtSummFromL3ExtSubnet(l3extSubnet.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_subnet_to_rt_summ")
		d.Partial(false)

	}

	d.SetId(l3extSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3ExtSubnetRead(d, m)
}

func resourceAciL3ExtSubnetUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Subnet: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	ExternalNetworkInstanceProfileDn := d.Get("external_network_instance_profile_dn").(string)

	l3extSubnetAttr := models.L3ExtSubnetAttributes{}
	if Aggregate, ok := d.GetOk("aggregate"); ok {
		l3extSubnetAttr.Aggregate = Aggregate.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extSubnetAttr.Annotation = Annotation.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		l3extSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		l3extSubnetAttr.Scope = Scope.(string)
	}
	l3extSubnet := models.NewL3ExtSubnet(fmt.Sprintf("extsubnet-[%s]", ip), ExternalNetworkInstanceProfileDn, desc, l3extSubnetAttr)

	l3extSubnet.Status = "modified"

	err := aciClient.Save(l3extSubnet)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("ip")

	d.Partial(false)

	if d.HasChange("relation_l3ext_rs_subnet_to_profile") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_subnet_to_profile")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsSubnetToProfileFromL3ExtSubnet(l3extSubnet.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsSubnetToProfileFromL3ExtSubnet(l3extSubnet.DistinguishedName, paramMap["tn_rtctrl_profile_name"].(string), paramMap["direction"].(string))
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_l3ext_rs_subnet_to_profile")
			d.Partial(false)
		}

	}
	if d.HasChange("relation_l3ext_rs_subnet_to_rt_summ") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_subnet_to_rt_summ")
		err = aciClient.DeleteRelationl3extRsSubnetToRtSummFromL3ExtSubnet(l3extSubnet.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationl3extRsSubnetToRtSummFromL3ExtSubnet(l3extSubnet.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_l3ext_rs_subnet_to_rt_summ")
		d.Partial(false)

	}

	d.SetId(l3extSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3ExtSubnetRead(d, m)

}

func resourceAciL3ExtSubnetRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extSubnet, err := getRemoteL3ExtSubnet(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3ExtSubnetAttributes(l3extSubnet, d)

	l3extRsSubnetToProfileData, err := aciClient.ReadRelationl3extRsSubnetToProfileFromL3ExtSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsSubnetToProfile %v", err)

	} else {
		d.Set("relation_l3ext_rs_subnet_to_profile", l3extRsSubnetToProfileData)
	}

	l3extRsSubnetToRtSummData, err := aciClient.ReadRelationl3extRsSubnetToRtSummFromL3ExtSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsSubnetToRtSumm %v", err)

	} else {
		d.Set("relation_l3ext_rs_subnet_to_rt_summ", l3extRsSubnetToRtSummData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3ExtSubnetDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extSubnet")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

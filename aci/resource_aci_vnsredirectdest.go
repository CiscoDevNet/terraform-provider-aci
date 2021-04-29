package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciDestinationofredirectedtraffic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciDestinationofredirectedtrafficCreate,
		Update: resourceAciDestinationofredirectedtrafficUpdate,
		Read:   resourceAciDestinationofredirectedtrafficRead,
		Delete: resourceAciDestinationofredirectedtrafficDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDestinationofredirectedtrafficImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"service_redirect_policy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dest_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vns_rs_redirect_health_group": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteDestinationofredirectedtraffic(client *client.Client, dn string) (*models.Destinationofredirectedtraffic, error) {
	vnsRedirectDestCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsRedirectDest := models.DestinationofredirectedtrafficFromContainer(vnsRedirectDestCont)

	if vnsRedirectDest.DistinguishedName == "" {
		return nil, fmt.Errorf("Destinationofredirectedtraffic %s not found", vnsRedirectDest.DistinguishedName)
	}

	return vnsRedirectDest, nil
}

func setDestinationofredirectedtrafficAttributes(vnsRedirectDest *models.Destinationofredirectedtraffic, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(vnsRedirectDest.DistinguishedName)
	d.Set("description", vnsRedirectDest.Description)
	if dn != vnsRedirectDest.DistinguishedName {
		d.Set("service_redirect_policy_dn", "")
	}
	vnsRedirectDestMap, _ := vnsRedirectDest.ToMap()

	d.Set("ip", vnsRedirectDestMap["ip"])

	d.Set("annotation", vnsRedirectDestMap["annotation"])
	d.Set("dest_name", vnsRedirectDestMap["destName"])
	d.Set("ip", vnsRedirectDestMap["ip"])
	d.Set("ip2", vnsRedirectDestMap["ip2"])
	d.Set("mac", vnsRedirectDestMap["mac"])
	d.Set("name_alias", vnsRedirectDestMap["nameAlias"])
	d.Set("pod_id", vnsRedirectDestMap["podId"])
	return d
}

func resourceAciDestinationofredirectedtrafficImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsRedirectDest, err := getRemoteDestinationofredirectedtraffic(aciClient, dn)

	if err != nil {
		return nil, err
	}
	vnsRedirectDestMap, _ := vnsRedirectDest.ToMap()
	ip := vnsRedirectDestMap["ip"]
	pDN := GetParentDn(dn, fmt.Sprintf("/RedirectDest_ip-[%s]", ip))
	d.Set("service_redirect_policy_dn", pDN)
	schemaFilled := setDestinationofredirectedtrafficAttributes(vnsRedirectDest, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDestinationofredirectedtrafficCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Destinationofredirectedtraffic: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	ServiceRedirectPolicyDn := d.Get("service_redirect_policy_dn").(string)

	vnsRedirectDestAttr := models.DestinationofredirectedtrafficAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsRedirectDestAttr.Annotation = Annotation.(string)
	} else {
		vnsRedirectDestAttr.Annotation = "{}"
	}
	if DestName, ok := d.GetOk("dest_name"); ok {
		vnsRedirectDestAttr.DestName = DestName.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		vnsRedirectDestAttr.Ip = Ip.(string)
	}
	if Ip2, ok := d.GetOk("ip2"); ok {
		vnsRedirectDestAttr.Ip2 = Ip2.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		vnsRedirectDestAttr.Mac = Mac.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsRedirectDestAttr.NameAlias = NameAlias.(string)
	}
	if PodId, ok := d.GetOk("pod_id"); ok {
		vnsRedirectDestAttr.PodId = PodId.(string)
	}
	vnsRedirectDest := models.NewDestinationofredirectedtraffic(fmt.Sprintf("RedirectDest_ip-[%s]", ip), ServiceRedirectPolicyDn, desc, vnsRedirectDestAttr)

	err := aciClient.Save(vnsRedirectDest)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationTovnsRsRedirectHealthGroup, ok := d.GetOk("relation_vns_rs_redirect_health_group"); ok {
		relationParam := relationTovnsRsRedirectHealthGroup.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationTovnsRsRedirectHealthGroup, ok := d.GetOk("relation_vns_rs_redirect_health_group"); ok {
		relationParam := relationTovnsRsRedirectHealthGroup.(string)
		err = aciClient.CreateRelationvnsRsRedirectHealthGroupFromDestinationofredirectedtraffic(vnsRedirectDest.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(vnsRedirectDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciDestinationofredirectedtrafficRead(d, m)
}

func resourceAciDestinationofredirectedtrafficUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Destinationofredirectedtraffic: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	ServiceRedirectPolicyDn := d.Get("service_redirect_policy_dn").(string)

	vnsRedirectDestAttr := models.DestinationofredirectedtrafficAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsRedirectDestAttr.Annotation = Annotation.(string)
	} else {
		vnsRedirectDestAttr.Annotation = "{}"
	}
	if DestName, ok := d.GetOk("dest_name"); ok {
		vnsRedirectDestAttr.DestName = DestName.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		vnsRedirectDestAttr.Ip = Ip.(string)
	}
	if Ip2, ok := d.GetOk("ip2"); ok {
		vnsRedirectDestAttr.Ip2 = Ip2.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		vnsRedirectDestAttr.Mac = Mac.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsRedirectDestAttr.NameAlias = NameAlias.(string)
	}
	if PodId, ok := d.GetOk("pod_id"); ok {
		vnsRedirectDestAttr.PodId = PodId.(string)
	}
	vnsRedirectDest := models.NewDestinationofredirectedtraffic(fmt.Sprintf("RedirectDest_ip-[%s]", ip), ServiceRedirectPolicyDn, desc, vnsRedirectDestAttr)

	vnsRedirectDest.Status = "modified"

	err := aciClient.Save(vnsRedirectDest)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_redirect_health_group") {
		_, newRelParam := d.GetChange("relation_vns_rs_redirect_health_group")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_redirect_health_group") {
		_, newRelParam := d.GetChange("relation_vns_rs_redirect_health_group")
		err = aciClient.DeleteRelationvnsRsRedirectHealthGroupFromDestinationofredirectedtraffic(vnsRedirectDest.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvnsRsRedirectHealthGroupFromDestinationofredirectedtraffic(vnsRedirectDest.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(vnsRedirectDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciDestinationofredirectedtrafficRead(d, m)

}

func resourceAciDestinationofredirectedtrafficRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsRedirectDest, err := getRemoteDestinationofredirectedtraffic(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setDestinationofredirectedtrafficAttributes(vnsRedirectDest, d)

	vnsRsRedirectHealthGroupData, err := aciClient.ReadRelationvnsRsRedirectHealthGroupFromDestinationofredirectedtraffic(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsRedirectHealthGroup %v", err)
		d.Set("relation_vns_rs_redirect_health_group", "")

	} else {
		d.Set("relation_vns_rs_redirect_health_group", vnsRsRedirectHealthGroupData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciDestinationofredirectedtrafficDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsRedirectDest")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

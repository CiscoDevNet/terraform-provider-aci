package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciL3outStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3outStaticRouteCreate,
		Update: resourceAciL3outStaticRouteUpdate,
		Read:   resourceAciL3outStaticRouteRead,
		Delete: resourceAciL3outStaticRouteDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outStaticRouteImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"fabric_node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"aggregate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rt_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bfd",
					"unspecified",
				}, false),
			},

			"relation_ip_rs_route_track": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteL3outStaticRoute(client *client.Client, dn string) (*models.L3outStaticRoute, error) {
	ipRoutePCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	ipRouteP := models.L3outStaticRouteFromContainer(ipRoutePCont)

	if ipRouteP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outStaticRoute %s not found", ipRouteP.DistinguishedName)
	}

	return ipRouteP, nil
}

func setL3outStaticRouteAttributes(ipRouteP *models.L3outStaticRoute, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(ipRouteP.DistinguishedName)
	d.Set("description", ipRouteP.Description)
	dn := d.Id()
	if dn != ipRouteP.DistinguishedName {
		d.Set("fabric_node_dn", "")
	}
	ipRoutePMap, _ := ipRouteP.ToMap()

	d.Set("ip", ipRoutePMap["ip"])

	d.Set("aggregate", ipRoutePMap["aggregate"])
	d.Set("annotation", ipRoutePMap["annotation"])
	d.Set("ip", ipRoutePMap["ip"])
	d.Set("name_alias", ipRoutePMap["nameAlias"])
	d.Set("pref", ipRoutePMap["pref"])
	if ipRoutePMap["rtCtrl"] == "" {
		d.Set("rt_ctrl", "unspecified")
	} else {
		d.Set("rt_ctrl", ipRoutePMap["rtCtrl"])
	}
	return d
}

func resourceAciL3outStaticRouteImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ipRouteP, err := getRemoteL3outStaticRoute(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3outStaticRouteAttributes(ipRouteP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outStaticRouteCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outStaticRoute: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	FabricNodeDn := d.Get("fabric_node_dn").(string)

	ipRoutePAttr := models.L3outStaticRouteAttributes{}
	if Aggregate, ok := d.GetOk("aggregate"); ok {
		ipRoutePAttr.Aggregate = Aggregate.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ipRoutePAttr.Annotation = Annotation.(string)
	} else {
		ipRoutePAttr.Annotation = "{}"
	}
	if Ip, ok := d.GetOk("ip"); ok {
		ipRoutePAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ipRoutePAttr.NameAlias = NameAlias.(string)
	}
	if Pref, ok := d.GetOk("pref"); ok {
		ipRoutePAttr.Pref = Pref.(string)
	}
	if RtCtrl, ok := d.GetOk("rt_ctrl"); ok {
		ipRoutePAttr.RtCtrl = RtCtrl.(string)
	}
	ipRouteP := models.NewL3outStaticRoute(fmt.Sprintf("rt-[%s]", ip), FabricNodeDn, desc, ipRoutePAttr)

	err := aciClient.Save(ipRouteP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("ip")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationToipRsRouteTrack, ok := d.GetOk("relation_ip_rs_route_track"); ok {
		relationParam := relationToipRsRouteTrack.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationToipRsRouteTrack, ok := d.GetOk("relation_ip_rs_route_track"); ok {
		relationParam := relationToipRsRouteTrack.(string)
		err = aciClient.CreateRelationipRsRouteTrackFromL3outStaticRoute(ipRouteP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ip_rs_route_track")
		d.Partial(false)

	}

	d.SetId(ipRouteP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outStaticRouteRead(d, m)
}

func resourceAciL3outStaticRouteUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outStaticRoute: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	FabricNodeDn := d.Get("fabric_node_dn").(string)

	ipRoutePAttr := models.L3outStaticRouteAttributes{}
	if Aggregate, ok := d.GetOk("aggregate"); ok {
		ipRoutePAttr.Aggregate = Aggregate.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ipRoutePAttr.Annotation = Annotation.(string)
	} else {
		ipRoutePAttr.Annotation = "{}"
	}
	if Ip, ok := d.GetOk("ip"); ok {
		ipRoutePAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ipRoutePAttr.NameAlias = NameAlias.(string)
	}
	if Pref, ok := d.GetOk("pref"); ok {
		ipRoutePAttr.Pref = Pref.(string)
	}
	if RtCtrl, ok := d.GetOk("rt_ctrl"); ok {
		ipRoutePAttr.RtCtrl = RtCtrl.(string)
	}
	ipRouteP := models.NewL3outStaticRoute(fmt.Sprintf("rt-[%s]", ip), FabricNodeDn, desc, ipRoutePAttr)

	ipRouteP.Status = "modified"

	err := aciClient.Save(ipRouteP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("ip")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_ip_rs_route_track") {
		_, newRelParam := d.GetChange("relation_ip_rs_route_track")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_ip_rs_route_track") {
		_, newRelParam := d.GetChange("relation_ip_rs_route_track")
		err = aciClient.DeleteRelationipRsRouteTrackFromL3outStaticRoute(ipRouteP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationipRsRouteTrackFromL3outStaticRoute(ipRouteP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ip_rs_route_track")
		d.Partial(false)

	}

	d.SetId(ipRouteP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outStaticRouteRead(d, m)

}

func resourceAciL3outStaticRouteRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ipRouteP, err := getRemoteL3outStaticRoute(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3outStaticRouteAttributes(ipRouteP, d)

	ipRsRouteTrackData, err := aciClient.ReadRelationipRsRouteTrackFromL3outStaticRoute(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation ipRsRouteTrack %v", err)
		d.Set("relation_ip_rs_route_track", "")

	} else {
		if _, ok := d.GetOk("relation_ip_rs_route_track"); ok {
			tfName := d.Get("relation_ip_rs_route_track").(string)
			if tfName != ipRsRouteTrackData {
				d.Set("relation_ip_rs_route_track", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outStaticRouteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ipRouteP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

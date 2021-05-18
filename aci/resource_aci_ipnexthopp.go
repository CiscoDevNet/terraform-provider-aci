package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciL3outStaticRouteNextHop() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3outStaticRouteNextHopCreate,
		Update: resourceAciL3outStaticRouteNextHopUpdate,
		Read:   resourceAciL3outStaticRouteNextHopRead,
		Delete: resourceAciL3outStaticRouteNextHopDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outStaticRouteNextHopImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"static_route_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nh_addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
				}, false),
			},

			"nexthop_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"prefix",
				}, false),
			},

			"relation_ip_rs_nexthop_route_track": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_ip_rs_nh_track_member": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteL3outStaticRouteNextHop(client *client.Client, dn string) (*models.L3outStaticRouteNextHop, error) {
	ipNexthopPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	ipNexthopP := models.L3outStaticRouteNextHopFromContainer(ipNexthopPCont)

	if ipNexthopP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outStaticRouteNextHop %s not found", ipNexthopP.DistinguishedName)
	}

	return ipNexthopP, nil
}

func setL3outStaticRouteNextHopAttributes(ipNexthopP *models.L3outStaticRouteNextHop, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(ipNexthopP.DistinguishedName)
	d.Set("description", ipNexthopP.Description)
	dn := d.Id()
	if dn != ipNexthopP.DistinguishedName {
		d.Set("static_route_dn", "")
	}
	ipNexthopPMap, _ := ipNexthopP.ToMap()

	d.Set("nh_addr", ipNexthopPMap["nhAddr"])

	d.Set("annotation", ipNexthopPMap["annotation"])
	d.Set("name_alias", ipNexthopPMap["nameAlias"])
	d.Set("pref", ipNexthopPMap["pref"])
	d.Set("nexthop_profile_type", ipNexthopPMap["type"])
	return d
}

func resourceAciL3outStaticRouteNextHopImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ipNexthopP, err := getRemoteL3outStaticRouteNextHop(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3outStaticRouteNextHopAttributes(ipNexthopP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outStaticRouteNextHopCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outStaticRouteNextHop: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	nhAddr := d.Get("nh_addr").(string)

	StaticRouteDn := d.Get("static_route_dn").(string)

	ipNexthopPAttr := models.L3outStaticRouteNextHopAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ipNexthopPAttr.Annotation = Annotation.(string)
	} else {
		ipNexthopPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ipNexthopPAttr.NameAlias = NameAlias.(string)
	}
	if Pref, ok := d.GetOk("pref"); ok {
		ipNexthopPAttr.Pref = Pref.(string)
	}
	if NexthopProfile_type, ok := d.GetOk("nexthop_profile_type"); ok {
		ipNexthopPAttr.NexthopProfile_type = NexthopProfile_type.(string)
	}
	ipNexthopP := models.NewL3outStaticRouteNextHop(fmt.Sprintf("nh-[%s]", nhAddr), StaticRouteDn, desc, ipNexthopPAttr)

	err := aciClient.Save(ipNexthopP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("nh_addr")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationToipRsNexthopRouteTrack, ok := d.GetOk("relation_ip_rs_nexthop_route_track"); ok {
		relationParam := relationToipRsNexthopRouteTrack.(string)
		checkDns = append(checkDns, relationParam)

	}
	if relationToipRsNHTrackMember, ok := d.GetOk("relation_ip_rs_nh_track_member"); ok {
		relationParam := relationToipRsNHTrackMember.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationToipRsNexthopRouteTrack, ok := d.GetOk("relation_ip_rs_nexthop_route_track"); ok {
		relationParam := relationToipRsNexthopRouteTrack.(string)
		err = aciClient.CreateRelationipRsNexthopRouteTrackFromL3outStaticRouteNextHop(ipNexthopP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ip_rs_nexthop_route_track")
		d.Partial(false)

	}
	if relationToipRsNHTrackMember, ok := d.GetOk("relation_ip_rs_nh_track_member"); ok {
		relationParam := relationToipRsNHTrackMember.(string)
		err = aciClient.CreateRelationipRsNHTrackMemberFromL3outStaticRouteNextHop(ipNexthopP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ip_rs_nh_track_member")
		d.Partial(false)

	}

	d.SetId(ipNexthopP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outStaticRouteNextHopRead(d, m)
}

func resourceAciL3outStaticRouteNextHopUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outStaticRouteNextHop: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	nhAddr := d.Get("nh_addr").(string)

	StaticRouteDn := d.Get("static_route_dn").(string)

	ipNexthopPAttr := models.L3outStaticRouteNextHopAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ipNexthopPAttr.Annotation = Annotation.(string)
	} else {
		ipNexthopPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ipNexthopPAttr.NameAlias = NameAlias.(string)
	}
	if Pref, ok := d.GetOk("pref"); ok {
		ipNexthopPAttr.Pref = Pref.(string)
	}
	if NexthopProfile_type, ok := d.GetOk("nexthop_profile_type"); ok {
		ipNexthopPAttr.NexthopProfile_type = NexthopProfile_type.(string)
	}
	ipNexthopP := models.NewL3outStaticRouteNextHop(fmt.Sprintf("nh-[%s]", nhAddr), StaticRouteDn, desc, ipNexthopPAttr)

	ipNexthopP.Status = "modified"

	err := aciClient.Save(ipNexthopP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("nh_addr")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_ip_rs_nexthop_route_track") {
		_, newRelParam := d.GetChange("relation_ip_rs_nexthop_route_track")
		checkDns = append(checkDns, newRelParam.(string))

	}
	if d.HasChange("relation_ip_rs_nh_track_member") {
		_, newRelParam := d.GetChange("relation_ip_rs_nh_track_member")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_ip_rs_nexthop_route_track") {
		_, newRelParam := d.GetChange("relation_ip_rs_nexthop_route_track")
		err = aciClient.DeleteRelationipRsNexthopRouteTrackFromL3outStaticRouteNextHop(ipNexthopP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationipRsNexthopRouteTrackFromL3outStaticRouteNextHop(ipNexthopP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ip_rs_nexthop_route_track")
		d.Partial(false)

	}
	if d.HasChange("relation_ip_rs_nh_track_member") {
		_, newRelParam := d.GetChange("relation_ip_rs_nh_track_member")
		err = aciClient.DeleteRelationipRsNHTrackMemberFromL3outStaticRouteNextHop(ipNexthopP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationipRsNHTrackMemberFromL3outStaticRouteNextHop(ipNexthopP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ip_rs_nh_track_member")
		d.Partial(false)

	}

	d.SetId(ipNexthopP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outStaticRouteNextHopRead(d, m)

}

func resourceAciL3outStaticRouteNextHopRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ipNexthopP, err := getRemoteL3outStaticRouteNextHop(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3outStaticRouteNextHopAttributes(ipNexthopP, d)

	ipRsNexthopRouteTrackData, err := aciClient.ReadRelationipRsNexthopRouteTrackFromL3outStaticRouteNextHop(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation ipRsNexthopRouteTrack %v", err)
		d.Set("relation_ip_rs_nexthop_route_track", "")

	} else {
		if _, ok := d.GetOk("relation_ip_rs_nexthop_route_track"); ok {
			tfName := d.Get("relation_ip_rs_nexthop_route_track").(string)
			if tfName != ipRsNexthopRouteTrackData {
				d.Set("relation_ip_rs_nexthop_route_track", "")
			}
		}

	}

	ipRsNHTrackMemberData, err := aciClient.ReadRelationipRsNHTrackMemberFromL3outStaticRouteNextHop(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation ipRsNHTrackMember %v", err)
		d.Set("relation_ip_rs_nh_track_member", "")

	} else {
		if _, ok := d.GetOk("relation_ip_rs_nh_track_member"); ok {
			tfName := d.Get("relation_ip_rs_nh_track_member").(string)
			if tfName != ipRsNHTrackMemberData {
				d.Set("relation_ip_rs_nh_track_member", "")
			}
		}

	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outStaticRouteNextHopDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ipNexthopP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

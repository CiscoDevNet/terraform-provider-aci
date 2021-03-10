package aci


import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciVMMCredential() *schema.Resource {
	return &schema.Resource{

		Read:   dataSourceAciVMMCredentialRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},


			"name": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
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



			"pwd": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},



			"usr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},


            }),
    }
}



func dataSourceAciVMMCredentialRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)


    rn := fmt.Sprintf("usracc-%s",name)
    VMMDomainDn := d.Get("vmm_domain_dn").(string)

    dn := fmt.Sprintf("%s/%s",VMMDomainDn,rn)



	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)

	if err != nil {
		return err
	}
	setVMMCredentialAttributes(vmmUsrAccP, d)
	return nil
}

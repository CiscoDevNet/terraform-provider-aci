package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciUserManagement() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciUserManagementRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"pwd_strength_check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"change_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"change_during_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"change_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiration_warn_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"history_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"no_change_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_duration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_login_block": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_failed_attempts": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_failed_attempts_window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maximum_validity_period": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"session_record_flags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ui_idle_timeout_seconds": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"webtoken_timeout_seconds": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciUserManagementRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("userext")
	dn := fmt.Sprintf("uni/%s", rn)
	aaaUserEp, err := getRemoteUserManagement(aciClient, dn)
	if err != nil {
		return err
	}

	_, err = aciClient.Get(dn + "/pwdprofile")
	if err == nil {
		aaaPwdProfileDn := dn + "/pwdprofile"
		aaaPwdProfile, err := getRemotePasswordChangeExpirationPolicy(aciClient, aaaPwdProfileDn)
		if err != nil {
			return err
		}
		_, err = setPasswordChangeExpirationPolicyAttributes(aaaPwdProfile, d)
		if err != nil {
			return nil
		}
	}

	_, err = aciClient.Get(dn + "/blockloginp")
	if err == nil {
		aaaBlockLoginProfileDn := dn + "/blockloginp"
		aaaBlockLoginProfile, err := getRemoteBlockUserLoginsPolicy(aciClient, aaaBlockLoginProfileDn)
		if err != nil {
			return err
		}
		_, err = setBlockUserLoginsPolicyAttributes(aaaBlockLoginProfile, d)
		if err != nil {
			return nil
		}
	}

	_, err = aciClient.Get(dn + "/pkiext/webtokendata")
	if err == nil {
		pkiWebTokenDn := dn + "/pkiext/webtokendata"
		pkiWebTokenData, err := getRemoteWebTokenData(aciClient, pkiWebTokenDn)
		if err != nil {
			return err
		}
		_, err = setWebTokenDataAttributes(pkiWebTokenData, d)
		if err != nil {
			return nil
		}
	}

	d.SetId(dn)
	setUserManagementAttributes(aaaUserEp, d)
	return nil
}

resource "aci_end_point_retention_policy" "test_ret_policy" {
    tenant_dn = aci_tenant.test_tenant.id
    name = "tf_test"
    description 		= "From terraform"
	annotation          = "tag_ret_pol"
	bounce_age_intvl    = "630"
	bounce_trig         = "protocol"
	hold_intvl          = "6"
	local_ep_age_intvl  = "900"
	move_freq           = "256"
	name_alias          = "alias_demo"
	remote_ep_age_intvl = "300"
}

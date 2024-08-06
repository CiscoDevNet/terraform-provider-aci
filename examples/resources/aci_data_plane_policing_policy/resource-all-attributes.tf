
resource "aci_data_plane_policing_policy" "full_example_tenant" {
  parent_dn            = aci_tenant.example.id
  admin_state          = "disabled"
  annotation           = "annotation"
  excessive_burst      = "unspecified"
  excessive_burst_unit = "giga"
  burst                = "unspecified"
  burst_unit           = "giga"
  conform_action       = "drop"
  conform_mark_cos     = "unspecified"
  conform_mark_dscp    = "unspecified"
  description          = "description_1"
  exceed_action        = "drop"
  exceed_mark_cos      = "unspecified"
  exceed_mark_dscp     = "unspecified"
  mode                 = "bit"
  name                 = "test_name"
  name_alias           = "name_alias_1"
  owner_key            = "owner_key_1"
  owner_tag            = "owner_tag_1"
  peak_rate            = "0"
  peak_rate_unit       = "giga"
  rate                 = "0"
  rate_unit            = "giga"
  sharing_mode         = "dedicated"
  type                 = "1R2C"
  violate_action       = "drop"
  violate_mark_cos     = "unspecified"
  violate_mark_dscp    = "unspecified"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}

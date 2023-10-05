
resource "aci_pim_route_map_entry" "example" {
  parent_dn   = aci_pim_route_map_policy.example.id
  action      = "deny"
  annotation  = "annotation"
  description = "description"
  grp         = "0.0.0.0"
  name        = "name"
  name_alias  = "name_alias"
  order       = "1"
  rp          = "0.0.0.0"
  src         = "1.1.1.1/30"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}


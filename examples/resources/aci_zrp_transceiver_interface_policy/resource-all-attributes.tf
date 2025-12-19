
resource "aci_zrp_transceiver_interface_policy" "full_example" {
  admin_state                  = "enabled"
  annotation                   = "annotation"
  maximum_chromatic_dispersion = "-50989"
  minimum_chromatic_dispersion = "-51999"
  dac_rate                     = "1x1.25"
  description                  = "description_1"
  dwdm_carrier_grid            = "50GHzITUchannel"
  fec_mode                     = "oFEC"
  frequency_100_mhz            = "1923504"
  frequency_50_ghz             = "19369"
  itu_channel_50_ghz           = "64"
  modulation                   = "16QAM"
  muxponder_mode               = "1x400"
  name                         = "test_name"
  name_alias                   = "name_alias_1"
  owner_key                    = "owner_key_1"
  owner_tag                    = "owner_tag_1"
  transmit_power               = "-120"
  wavelength_50_ghz            = "1545733"
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

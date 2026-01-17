
resource "aci_zr_transceiver_interface_policy" "full_example" {
  admin_state                  = "enabled"
  annotation                   = "annotation"
  maximum_chromatic_dispersion = "-1390"
  minimum_chromatic_dispersion = "-2391"
  dac_rate                     = "1x1"
  description                  = "description_1"
  dwdm_carrier_grid            = "50GHzFrequency"
  fec_mode                     = "cFEC"
  frequency_100_mhz            = "1912504"
  frequency_50_ghz             = "19415"
  itu_channel_50_ghz           = "23"
  modulation                   = "16QAM"
  muxponder_mode               = "1x400"
  name                         = "test_name"
  name_alias                   = "name_alias_1"
  owner_key                    = "owner_key_1"
  owner_tag                    = "owner_tag_1"
  transmit_power               = "20"
  wavelength_50_ghz            = "1563742"
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

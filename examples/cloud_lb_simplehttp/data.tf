
resource "time_sleep" "wait_for_cloudLB" {
  depends_on      = [aci_rest.appliedServiceGraph]
  create_duration = "10s"
}

data "aws_lb" "wwwalb" {
  depends_on = [time_sleep.wait_for_cloudLB]
  name       = "${var.name}-${join("", regex("ctxprofile-(.*?)/", var.subnet_a_dn))}"
}

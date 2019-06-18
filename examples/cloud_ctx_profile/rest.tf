resource "aci_rest" "rest_l3_ext_out" {
  path       = "api/node/mo/tn-cl-workshop/ctxprofile-cloudcontext-test.json"
  class_name = "cloudCtxProfile"

  content = {
    "name" = "cloudcontext-test"
    "children" =  [
        {"cloudRsToCtx" = {
            "attributes"= {
                "status"= "created,modified",
                "tnFvCtxName"= "user-vrf-demo1"
              }
          }
        }
    ]
}


}
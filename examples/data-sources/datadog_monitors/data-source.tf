data "datadog_monitors" "default" {
  
  filter {
      name = "type"
      values = [ "query alert" ]
  }

}

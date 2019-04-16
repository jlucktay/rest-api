locals {
  gcp_region = "${substr(var.gcp_zone, 0, length(var.gcp_zone)-2)}"
}

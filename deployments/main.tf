provider "google" {
  project = "${var.gcp_project_id}"
  region  = "${local.gcp_region}"
  zone    = "${var.gcp_zone}"
}

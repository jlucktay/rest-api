variable "gcp_project_id" {
  description = "Should contain the GCP project ID (listed in many places, like [the dashboard](https://console.cloud.google.com/home/dashboard)). Example: `export TF_VAR_gcp_project_id=single-axis-237719`"
  type        = "string"
}

variable "gcp_zone" {
  description = "Should contain a GCP zone, from which the region will also be derived. Can be looked up on [this page](https://cloud.google.com/compute/docs/regions-zones/) or listed with `gcloud compute zones list`. Example: `export TF_VAR_gcp_zone=us-east1-b`"
  type        = "string"
}

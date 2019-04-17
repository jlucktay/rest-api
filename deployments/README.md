# Deploying to GCP

1. Create/identify a project on GCP to which you will deploy resources.
1. Create a service account on [this page](https://console.cloud.google.com/apis/credentials/serviceaccountkey) with
the following permissions:
    - TODO add specific list of permissions
1. Store the resulting JSON credentials file somewhere secure.
    - **Do not check this file into source control!**
1. Set up the following environment variables:
    1. `GOOGLE_CLOUD_KEYFILE_JSON` should contain the location of the credentials file.
    Example: `export GOOGLE_CLOUD_KEYFILE_JSON=/path/to/credentials.json`
    1. `TF_VAR_gcp_project_id` should contain the GCP project ID (listed in many places, like
    [the dashboard](https://console.cloud.google.com/home/dashboard)).
    Example: `export TF_VAR_gcp_project_id=single-axis-237719`
    1. `TF_VAR_gcp_zone` should contain a GCP zone, from which the region will also be derived. Can be looked up on
    [this page](https://cloud.google.com/compute/docs/regions-zones/) or listed with `gcloud compute zones list`.
    Example: `export TF_VAR_gcp_zone=europe-west2-b`
There is an [example Bash script](env.example.sh) to help set these up.

Make note of the following:

- the example Bash script **must be dot sourced** rather than executed directly
- the Bash environment variable names are case-sensitive

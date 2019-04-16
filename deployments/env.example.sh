#!/usr/bin/env bash
if [[ "${BASH_SOURCE[0]}" = "${0}" ]]; then
    echo "This script '${BASH_SOURCE[0]}' must be sourced, like so:"
    echo "    $(tput setab 7 ; tput setaf 0). ${BASH_SOURCE[0]}$(tput sgr0)"
    exit 1
fi

PreviousIFS=$IFS
IFS=$'\n\t'

echo -n "Exporting GCP variables for Terraform... "

export GOOGLE_CLOUD_KEYFILE_JSON=/path/to/credentials.json
export TF_VAR_gcp_project_id=single-axis-237719
export TF_VAR_gcp_zone=us-east1-b

IFS=$PreviousIFS

echo "Done!"

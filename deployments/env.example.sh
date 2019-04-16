#!/usr/bin/env bash
# Thank you: https://github.com/anordal/shellharden/blob/master/how_to_do_things_safely_in_bash.md#how-to-begin-a-bash-script
if test "$BASH" = "" || "$BASH" -uc "a=();true \"\${a[@]}\"" 2>/dev/null; then
    # Bash 4.4, Zsh
    set -euo pipefail
else
    # Bash 4.3 and older chokes on empty arrays with set -u.
    set -eo pipefail
fi

shopt -s nullglob globstar
PreviousIFS=$IFS
IFS=$'\n\t'

if [[ "${BASH_SOURCE[0]}" = "${0}" ]]; then
    echo "This script '${BASH_SOURCE[0]}' must be sourced, like so:"
    echo "    $(tput setab 7 ; tput setaf 0). ${BASH_SOURCE[0]}$(tput sgr0)"
    exit 1
fi

echo -n "Exporting GCP variables for Terraform... "

export GOOGLE_CLOUD_KEYFILE_JSON=/path/to/credentials.json
export TF_VAR_gcp_project_id=single-axis-237719
export TF_VAR_gcp_zone=us-east1-b

IFS=$PreviousIFS

echo "Done!"

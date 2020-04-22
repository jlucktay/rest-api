#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

if ! command -v jq > /dev/null; then
  echo "'jq' not found! Please install: https://stedolan.github.io/jq/download/"
  exit 1
fi

secrets_file="$(cd "$(dirname "${BASH_SOURCE[0]}")" > /dev/null 2>&1 && pwd)/secrets.trello.json"

if ! [[ -r $secrets_file ]]; then
  echo "The Trello secrets file '$secrets_file' could not be read. See the example JSON file alongside this script."
  exit 1
fi

board=$(jq --exit-status --raw-output '.board' "$secrets_file")
finished_list_id=$(jq --exit-status --raw-output '.finishedListId' "$secrets_file")
key=$(jq --exit-status --raw-output '.key' "$secrets_file")
token=$(jq --exit-status --raw-output '.token' "$secrets_file")

trello_get_url="https://api.trello.com/1/boards/$board/cards/?key=$key&token=$token"

curl --request GET --silent "$trello_get_url" \
  | jq --exit-status --raw-output '.[] | select( .idList != "'"$finished_list_id"'" ) | .name' \
  | sort -f

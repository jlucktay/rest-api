#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

if ! command -v jq > /dev/null; then
  echo "'jq' not found! Please install: https://stedolan.github.io/jq/download/"
  exit 1
fi

SecretsFile="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )/secrets.trello.json"

if ! [[ -f $SecretsFile ]]; then
  echo "The Trello secrets file '$SecretsFile' does not exist. See the example JSON file alongside this script."
  exit 1
fi

Board=$(jq -r '.board' "$SecretsFile")
FinishedListId=$(jq -r '.finishedListId' "$SecretsFile")
Key=$(jq -r '.key' "$SecretsFile")
Token=$(jq -r '.token' "$SecretsFile")

URL="https://api.trello.com/1/boards/$Board/cards/?key=$Key&token=$Token"

curl --silent "$URL" | jq -r '.[] | select( .idList != "'"$FinishedListId"'" ) | .name' | sort -f

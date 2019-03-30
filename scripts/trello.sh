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

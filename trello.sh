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

Board=$(jq -r '.board' secrets.trello.json)
FinishedListId=$(jq -r '.finishedListId' secrets.trello.json)
Key=$(jq -r '.key' secrets.trello.json)
Token=$(jq -r '.token' secrets.trello.json)

URL="https://api.trello.com/1/boards/$Board/cards/?key=$Key&token=$Token"

curl --silent "$URL" | jq -r '.[] | select( .idList != "'"$FinishedListId"'" ) | .name' | sort -f

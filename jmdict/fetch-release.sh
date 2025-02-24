#!/bin/bash

JSON_RESPONSE=$(curl -s https://api.github.com/repos/scriptin/jmdict-simplified/releases/latest)

# Matches any string that starts with "jmdict-eng-common-" followed by any characters. Must end with .zip too.
LATEST_RELEASE=$(echo "$JSON_RESPONSE" | jq '.assets[] | select(.name | test("jmdict-eng-common-.*\\.zip"))')

if [ -z "$LATEST_RELEASE" ]; then
    echo "No matching item found."
else
    # Log the release name
    echo "Matching item name: $LATEST_RELEASE"

    # Download and uncompress the file
    BROWSER_DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | jq -r '.browser_download_url')

    # Follow HTTP redirections (301) https://askubuntu.com/a/1036492
    curl -sLO -H 'Accept-encoding: gzip' "$BROWSER_DOWNLOAD_URL"
fi

read -p "Press any key to continue" x
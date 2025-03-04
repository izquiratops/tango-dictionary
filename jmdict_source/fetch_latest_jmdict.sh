#!/bin/bash

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "darwin"* ]]; then
    OS="UNIX"
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    OS="Windows"
else
    echo "Unsupported operating system"
    exit 1
fi

JSON_RESPONSE=$(curl -s https://api.github.com/repos/scriptin/jmdict-simplified/releases/latest)

# Matches any string that starts with "jmdict-eng-common-" followed by any characters. Must end with .zip too.
LATEST_RELEASE=$(echo "$JSON_RESPONSE" | jq '.assets[] | select(.name | test("jmdict-eng-common-.*\\.zip"))')

if [ -z "$LATEST_RELEASE" ]; then
    echo "No matching item found."
else
    echo "Matching item name: $LATEST_RELEASE"
    BROWSER_DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | jq -r '.browser_download_url')

    # Follow HTTP redirections (301) https://askubuntu.com/a/1036492
    TMP_FILE="jmdict.zip"
    curl -sL -H 'Accept-encoding: gzip' "$BROWSER_DOWNLOAD_URL" -o "$TMP_FILE"

    if [ "$OS" = "UNIX" ]; then
        echo "Extracting with unzip..."
        unzip -o "$TMP_FILE"
    else
        # Trying to unzip it with PowerShell using tar -xf doesn't work and I don't know why!
        echo "Extracting with PowerShell..."
        powershell -command "Expand-Archive -Path "$TMP_FILE" -DestinationPath . -Force"
    fi

    # Zip is not needed anymore 💫
    rm jmdict.zip
fi

read -p "Press any key to continue" x
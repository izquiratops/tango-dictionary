#!/bin/bash

SOURCE_FOLDER=""
REMOTE_USER=""
REMOTE_HOST=""
REMOTE_PATH=""

usage() {
    echo "Usage: $0 -f <source_folder> -u <remote_user> -h <remote_host> -p <remote_path>"
    echo "Example: $0 -f /path/to/folder -u user -h domain.com -p /root/destination"
    exit 1
}

while getopts "f:u:h:p:" opt; do
    case $opt in
        f) SOURCE_FOLDER="$OPTARG";;
        u) REMOTE_USER="$OPTARG";;
        h) REMOTE_HOST="$OPTARG";;
        p) REMOTE_PATH="$OPTARG";;
        ?) usage;;
    esac
done

if [ -z "$SOURCE_FOLDER" ] || [ -z "$REMOTE_USER" ] || [ -z "$REMOTE_HOST" ] || [ -z "$REMOTE_PATH" ]; then
    echo "Error: Missing required parameters!"
    usage
fi

if [ ! -d "$SOURCE_FOLDER" ]; then
    echo "Error: Source folder '$SOURCE_FOLDER' not found!"
    exit 1
fi

SOURCE_FOLDER=${SOURCE_FOLDER%/}

echo "Sending folder $SOURCE_FOLDER to $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH"

if scp -r "$SOURCE_FOLDER" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH"; then
    echo "Folder transfer completed successfully!"
else
    echo "Error: Folder transfer failed!"
    exit 1
fi
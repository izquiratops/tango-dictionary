#!/bin/bash

# Check if directory path is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <directory_path>"
    read -p "Press any key to exit" x
    exit 1
fi

DIR_PATH="$1"
FILE_COUNT=$(find "$DIR_PATH" -type f -not -name "*.gz" | wc -l)

if [ $FILE_COUNT -eq 0 ]; then
    echo "No files found to compress in '$DIR_PATH'"
    read -p "Press any key to exit" x
    exit 0
fi

find "$DIR_PATH" -type f -not -name "*.gz" | while read file; do
    echo "Compressing: $file"
    gzip -k -f "$file"
    
    # Check if compression was successful
    if [ $? -eq 0 ]; then
        echo "Successfully compressed: $file"
    else
        echo "Failed to compress: $file"
    fi
done

read -p "Press any key to exit" x
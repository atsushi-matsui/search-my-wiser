#!/bin/bash

# URL of the Wikipedia dump
URL="https://dumps.wikimedia.org/jawiki/latest/jawiki-latest-pages-articles.xml.bz2"

# カレントディレクトリを取得
CURRENT=$(cd $(dirname $0); pwd)
# 一つ上のディレクトリを取得
CURRENT_PATH=$(dirname "$CURRENT")
# Output directory
OUTPUT_DIR="${CURRENT_PATH}/files"

# Output file path
OUTPUT_FILE="${OUTPUT_DIR}/jawiki-latest-pages-articles.xml.bz2"
OUTPUT_EXTRACTED_FILE="${OUTPUT_DIR}/jawiki-latest-pages-articles.xml"

# Create the output directory if it doesn't exist
mkdir -p ${OUTPUT_DIR}

# Check if the file already exists
if [ -f "${OUTPUT_EXTRACTED_FILE}" ]; then
    echo "${OUTPUT_EXTRACTED_FILE} already exists. Skipping download."
    exit 0
else
    # Download the file
    echo "Downloading ${OUTPUT_FILE} from ${URL}..."
    curl -o ${OUTPUT_FILE} ${URL}

    # Check if the download was successful
    if [ $? -eq 0 ]; then
        echo "Download completed successfully."
    else
        echo "Download failed. Exiting."
        exit 1
    fi
fi

# Extract the .bz2 file using tar
echo "Extracting ${OUTPUT_FILE}..."
tar -xjf ${OUTPUT_FILE} -C ${OUTPUT_DIR}

# Check if the extraction was successful
if [ $? -eq 0 ]; then
    echo "Extraction completed successfully."
else
    echo "Extraction failed. Exiting."
    exit 1
fi

echo "Done."

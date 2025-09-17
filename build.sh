#!/bin/bash
set -e

APP_NAME="ffcutter"

OUTPUT_DIR="./rpmdude_build/SOURCES/"


echo "Building app..."
CGO_ENABLED=0 go build -ldflags="-s -w" -o "$APP_NAME"
cp "$APP_NAME" "$OUTPUT_DIR/"

echo "Build done: $OUTPUT_DIR/$APP_NAME"

rpmdude build

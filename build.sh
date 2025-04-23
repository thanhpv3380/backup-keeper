#!/bin/bash

APP_NAME="backup-keeper"
BUILD_DIR="build"
MAIN_FILE="./cmd/main.go"

mkdir -p "$BUILD_DIR"

echo "🔨 Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/${APP_NAME}-linux" "$MAIN_FILE"

echo "✅ Done. Built files in ./$BUILD_DIR/"

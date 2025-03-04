#!/bin/bash
set -e

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please install Protocol Buffers compiler (protoc) from https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null || ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Error: protoc-gen-go or protoc-gen-go-grpc is not installed"
    echo "Please install them with:"
    echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Directory containing .proto files
PROTO_DIR="./pkg/proto"

# Output directory for generated Go files
GO_OUT_DIR="./pkg/proto"

# Create output directory if it doesn't exist
mkdir -p "$GO_OUT_DIR"

echo "Generating Go code from Protocol Buffers definitions..."

# Generate Go code from .proto files
protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$GO_OUT_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$GO_OUT_DIR" --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR"/*.proto

echo "Done!" 
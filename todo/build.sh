#!/bin/bash

# Smart Todo CLI Build Script
# Builds cross-platform binaries for distribution

set -e

VERSION=${1:-"v1.0.0"}
BINARY_NAME="todo"
BUILD_DIR="dist"

echo "🚀 Building Smart Todo CLI ${VERSION}"

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build targets
TARGETS=(
    "linux/amd64"
    "linux/arm64" 
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

echo "📦 Building for multiple platforms..."

for target in "${TARGETS[@]}"; do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    
    output_name="${BINARY_NAME}-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    
    echo "Building ${output_name}..."
    
    GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o ${BUILD_DIR}/${output_name} \
        .
done

echo "✅ Build complete! Binaries available in ${BUILD_DIR}/"
echo ""
echo "📁 Generated files:"
ls -la ${BUILD_DIR}/

echo ""
echo "🎯 Quick deployment commands:"
echo "  Linux/macOS: curl -L <url>/${BINARY_NAME}-linux-amd64 -o ${BINARY_NAME} && chmod +x ${BINARY_NAME}"
echo "  Windows: Download ${BINARY_NAME}-windows-amd64.exe"
echo ""
echo "💡 Don't forget to:"
echo "  1. Test binaries on target platforms"
echo "  2. Create GitHub release with ${VERSION}"
echo "  3. Upload binaries as release assets"
echo "  4. Update installation instructions"

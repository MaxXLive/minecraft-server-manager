#!/bin/bash

cd /Users/MERMACK/Projects/minecraft-server-manager
#!/bin/bash

# Output directory
OUTPUT_DIR="build"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# List of target OS and architectures
TARGETS=(
    "windows amd64"
    "windows arm64"
    "darwin amd64"
    "darwin arm64"
    "linux amd64"
    "linux arm64"
)

# Default values for OS and ARCH
BUILD_OS=""
BUILD_ARCH=""

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        --os)
            BUILD_OS="$2"
            shift 2
            ;;
        --arch)
            BUILD_ARCH="$2"
            shift 2
            ;;
        *)
            echo -e "\033[31mUnknown option: $1\033[0m"
            exit 1
            ;;
    esac
done

# Function to build for a specific OS and ARCH
build_target() {
    local OS=$1
    local ARCH=$2
    local OUTPUT_FILE="$OUTPUT_DIR/minecraft-server-manager_${OS}_${ARCH}"

    # Windows executables need .exe extension
    if [ "$OS" == "windows" ]; then
        OUTPUT_FILE+=".exe"
    fi

    echo "Building for $OS $ARCH..."

    # Build the binary
    GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT_FILE" main.go

    if [ $? -eq 0 ]; then
        echo "✅  Successfully built: $OUTPUT_FILE"
    else
        echo -e "\033[31m❌ Failed to build for $OS $ARCH\033[0m"
    fi
}

# Loop through targets and build
for target in "${TARGETS[@]}"; do
    OS=$(echo "$target" | awk '{print $1}')
    ARCH=$(echo "$target" | awk '{print $2}')

    # If specific OS or ARCH is set, filter targets
    if [[ -n "$BUILD_OS" && "$BUILD_OS" != "$OS" ]]; then
        continue
    fi
    if [[ -n "$BUILD_ARCH" && "$BUILD_ARCH" != "$ARCH" ]]; then
        continue
    fi

    build_target "$OS" "$ARCH"
done

echo "🎉 Build process complete. Binaries are in the /$OUTPUT_DIR directory."

# Detect current OS and architecture
CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
CURRENT_ARCH=$(uname -m)

# Map architecture names
if [ "$CURRENT_ARCH" == "x86_64" ]; then
    CURRENT_ARCH="amd64"
elif [ "$CURRENT_ARCH" == "aarch64" ] || [ "$CURRENT_ARCH" == "arm64" ]; then
    CURRENT_ARCH="arm64"
fi

# Build the path to the binary
BINARY_PATH="$(pwd)/$OUTPUT_DIR/minecraft-server-manager_${CURRENT_OS}_${CURRENT_ARCH}"

if [ -f "$BINARY_PATH" ]; then
    echo ""
    echo "To use 'msm' in this session, run:"
    echo "  source <(echo 'msm() { $BINARY_PATH \"\$@\"; }')"
    echo ""
    echo "Or add this to your ~/.bashrc or ~/.zshrc for permanent use:"
    echo "  alias msm='$BINARY_PATH'"
fi

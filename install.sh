#!/bin/bash
# Quick installation script for imgconvert

set -e

echo "imgconvert - Installation Script"
echo "================================"
echo ""

# Check if running from project directory
if [ ! -f "cmd/imgconvert/main.go" ]; then
    echo "Error: Please run this from the imgconvert project directory"
    exit 1
fi

# Build
echo "Building imgconvert..."
export PATH="$HOME/gosdk/bin:$PATH"
go build -o imgconvert -ldflags "-s -w" ./cmd/imgconvert

echo "✓ Build complete"
echo ""

# Offer installation options
echo "Choose installation method:"
echo ""
echo "1. Install globally (recommended) - requires sudo"
echo "   → /usr/local/bin/imgconvert"
echo ""
echo "2. Add to PATH - no sudo required"
echo "   → Update ~/.zshrc or ~/.bashrc"
echo ""
echo "3. Create alias - no sudo required"
echo "   → Add alias to shell config"
echo ""
echo "4. Skip - just build (use from this directory)"
echo ""

read -p "Enter choice (1-4): " choice

case $choice in
    1)
        echo ""
        echo "Installing to /usr/local/bin/ ..."
        sudo cp imgconvert /usr/local/bin/
        echo "✓ Installed successfully!"
        echo ""
        echo "Test it: imgconvert --version"
        ;;
    2)
        INSTALL_DIR=$(pwd)
        SHELL_CONFIG=""
        
        if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
            SHELL_CONFIG="$HOME/.zshrc"
        elif [ -f "$HOME/.bashrc" ]; then
            SHELL_CONFIG="$HOME/.bashrc"
        else
            SHELL_CONFIG="$HOME/.profile"
        fi
        
        echo "" >> "$SHELL_CONFIG"
        echo "# imgconvert" >> "$SHELL_CONFIG"
        echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> "$SHELL_CONFIG"
        
        echo "✓ Added to PATH in $SHELL_CONFIG"
        echo ""
        echo "Run: source $SHELL_CONFIG"
        echo "Then: imgconvert --version"
        ;;
    3)
        INSTALL_DIR=$(pwd)
        SHELL_CONFIG=""
        
        if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
            SHELL_CONFIG="$HOME/.zshrc"
        elif [ -f "$HOME/.bashrc" ]; then
            SHELL_CONFIG="$HOME/.bashrc"
        else
            SHELL_CONFIG="$HOME/.profile"
        fi
        
        echo "" >> "$SHELL_CONFIG"
        echo "# imgconvert" >> "$SHELL_CONFIG"
        echo "alias imgconvert=\"$INSTALL_DIR/imgconvert\"" >> "$SHELL_CONFIG"
        
        echo "✓ Added alias to $SHELL_CONFIG"
        echo ""
        echo "Run: source $SHELL_CONFIG"
        echo "Then: imgconvert --version"
        ;;
    4)
        echo ""
        echo "Build complete. Use from this directory:"
        echo "  ./imgconvert --help"
        ;;
    *)
        echo "Invalid choice. Build complete, use: ./imgconvert"
        ;;
esac

echo ""
echo "Quick start:"
echo "  imgconvert photo.jpg"
echo "  imgconvert -r /path/to/photos"
echo "  imgconvert --help"

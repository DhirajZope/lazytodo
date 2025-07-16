#!/bin/bash

# LazyTodo Linux Installation Script
# This script downloads and installs LazyTodo for Linux

set -e

# Default configuration
VERSION="${1:-latest}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
CONFIG_DIR="${HOME}/.config/lazytodo"
REPO_URL="https://github.com/yourusername/lazytodo"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Emoji support check
if locale | grep -q "UTF-8"; then
    CHECKMARK="✅"
    CROSS="❌"
    ROCKET="🚀"
    PACKAGE="📦"
    FOLDER="📁"
    GEAR="🔧"
    BOOK="📚"
    LINK="🔗"
    TEST="🧪"
    WARNING="⚠️"
    DESKTOP="🖥️"
else
    CHECKMARK="[OK]"
    CROSS="[ERROR]"
    ROCKET="[INSTALL]"
    PACKAGE="[DOWNLOAD]"
    FOLDER="[DIR]"
    GEAR="[CONFIG]"
    BOOK="[DOCS]"
    LINK="[LINK]"
    TEST="[TEST]"
    WARNING="[WARNING]"
    DESKTOP="[DESKTOP]"
fi

# Helper functions
print_info() {
    echo -e "${CYAN}$1${NC}"
}

print_success() {
    echo -e "${GREEN}$CHECKMARK $1${NC}"
}

print_error() {
    echo -e "${RED}$CROSS $1${NC}"
    exit 1
}

print_warning() {
    echo -e "${YELLOW}$WARNING $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect architecture
detect_arch() {
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l|armv6l)
            echo "arm"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            ;;
    esac
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "darwin"
    else
        print_error "Unsupported OS: $OSTYPE"
    fi
}

# Download and extract
download_release() {
    local os arch download_url asset_name temp_file
    os=$(detect_os)
    arch=$(detect_arch)
    
    print_info "$PACKAGE Downloading LazyTodo..."
    
    # Get release information
    local api_url="https://api.github.com/repos/yourusername/lazytodo/releases/latest"
    if [[ "$VERSION" != "latest" ]]; then
        api_url="https://api.github.com/repos/yourusername/lazytodo/releases/tags/$VERSION"
    fi
    
    if command_exists curl; then
        local release_info
        release_info=$(curl -s "$api_url")
    elif command_exists wget; then
        local release_info
        release_info=$(wget -qO- "$api_url")
    else
        print_error "Neither curl nor wget found. Please install one of them."
    fi
    
    # Extract download URL for our platform
    asset_name="lazytodo-${os}-${arch}.tar.gz"
    download_url=$(echo "$release_info" | grep -o "\"browser_download_url\": \"[^\"]*${asset_name}\"" | cut -d '"' -f 4)
    
    if [[ -z "$download_url" ]]; then
        print_error "No binary found for ${os}-${arch} in release"
    fi
    
    print_info "$PACKAGE Downloading: $asset_name"
    
    # Download
    temp_file="/tmp/lazytodo-${os}-${arch}.tar.gz"
    if command_exists curl; then
        curl -L -o "$temp_file" "$download_url"
    elif command_exists wget; then
        wget -O "$temp_file" "$download_url"
    fi
    
    print_success "Downloaded successfully"
    echo "$temp_file"
}

# Create directories
create_directories() {
    print_info "$FOLDER Creating directories..."
    
    # Create install directory
    if [[ ! -d "$INSTALL_DIR" ]]; then
        mkdir -p "$INSTALL_DIR"
        print_success "Created: $INSTALL_DIR"
    else
        print_info "Directory exists: $INSTALL_DIR"
    fi
    
    # Create config directory
    if [[ ! -d "$CONFIG_DIR" ]]; then
        mkdir -p "$CONFIG_DIR"
        print_success "Created: $CONFIG_DIR"
    else
        print_info "Config directory exists: $CONFIG_DIR"
    fi
}

# Install binary
install_binary() {
    local temp_file="$1"
    
    print_info "$PACKAGE Extracting and installing..."
    
    # Extract
    tar -xzf "$temp_file" -C "/tmp/"
    
    # Find the binary (it might be in a subdirectory)
    local binary_path
    binary_path=$(find "/tmp" -name "lazytodo" -type f -perm -u+x | head -1)
    
    if [[ -z "$binary_path" ]]; then
        print_error "Binary not found in downloaded archive"
    fi
    
    # Install
    cp "$binary_path" "$INSTALL_DIR/lazytodo"
    chmod +x "$INSTALL_DIR/lazytodo"
    
    # Cleanup
    rm -f "$temp_file"
    rm -rf "$(dirname "$binary_path")" 2>/dev/null || true
    
    print_success "Binary installed to $INSTALL_DIR/lazytodo"
}

# Configure PATH
configure_path() {
    print_info "$GEAR Configuring PATH..."
    
    # Check if directory is already in PATH
    if echo "$PATH" | grep -q "$INSTALL_DIR"; then
        print_info "Install directory already in PATH"
        return
    fi
    
    # Add to shell profile
    local shell_profile=""
    if [[ -n "$BASH_VERSION" ]]; then
        if [[ -f "$HOME/.bashrc" ]]; then
            shell_profile="$HOME/.bashrc"
        elif [[ -f "$HOME/.bash_profile" ]]; then
            shell_profile="$HOME/.bash_profile"
        fi
    elif [[ -n "$ZSH_VERSION" ]]; then
        shell_profile="$HOME/.zshrc"
    elif [[ -f "$HOME/.profile" ]]; then
        shell_profile="$HOME/.profile"
    fi
    
    if [[ -n "$shell_profile" ]]; then
        echo "" >> "$shell_profile"
        echo "# Added by LazyTodo installer" >> "$shell_profile"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$shell_profile"
        print_success "Added $INSTALL_DIR to PATH in $shell_profile"
        print_warning "Please restart your terminal or run: source $shell_profile"
    else
        print_warning "Could not detect shell profile. Please add $INSTALL_DIR to your PATH manually"
    fi
}

# Create desktop entry (Linux only)
create_desktop_entry() {
    if [[ "$(detect_os)" != "linux" ]]; then
        return
    fi
    
    local desktop_dir="$HOME/.local/share/applications"
    if [[ ! -d "$desktop_dir" ]]; then
        mkdir -p "$desktop_dir"
    fi
    
    cat > "$desktop_dir/lazytodo.desktop" << EOF
[Desktop Entry]
Name=LazyTodo
Comment=Beautiful Terminal Todo Manager
Exec=gnome-terminal -- $INSTALL_DIR/lazytodo
Icon=text-editor
Terminal=true
Type=Application
Categories=Office;Productivity;
Keywords=todo;task;productivity;terminal;
EOF
    
    print_success "Desktop entry created"
}

# Main installation process
main() {
    echo
    print_info "$ROCKET LazyTodo Linux Installer"
    print_info "=============================="
    echo
    
    # Check dependencies
    print_info "Checking dependencies..."
    if ! command_exists curl && ! command_exists wget; then
        print_error "Either curl or wget is required for installation"
    fi
    
    if ! command_exists tar; then
        print_error "tar is required for installation"
    fi
    
    print_success "Dependencies check passed"
    
    # Create directories
    create_directories
    
    # Download and install
    local temp_file
    temp_file=$(download_release)
    install_binary "$temp_file"
    
    # Configure environment
    configure_path
    
    # Create desktop entry (optional)
    if [[ "$(detect_os)" == "linux" ]]; then
        read -p "$DESKTOP Create desktop entry? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            create_desktop_entry
        fi
    fi
    
    # Installation complete
    echo
    print_success "Installation Complete!"
    print_info "======================"
    print_info "Installation Directory: $INSTALL_DIR"
    print_info "Executable: $INSTALL_DIR/lazytodo"
    print_info "Config Directory: $CONFIG_DIR"
    echo
    print_info "$BOOK Getting Started:"
    print_info "  • Run: lazytodo (after restarting terminal)"
    print_info "  • Or run directly: $INSTALL_DIR/lazytodo"
    print_info "  • Press 'n' to create your first todo list"
    print_info "  • Press '?' for help"
    echo
    print_info "$LINK Documentation: $REPO_URL"
    echo
    
    # Test installation
    read -p "$TEST Test installation now? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Testing installation..."
        if "$INSTALL_DIR/lazytodo" --version 2>/dev/null; then
            print_success "Installation test successful!"
        else
            print_warning "Could not test installation, but binary is installed"
        fi
    fi
    
    print_success "LazyTodo is ready! Enjoy organizing your tasks! $CHECKMARK"
}

# Show help
show_help() {
    cat << EOF
LazyTodo Linux Installation Script

Usage: $0 [VERSION]

Arguments:
  VERSION     Version to install (default: latest)

Environment Variables:
  INSTALL_DIR Directory to install binary (default: \$HOME/.local/bin)

Examples:
  $0                    # Install latest version
  $0 v1.0.0            # Install specific version
  INSTALL_DIR=/usr/local/bin $0  # Install to custom directory

For more information, visit: $REPO_URL
EOF
}

# Handle arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac 
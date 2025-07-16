#!/bin/bash

# LazyTodo macOS Installation Script
# This script downloads and installs LazyTodo for macOS

set -e

# Default configuration
VERSION="${1:-latest}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
CONFIG_DIR="${HOME}/.config/lazytodo"
REPO_URL="https://github.com/yourusername/lazytodo"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Emojis (macOS Terminal supports these well)
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
APPLE="🍎"

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

# Check if we have sudo privileges for /usr/local/bin
check_sudo() {
    if [[ "$INSTALL_DIR" == "/usr/local/bin" ]] && [[ ! -w "$INSTALL_DIR" ]]; then
        print_warning "Administrator privileges required for installation to $INSTALL_DIR"
        print_info "You can also install to user directory: INSTALL_DIR=\"\$HOME/.local/bin\" $0"
        return 1
    fi
    return 0
}

# Detect architecture for Apple Silicon / Intel
detect_arch() {
    local arch
    arch=$(uname -m)
    case $arch in
        x86_64|amd64)
            echo "amd64"
            ;;
        arm64)
            echo "arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            ;;
    esac
}

# Download and extract
download_release() {
    local arch download_url asset_name temp_file
    arch=$(detect_arch)
    
    print_info "$PACKAGE Downloading LazyTodo for macOS..."
    
    # Get release information
    local api_url="https://api.github.com/repos/yourusername/lazytodo/releases/latest"
    if [[ "$VERSION" != "latest" ]]; then
        api_url="https://api.github.com/repos/yourusername/lazytodo/releases/tags/$VERSION"
    fi
    
    if command_exists curl; then
        local release_info
        release_info=$(curl -s "$api_url")
    else
        print_error "curl is required for installation. Please install it via: brew install curl"
    fi
    
    # Extract download URL for our platform
    asset_name="lazytodo-darwin-${arch}.tar.gz"
    download_url=$(echo "$release_info" | grep -o "\"browser_download_url\": \"[^\"]*${asset_name}\"" | cut -d '"' -f 4)
    
    if [[ -z "$download_url" ]]; then
        print_error "No binary found for macOS-${arch} in release"
    fi
    
    print_info "$PACKAGE Downloading: $asset_name"
    
    # Download
    temp_file="/tmp/lazytodo-darwin-${arch}.tar.gz"
    curl -L -o "$temp_file" "$download_url"
    
    print_success "Downloaded successfully"
    echo "$temp_file"
}

# Create directories
create_directories() {
    print_info "$FOLDER Creating directories..."
    
    # Create install directory
    if [[ ! -d "$INSTALL_DIR" ]]; then
        if [[ "$INSTALL_DIR" == "/usr/local/bin" ]]; then
            sudo mkdir -p "$INSTALL_DIR"
        else
            mkdir -p "$INSTALL_DIR"
        fi
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
    binary_path=$(find "/tmp" -name "lazytodo" -type f -perm +111 | head -1)
    
    if [[ -z "$binary_path" ]]; then
        print_error "Binary not found in downloaded archive"
    fi
    
    # Install
    if [[ "$INSTALL_DIR" == "/usr/local/bin" ]] && [[ ! -w "$INSTALL_DIR" ]]; then
        sudo cp "$binary_path" "$INSTALL_DIR/lazytodo"
        sudo chmod +x "$INSTALL_DIR/lazytodo"
    else
        cp "$binary_path" "$INSTALL_DIR/lazytodo"
        chmod +x "$INSTALL_DIR/lazytodo"
    fi
    
    # Cleanup
    rm -f "$temp_file"
    rm -rf "$(dirname "$binary_path")" 2>/dev/null || true
    
    print_success "Binary installed to $INSTALL_DIR/lazytodo"
}

# Configure PATH for macOS
configure_path() {
    print_info "$GEAR Configuring PATH..."
    
    # /usr/local/bin is usually in PATH by default
    if [[ "$INSTALL_DIR" == "/usr/local/bin" ]]; then
        print_info "Installation directory is in default PATH"
        return
    fi
    
    # Check if directory is already in PATH
    if echo "$PATH" | grep -q "$INSTALL_DIR"; then
        print_info "Install directory already in PATH"
        return
    fi
    
    # Detect shell and add to appropriate profile
    local shell_profile=""
    local current_shell
    current_shell=$(basename "$SHELL")
    
    case $current_shell in
        bash)
            if [[ -f "$HOME/.bash_profile" ]]; then
                shell_profile="$HOME/.bash_profile"
            elif [[ -f "$HOME/.bashrc" ]]; then
                shell_profile="$HOME/.bashrc"
            else
                shell_profile="$HOME/.bash_profile"
                touch "$shell_profile"
            fi
            ;;
        zsh)
            shell_profile="$HOME/.zshrc"
            if [[ ! -f "$shell_profile" ]]; then
                touch "$shell_profile"
            fi
            ;;
        fish)
            local fish_config_dir="$HOME/.config/fish"
            mkdir -p "$fish_config_dir"
            shell_profile="$fish_config_dir/config.fish"
            if [[ ! -f "$shell_profile" ]]; then
                touch "$shell_profile"
            fi
            # Fish uses different syntax
            echo "" >> "$shell_profile"
            echo "# Added by LazyTodo installer" >> "$shell_profile"
            echo "set -gx PATH \$PATH $INSTALL_DIR" >> "$shell_profile"
            print_success "Added $INSTALL_DIR to PATH in $shell_profile"
            print_warning "Please restart your terminal or run: source $shell_profile"
            return
            ;;
        *)
            shell_profile="$HOME/.profile"
            if [[ ! -f "$shell_profile" ]]; then
                touch "$shell_profile"
            fi
            ;;
    esac
    
    if [[ -n "$shell_profile" ]] && [[ "$current_shell" != "fish" ]]; then
        echo "" >> "$shell_profile"
        echo "# Added by LazyTodo installer" >> "$shell_profile"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$shell_profile"
        print_success "Added $INSTALL_DIR to PATH in $shell_profile"
        print_warning "Please restart your terminal or run: source $shell_profile"
    fi
}

# Check for Homebrew and suggest it
check_homebrew() {
    if command_exists brew; then
        print_info "🍺 Homebrew detected! Consider creating a Homebrew formula for easier updates."
        print_info "   You can still use this installer or check if LazyTodo is available via: brew search lazytodo"
    else
        print_info "💡 Tip: Consider installing Homebrew (https://brew.sh) for easier package management on macOS"
    fi
}

# Create application bundle (optional)
create_app_bundle() {
    local app_dir="/Applications/LazyTodo.app"
    local contents_dir="$app_dir/Contents"
    local macos_dir="$contents_dir/MacOS"
    local resources_dir="$contents_dir/Resources"
    
    print_info "📱 Creating macOS application bundle..."
    
    # Create directory structure
    sudo mkdir -p "$macos_dir" "$resources_dir"
    
    # Copy binary
    sudo cp "$INSTALL_DIR/lazytodo" "$macos_dir/LazyTodo"
    
    # Create Info.plist
    sudo tee "$contents_dir/Info.plist" > /dev/null << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>LazyTodo</string>
    <key>CFBundleIdentifier</key>
    <string>com.lazytodo.app</string>
    <key>CFBundleName</key>
    <string>LazyTodo</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF
    
    # Create wrapper script that opens Terminal
    sudo tee "$macos_dir/LazyTodo" > /dev/null << EOF
#!/bin/bash
osascript -e 'tell application "Terminal" to do script "$INSTALL_DIR/lazytodo"'
EOF
    
    sudo chmod +x "$macos_dir/LazyTodo"
    
    print_success "Application bundle created at $app_dir"
}

# Main installation process
main() {
    echo
    print_info "$ROCKET$APPLE LazyTodo macOS Installer"
    print_info "================================="
    echo
    
    # Check macOS version
    local macos_version
    macos_version=$(sw_vers -productVersion)
    print_info "macOS Version: $macos_version"
    
    # Check for Homebrew
    check_homebrew
    echo
    
    # Check dependencies
    print_info "Checking dependencies..."
    if ! command_exists curl; then
        print_error "curl is required. Install with: brew install curl"
    fi
    
    if ! command_exists tar; then
        print_error "tar is required but should be available by default"
    fi
    
    print_success "Dependencies check passed"
    
    # Check permissions
    if ! check_sudo && [[ "$INSTALL_DIR" == "/usr/local/bin" ]]; then
        print_info "Switching to user installation directory..."
        INSTALL_DIR="$HOME/.local/bin"
    fi
    
    # Create directories
    create_directories
    
    # Download and install
    local temp_file
    temp_file=$(download_release)
    install_binary "$temp_file"
    
    # Configure environment
    configure_path
    
    # Create app bundle (optional)
    read -p "📱 Create macOS application bundle in /Applications? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        create_app_bundle
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
    print_info "🍺 Consider: brew install lazytodo (if available)"
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
    
    print_success "LazyTodo is ready for macOS! Happy organizing! $CHECKMARK"
}

# Show help
show_help() {
    cat << EOF
LazyTodo macOS Installation Script

Usage: $0 [VERSION]

Arguments:
  VERSION     Version to install (default: latest)

Environment Variables:
  INSTALL_DIR Directory to install binary (default: /usr/local/bin)

Examples:
  $0                    # Install latest version
  $0 v1.0.0            # Install specific version
  INSTALL_DIR="\$HOME/.local/bin" $0  # Install to user directory

macOS Notes:
  • Requires macOS 10.12 or later
  • Universal binary supports both Intel and Apple Silicon
  • Creates optional .app bundle for Launchpad integration

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
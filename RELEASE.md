# Release Guide

This document outlines the process for creating releases of LazyTodo using GoReleaser.

## 🚀 Release Process

### Pre-Release Checklist

- [ ] All tests are passing (check GitHub Actions CI)
- [ ] Version number has been determined following [Semantic Versioning](https://semver.org/)
- [ ] CHANGELOG.md has been updated with new features, fixes, and breaking changes
- [ ] Documentation is up to date
- [ ] Installation scripts have been tested on target platforms
- [ ] No critical security vulnerabilities in dependencies
- [ ] GoReleaser configuration is tested locally

### Release Types

#### Major Release (x.0.0)
- Breaking changes
- Major new features
- Significant architecture changes

#### Minor Release (x.y.0)
- New features (backward compatible)
- Deprecations (with migration path)
- Performance improvements

#### Patch Release (x.y.z)
- Bug fixes
- Security patches
- Documentation updates

### Creating a Release

#### Method 1: Git Tag (Recommended)

1. **Create and Push Tag**:
   ```bash
   # Create annotated tag
   git tag -a v1.0.0 -m "Release version 1.0.0"
   
   # Push tag to trigger release
   git push origin v1.0.0
   ```

2. **GoReleaser will automatically**:
   - Build binaries for all supported platforms
   - Create GitHub release with generated changelog
   - Upload release assets with checksums
   - Update Homebrew tap (if configured)
   - Build and push Docker images
   - Create Snap packages

#### Method 2: GitHub Actions Manual Trigger

1. **Manual Trigger via GitHub UI**:
   - Go to Actions tab in GitHub
   - Select "Release" workflow
   - Click "Run workflow"
   - Enter version (e.g., `v1.0.0`)
   - Mark as pre-release if needed
   - Click "Run workflow"

### Post-Release Tasks

- [ ] Verify release assets are uploaded correctly
- [ ] Test installation scripts with new release
- [ ] Verify Homebrew tap was updated
- [ ] Test Docker images
- [ ] Announce release on relevant channels
- [ ] Monitor for issues with the new release

## 📦 Supported Platforms & Package Managers

GoReleaser builds and distributes LazyTodo through multiple channels:

### **Binary Releases**

| Platform | Architecture | Archive Format |
|----------|-------------|----------------|
| Windows  | x64         | `lazytodo-windows-amd64.zip` |
| Windows  | ARM64       | `lazytodo-windows-arm64.zip` |
| Linux    | x64         | `lazytodo-linux-amd64.tar.gz` |
| Linux    | ARM64       | `lazytodo-linux-arm64.tar.gz` |
| Linux    | ARM         | `lazytodo-linux-armv7.tar.gz` |
| macOS    | Intel       | `lazytodo-darwin-amd64.tar.gz` |
| macOS    | Apple Silicon | `lazytodo-darwin-arm64.tar.gz` |
| FreeBSD  | x64         | `lazytodo-freebsd-amd64.tar.gz` |

### **Package Managers**

- **Homebrew** (macOS/Linux): `brew install yourusername/tap/lazytodo`
- **Snap** (Linux): `snap install lazytodo`
- **Docker**: `docker run ghcr.io/yourusername/lazytodo:latest`

## 🔧 GoReleaser Configuration

### Key Features in `.goreleaser.yml`

- **Cross-compilation**: Automatic builds for all platforms
- **Archive Creation**: Platform-specific formats (zip for Windows, tar.gz for others)
- **Checksums**: SHA256 verification files
- **Homebrew Integration**: Automatic tap updates
- **Docker Images**: Multi-platform container builds
- **Snap Packages**: Linux universal packages
- **Changelog Generation**: Automatic from conventional commits

### Environment Variables

Required secrets for full functionality:

- `GITHUB_TOKEN`: Automatically provided by GitHub Actions
- `HOMEBREW_TAP_GITHUB_TOKEN`: Personal access token for Homebrew tap updates

### Build Configuration

Binaries are built with:

```bash
go build -ldflags="-s -w -X main.version={{.Version}} -X main.buildTime={{.Date}} -X main.gitCommit={{.FullCommit}}"
```

## 🧪 Testing Releases

### Local Testing with GoReleaser

1. **Install GoReleaser**:
   ```bash
   # macOS
   brew install goreleaser
   
   # Linux
   curl -sfL https://goreleaser.com/static/run | bash
   
   # Or download from GitHub releases
   ```

2. **Test Build (without publishing)**:
   ```bash
   # Dry run to check configuration
   goreleaser check
   
   # Build snapshot without releasing
   goreleaser build --snapshot --clean
   
   # Full release test (creates local dist/ folder)
   goreleaser release --snapshot --clean
   ```

3. **Validate Configuration**:
   ```bash
   # Check .goreleaser.yml syntax
   goreleaser check
   
   # Generate config documentation
   goreleaser jsonschema -o goreleaser.schema.json
   ```

### Pre-Release Testing

1. **Build Verification**:
   ```bash
   # Test local build
   go build -o lazytodo cmd/main.go
   ./lazytodo --version
   ```

2. **Cross-Platform Testing**:
   ```bash
   # Test with GoReleaser
   goreleaser build --snapshot --clean
   
   # Test specific platforms
   find dist/ -name "lazytodo*" -executable
   ```

3. **Installation Script Testing**:
   ```bash
   # Test script syntax
   bash -n scripts/install-linux.sh
   bash -n scripts/install-mac.sh
   
   # Test PowerShell syntax (if on Windows/WSL)
   pwsh -Command "Get-Content scripts/install-windows.ps1 | Out-Null"
   ```

### Post-Release Verification

1. **Download and Test**:
   ```bash
   # Test Homebrew installation
   brew install yourusername/tap/lazytodo
   lazytodo --version
   
   # Test script installation
   curl -fsSL https://raw.githubusercontent.com/yourusername/lazytodo/main/scripts/install-linux.sh | bash
   ```

2. **Docker Testing**:
   ```bash
   # Test Docker image
   docker run --rm ghcr.io/yourusername/lazytodo:latest --version
   ```

## 🔄 Hotfix Process

For critical bugs that need immediate fixes:

1. **Create Hotfix Branch**:
   ```bash
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. **Apply Fix and Test**:
   ```bash
   # Make necessary changes
   git commit -m "fix: critical bug description"
   
   # Test with GoReleaser
   goreleaser build --snapshot --clean
   ```

3. **Release Hotfix**:
   ```bash
   git tag -a v1.0.1 -m "Hotfix release v1.0.1"
   git push origin v1.0.1
   ```

4. **Merge Back**:
   ```bash
   git checkout main
   git merge hotfix/v1.0.1
   git push origin main
   ```

## 📝 Release Notes

GoReleaser automatically generates release notes from conventional commits:

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

### Supported Types for Changelog

- `feat`: 🚀 New Features
- `fix`: 🐛 Bug Fixes  
- `perf`: ⚡ Performance Improvements
- Others: 🔧 Other Changes

### Example Commits

```bash
git commit -m "feat(ui): add dark mode support"
git commit -m "fix(storage): resolve database connection issue"
git commit -m "perf(rendering): optimize window drawing"
```

## 🚨 Emergency Procedures

### Yanking a Release

If a release has critical issues:

1. **Mark as Pre-release** in GitHub UI
2. **Create Hotfix** following the hotfix process
3. **Update Package Managers**:
   ```bash
   # Remove from Homebrew tap if needed
   # Contact Snap store for removal
   ```
4. **Communicate** the issue to users

### Configuration Issues

If GoReleaser configuration fails:

1. **Test Locally**:
   ```bash
   goreleaser check
   goreleaser build --snapshot --clean
   ```

2. **Check GitHub Actions Logs**:
   - Review the workflow run logs
   - Check for missing environment variables
   - Verify cross-compilation tools

3. **Rollback Configuration**:
   ```bash
   git revert <bad-commit>
   git tag -d v1.0.0  # Delete bad tag
   git push origin :refs/tags/v1.0.0
   ```

## 📊 Release Metrics

Track these metrics for each release:

- Download counts by platform and package manager
- Installation success rates via different methods
- Issue reports within 48 hours
- Performance of different distribution channels

## 🤝 Communication

### Release Channels

GoReleaser automatically handles:
- ✅ GitHub Releases
- ✅ Homebrew Tap updates
- ✅ Docker registry pushes
- ✅ Snap store uploads (with proper store connection)

Manual announcements:
- GitHub Discussions
- Project README updates
- Social media

### Package Manager Setup

#### Homebrew Tap Setup

1. **Create Homebrew Tap Repository**:
   ```bash
   # Create new repository: homebrew-tap
   # Enable GitHub Pages (optional)
   ```

2. **Configure Secrets**:
   - Create `HOMEBREW_TAP_GITHUB_TOKEN` secret
   - Grant token access to the tap repository

3. **First Release**:
   ```bash
   # GoReleaser will create Formula/lazytodo.rb automatically
   ```

#### Snap Store Setup

1. **Register Snap Name**: Register "lazytodo" on snapcraft.io
2. **Store Credentials**: Configure snapcraft store credentials
3. **Auto-publish**: Configure in `.goreleaser.yml`

## 🎯 Best Practices

1. **Always test locally** with `goreleaser build --snapshot --clean`
2. **Use conventional commits** for better changelogs
3. **Keep .goreleaser.yml** up to date with GoReleaser releases
4. **Monitor package managers** for successful updates
5. **Test installation methods** after each release
6. **Keep dependencies minimal** for easier cross-compilation

## 📚 Resources

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [Homebrew Tap Creation](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [Snapcraft Documentation](https://snapcraft.io/docs) 
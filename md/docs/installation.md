---
title: Installation
description: Installation documentation for Base Framework.
---

# Installation
Get started with Base Framework by installing the CLI tool on your system.
## Installation
### Quick Install
Install Base CLI with our installation script. This will download the latest version and add it to your PATH:
Terminal
Copy
```
curl -fsSL https://get.base.al | sh
```
### Manual Installation
Download the latest release for your platform:
Platform
Architecture
Download
macOS
amd64 / arm64
GitHub Releases
Linux
amd64 / arm64
GitHub Releases
Windows
amd64
GitHub Releases
### Verify Installation
After installation, verify that Base CLI is correctly installed and accessible:
Terminal
Copy
```
base version
```
You should see output similar to:
```
Base CLI v2.0.4
Commit: 690c60cbd06ea1fd6415b2b44d1de3bca476deb4
Built: 2025-08-21T01:17:14Z
Go version: go1.24.5
```
### Next Steps
Once Base CLI is installed, you can:
- **Create a new project:** base new my-project
- **Generate structures:** base g User name:string email:string
- **Start development server:** base start -r
- **View help:** base help
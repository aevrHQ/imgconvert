# Installation Guide - Use imgconvert Anywhere

## Option 1: Install Globally (Recommended)

Copy the binary to a location in your PATH:

```bash
cd /Users/miracleio/Documents/devprojects/image-converter
./build.sh
sudo cp imgconvert /usr/local/bin/
```

Now you can use it from anywhere:
```bash
cd ~/Pictures
imgconvert *.jpg

cd /any/directory
imgconvert -r .
```

## Option 2: Add to Your PATH (No sudo required)

Add the project directory to your PATH:

```bash
# For Zsh (macOS default)
echo 'export PATH="$HOME/Documents/devprojects/image-converter:$PATH"' >> ~/.zshrc
source ~/.zshrc

# For Bash
echo 'export PATH="$HOME/Documents/devprojects/image-converter:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Now use it from anywhere:
```bash
imgconvert photo.jpg
```

## Option 3: Create an Alias

Add to your shell config:

```bash
# For Zsh
echo 'alias imgconvert="$HOME/Documents/devprojects/image-converter/imgconvert"' >> ~/.zshrc
source ~/.zshrc

# For Bash
echo 'alias imgconvert="$HOME/Documents/devprojects/image-converter/imgconvert"' >> ~/.bashrc
source ~/.bashrc
```

## Option 4: Install via Go (For systems with Go)

If you want to install on other machines with Go:

```bash
# From any directory
go install github.com/aevrHQ/imgconvert/cmd/imgconvert@latest

# Or from source
cd /path/to/image-converter
go install ./cmd/imgconvert
```

This installs to `$GOPATH/bin` (usually `~/go/bin`)

## Option 5: Copy Binary Directly

Simply copy the binary wherever you need it:

```bash
# Copy to specific project
cp imgconvert /path/to/other/project/

# Use it there
cd /path/to/other/project
./imgconvert *.jpg
```

## Option 6: Build for Other Platforms

Build for different operating systems:

```bash
cd /Users/miracleio/Documents/devprojects/image-converter
export PATH="$HOME/gosdk/bin:$PATH"

# For Linux
GOOS=linux GOARCH=amd64 go build -o imgconvert-linux ./cmd/imgconvert

# For Windows
GOOS=windows GOARCH=amd64 go build -o imgconvert.exe ./cmd/imgconvert

# For macOS Intel
GOOS=darwin GOARCH=amd64 go build -o imgconvert-macos-intel ./cmd/imgconvert

# For macOS Apple Silicon (M1/M2/M3)
GOOS=darwin GOARCH=arm64 go build -o imgconvert-macos-arm ./cmd/imgconvert
```

Transfer these binaries to other machines and run them directly.

## Quick Setup Script

Run this for instant global installation:

```bash
cd /Users/miracleio/Documents/devprojects/image-converter
./build.sh && sudo cp imgconvert /usr/local/bin/
imgconvert --version
```

## Verify Installation

After any method above, test it:

```bash
# Check if accessible
which imgconvert

# Test it
imgconvert --version

# Use it from any directory
cd ~
imgconvert --help
```

## Recommended: Global Installation

**Best approach for daily use:**

```bash
cd /Users/miracleio/Documents/devprojects/image-converter
./build.sh
sudo cp imgconvert /usr/local/bin/
```

Then from ANYWHERE on your system:
```bash
cd ~/Pictures/vacation-2024
imgconvert -r --preset web .

cd ~/Desktop
imgconvert photo.jpg

cd /any/path
imgconvert -q 90 *.png
```

## Uninstall

If you ever want to remove it:

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/imgconvert

# If added to PATH, remove from ~/.zshrc or ~/.bashrc

# If using alias, remove from ~/.zshrc or ~/.bashrc
```

## Share with Others

### Simple sharing:
```bash
# Build
cd /Users/miracleio/Documents/devprojects/image-converter
./build.sh

# Send the binary to someone
# They can run it directly:
./imgconvert photo.jpg
```

### Professional sharing:
Create a release package:
```bash
cd /Users/miracleio/Documents/devprojects/image-converter
./build.sh

# Create distribution
mkdir -p dist
cp imgconvert dist/
cp README.md dist/
cp docs/quickstart.md dist/USAGE.txt
cp LICENSE dist/

# Create archive
tar -czf imgconvert-v0.1.0-macos.tar.gz -C dist .
```

Now share `imgconvert-v0.1.0-macos.tar.gz` with others!

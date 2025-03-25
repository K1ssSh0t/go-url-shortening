#!/bin/bash

# Script to update Go to the latest stable version

set -e # Exit immediately if a command exits with a non-zero status.

# --- Configuration ---
GO_INSTALL_DIR="/usr/local/go" # Default Go installation directory
TMP_DIR="/tmp"
LATEST_VERSION_URL="https://go.dev/dl/"

# --- Helper Functions ---

log_info() {
  echo "[INFO] $1"
}

log_warning() {
  echo "[WARNING] $1"
}

log_error() {
  echo "[ERROR] $1"
}

check_sudo() {
  if [[ "$EUID" -ne 0 ]]; then
    log_warning "This script may require sudo privileges for certain operations."
  fi
}

get_current_go_version() {
  if command -v go &> /dev/null; then
    go version
  else
    echo "Go is not currently installed or not in your PATH."
    return 1
  fi
}

get_latest_go_version() {
  log_info "Fetching the latest stable Go version from $LATEST_VERSION_URL..."
  latest_version=$(curl -s "$LATEST_VERSION_URL" | grep -oE 'go+\.+\.+' | head -n 1)
  if [[ -z "$latest_version" ]]; then
    log_error "Failed to retrieve the latest Go version."
    return 1
  fi
  echo "$latest_version"
}

download_go() {
  local version="$1"
  local os=$(uname -s | tr '[:upper:]' '[:lower:]')
  local arch=$(uname -m)
  local filename="go${version}.${os}-${arch}.tar.gz"
  local download_url="$LATEST_VERSION_URL/$filename"
  local tmp_file="$TMP_DIR/$filename"

  log_info "Downloading Go version $version for $os/$arch from $download_url..."
  if wget -q "$download_url" -O "$tmp_file"; then
    echo "$tmp_file"
  else
    log_error "Failed to download Go version $version."
    return 1
  fi
}

install_go() {
  local tmp_file="$1"

  check_sudo

  log_info "Removing existing Go installation from $GO_INSTALL_DIR..."
  if sudo rm -rf "$GO_INSTALL_DIR"; then
    log_info "Extracting the new Go version to $GO_INSTALL_DIR..."
    if sudo tar -C /usr/local -xzf "$tmp_file"; then
      log_info "Go version updated successfully!"
    else
      log_error "Failed to extract the Go archive."
      return 1
    fi
  else
    log_error "Failed to remove the existing Go installation. Please ensure you have the necessary permissions."
    return 1
  fi
}

cleanup() {
  local tmp_file="$1"
  if [[ -f "$tmp_file" ]]; then
    log_info "Cleaning up temporary file: $tmp_file"
    rm "$tmp_file"
  fi
}

update_path_instructions() {
  log_info "-------------------------------------------------------------------"
  log_info "Important: You might need to update your system's PATH environment variable."
  log_info "Add the following line to your shell configuration file (e.g., ~/.bashrc, ~/.zshrc):"
  log_info "export PATH=$PATH:$GO_INSTALL_DIR/bin"
  log_info "Then, apply the changes by running: source ~/.bashrc or source ~/.zshrc (or similar)."
  log_info "-------------------------------------------------------------------"
}

# --- Main Script ---

log_info "Starting Go version update script..."
check_sudo

log_info "Current Go version:"
get_current_go_version

log_info ""

latest_version=$(get_latest_go_version)
if [[ $? -ne 0 ]]; then
  exit 1
fi
log_info "Latest stable Go version available: $latest_version"

tmp_file=$(download_go "$latest_version")
if [[ $? -ne 0 ]]; then
  exit 1
fi

install_go "$tmp_file"
if [[ $? -ne 0 ]]; then
  cleanup "$tmp_file"
  exit 1
fi

cleanup "$tmp_file"

update_path_instructions

log_info ""
log_info "Please restart your terminal or source your shell configuration file to apply the PATH changes."
log_info "You can verify the updated Go version by running: go version"

log_info "Go version update script finished."

exit 0
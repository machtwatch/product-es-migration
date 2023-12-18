#!/usr/bin/env bash

# This script is used for setting up the development environment for the Catalyst SDK.
# It performs various checks and installations for required tools and sets up Git hooks.

echo -e "Checking available Homebrew.."
# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Homebrew is not installed. Downloading..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi
echo -e "Homebrew OK \n"

echo "Checking available Vault..."
if ! command -v vault &> /dev/null; then
    echo "Vault is not installed. Proceed to download and install"
    brew install vault
fi
echo -e "Vault OK \n"


echo -e "\n✨ Setup success!\n"

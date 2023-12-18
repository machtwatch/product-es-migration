#!/bin/bash

# Read the input parameters: GitHub token, Vault secret, and Environment
github_token=$1
vault_secret=$2
environment=$3
env_output_file_name=".env"
hide_url=false

# Check if the GitHub token is provided; exit if not
if [[ -z $github_token ]]; then
    echo "Your local GITHUB_TOKEN env variable is not set, see: https://jamtangan.atlassian.net/wiki/spaces/EN/pages/2344190011/Code+Repository#Work-Station-Setup for guidance.";
    exit 1
fi

# Check if the Vault secret is provided; exit if not
if [[ -z $vault_secret ]]; then
    echo "Vault secret variable has not been provided.";
    exit 1
fi

# Check if the Environment is provided; if not, set default values for development environment
if [[ -z $environment ]]; then
    environment="development"
    env_output_file_name=".env"
    hide_url=true
fi

# Set the Vault address
export VAULT_ADDR=http://vault.ctlyst.id:8200/

# Authenticate to Vault using GitHub token
vault login -method=github token=$github_token > /dev/null

# Retrieve the specified Vault secret for the given environment, extract key-value pairs starting from the 16th line,
# and optionally replace URLs with "localhost" if hide_url is set to true.
vault kv get $environment/$vault_secret | awk 'NR >= 16 {print $1 "=" $2}' |  { if $hide_url; then cat | sed -E 's#(http[s]?://)[^/ ,]+#\1localhost#g';else cat ;fi; } > $env_output_file_name

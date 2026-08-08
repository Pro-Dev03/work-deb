#!/bin/bash

# Script to setup GitHub Secrets using GitHub API
# This requires Python with libsodium for encryption

set -e

REPO="titojbxn-cell/worktrack-v2"
API_BASE="https://api.github.com"

echo "=== GitHub Secrets Setup (API Method) ==="
echo ""

# Check for required tools
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is required but not installed."
    exit 1
fi

# Install libsodium if not available
python3 -c "import nacl" 2>/dev/null || {
    echo "📦 Installing PyNaCl for encryption..."
    pip3 install --user pynacl || {
        echo "❌ Failed to install PyNaCl. Please install manually:"
        echo "   pip3 install pynacl"
        exit 1
    }
}

# Get GitHub token
read -p "Enter your GitHub Personal Access Token (with repo scope): " GITHUB_TOKEN

if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ GitHub token is required."
    exit 1
fi

echo ""
echo "📋 Repository: $REPO"
echo ""

# Get repository public key
echo "🔑 Getting repository public key..."
KEY_RESPONSE=$(curl -s -X GET \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "$API_BASE/repos/$REPO/actions/secrets/public-key")

KEY_ID=$(echo "$KEY_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['key_id'])")
PUBLIC_KEY=$(echo "$KEY_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['key'])")

echo "✅ Public key retrieved (ID: $KEY_ID)"
echo ""

# Function to encrypt and add secret
add_secret() {
    local secret_name=$1
    local secret_value=$2

    echo "🔐 Encrypting $secret_name..."

    # Encrypt using Python
    ENCRYPTED_VALUE=$(python3 <<EOF
import base64
import json
import nacl.pwhash
import nacl.secret
import nacl.utils
from nacl import encoding, public

def encrypt(public_key: str, secret_value: str) -> str:
    """Encrypt a Unicode string using the public key."""
    public_key = public.PublicKey(public_key.encode("utf-8"), encoding.Base64Encoder())
    sealed_box = public.SealedBox(public_key)
    encrypted = sealed_box.encrypt(secret_value.encode("utf-8"))
    return encrypted.decode("utf-8")

encrypted = encrypt("$PUBLIC_KEY", "$secret_value")
print(encrypted)
EOF
)

    # Add secret to GitHub
    echo "📤 Adding $secret_name to GitHub..."
    curl -s -X PUT \
      -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github+json" \
      "$API_BASE/repos/$REPO/actions/secrets/$secret_name" \
      -d "{\"encrypted_value\":\"$ENCRYPTED_VALUE\",\"key_id\":\"$KEY_ID\"}" \
      > /dev/null

    echo "✅ $secret_name added successfully"
}

# Get secrets from user
echo "Please provide the following information:"
echo ""

read -p "1. Enter your Codemagic API Token: " CODEMAGIC_API_TOKEN
read -p "2. Enter your Codemagic App ID: " CODEMAGIC_APP_ID
read -p "3. Enter your VITE API Base URL (e.g., https://api.example.com/api/v1): " VITE_API_BASE_URL

echo ""
echo "🔐 Adding secrets..."

# Add secrets
add_secret "CODEMAGIC_API_TOKEN" "$CODEMAGIC_API_TOKEN"
add_secret "CODEMAGIC_APP_ID" "$CODEMAGIC_APP_ID"
add_secret "VITE_API_BASE_URL" "$VITE_API_BASE_URL"

echo ""
echo "✅ All secrets added successfully!"
echo ""
echo "Next steps:"
echo "1. Go to Codemagic (codemagic.io)"
echo "2. Add your project (worktrack-v2)"
echo "3. Click 'Check for configuration file'"
echo "4. The build will be triggered automatically on next push"

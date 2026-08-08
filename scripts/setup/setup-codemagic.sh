#!/bin/bash

# Script to setup GitHub Secrets for Codemagic integration
# Usage: ./setup-codemagic.sh
#
# 📝 وظيفة هذا السكريبت:
# - إعداد GitHub Secrets للتكامل مع Codemagic
# - تخزين المفاتيح الحساسة بشكل آمن
# - تهيئة البناء السحابي التلقائي
#
# 🚀 الاستخدام: ./scripts/setup/setup-codemagic.sh

set -e

echo "=== Codemagic GitHub Secrets Setup ==="
echo ""

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed."
    echo "Please install it first: https://cli.github.com/"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Not authenticated with GitHub CLI."
    echo "Please run: gh auth login"
    exit 1
fi

REPO="titojbxn-cell/worktrack-v2"

echo "📋 Repository: $REPO"
echo ""

# Get secrets from user
echo "Please provide the following information:"
echo ""

read -p "1. Enter your Codemagic API Token: " CODEMAGIC_API_TOKEN
read -p "2. Enter your Codemagic App ID: " CODEMAGIC_APP_ID
read -p "3. Enter your VITE API Base URL (e.g., https://api.example.com/api/v1): " VITE_API_BASE_URL

echo ""
echo "🔐 Adding secrets to GitHub..."

# Add secrets
gh secret set CODEMAGIC_API_TOKEN -b"$CODEMAGIC_API_TOKEN" --repo "$REPO"
gh secret set CODEMAGIC_APP_ID -b"$CODEMAGIC_APP_ID" --repo "$REPO"
gh secret set VITE_API_BASE_URL -b"$VITE_API_BASE_URL" --repo "$REPO"

echo ""
echo "✅ Secrets added successfully!"
echo ""
echo "Next steps:"
echo "1. Go to Codemagic (codemagic.io)"
echo "2. Add your project (worktrack-v2)"
echo "3. Click 'Check for configuration file'"
echo "4. The build will be triggered automatically on next push"

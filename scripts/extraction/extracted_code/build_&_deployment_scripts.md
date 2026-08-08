# Build & Deployment Scripts

## تم الاستخراج في: 2026-08-06 02:09:34

## عدد الملفات: 9

---

## 📄 build-all-release.sh

```bash
#!/bin/bash

# WorkTrack - بناء جميع التطبيقات للإصدار الإنتاجي
# هذا السكريبت يبني نسخة Release موقعة لجميع التطبيقات الثلاثة

echo "🚀 بدء بناء تطبيقات WorkTrack للإصدار الإنتاجي"
echo "=================================================="
echo ""

# التحقق من وجود Keystore
if [ ! -f "android-keystores/worktrack-release.keystore" ]; then
    echo "❌ خطأ: لم يتم العثور على ملف Keystore"
    echo "يرجى تشغيل ./create-keystore.sh أولاً"
    exit 1
fi

# التحقق من كلمات المرور
echo "⚠️  تحذير: تأكد من تغيير كلمات المرور في ملفات keystore.properties"
echo ""
read -p "هل قمت بتغيير كلمات المرور؟ (y/n): " CONFIRM_PASSWORDS

if [ "$CONFIRM_PASSWORDS" != "y" ] && [ "$CONFIRM_PASSWORDS" != "Y" ]; then
    echo "❌ يرجى تغيير كلمات المرور أولاً في ملفات keystore.properties"
    exit 1
fi

# بناء تطبيق الموظف
echo ""
echo "📱 بناء تطبيق الموظف (Worker App)..."
echo "-----------------------------------"
cd frontend-worker-pwa

# التحقق من وجود node_modules
if [ ! -d "node_modules" ]; then
    echo "📦 تثبيت الاعتماديات..."
    npm install
fi

# بناء التطبيق
echo "🔨 بناء التطبيق..."
npm run build

if [ $? -ne 0 ]; then
    echo "❌ فشل بناء تطبيق الموظف"
    cd ..
    exit 1
fi

# مزامنة مع Android
echo "🔄 مزامنة مع Android..."
npx cap sync android

if [ $? -ne 0 ]; then
    echo "❌ فشلت المزامنة مع Android"
    cd ..
    exit 1
fi

# بناء Release APK
echo "📦 بناء Release APK..."
cd android
./gradlew assembleRelease

if [ $? -ne 0 ]; then
    echo "❌ فشل بناء Release APK"
    cd ../..
    exit 1
fi

echo "✅ تم بناء تطبيق الموظف بنجاح"
cd ../..

# بناء تطبيق الإدارة
echo ""
echo "📱 بناء تطبيق الإدارة (Admin App)..."
echo "-------------------------------------"
cd frontend-admin-dashboard

# التحقق من وجود node_modules
if [ ! -d "node_modules" ]; then
    echo "📦 تثبيت الاعتماديات..."
    npm install
fi

# بناء التطبيق
echo "🔨 بناء التطبيق..."
npm run build

if [ $? -ne 0 ]; then
    echo "❌ فشل بناء تطبيق الإدارة"
    cd ..
    exit 1
fi

# مزامنة مع Android
echo "🔄 مزامنة مع Android..."
npx cap sync android

if [ $? -ne 0 ]; then
    echo "❌ فشلت المزامنة مع Android"
    cd ..
    exit 1
fi

# بناء Release APK
echo "📦 بناء Release APK..."
cd android
./gradlew assembleRelease

if [ $? -ne 0 ]; then
    echo "❌ فشل بناء Release APK"
    cd ../..
    exit 1
fi

echo "✅ تم بناء تطبيق الإدارة بنجاح"
cd ../..

# بناء تطبيق العميل
echo ""
echo "📱 بناء تطبيق العميل (Client App)..."
echo "--------------------------------------"
cd frontend-client-portal

# التحقق من وجود node_modules
if [ ! -d "node_modules" ]; then
    echo "📦 تثبيت الاعتماديات..."
    npm install
fi

# بناء التطبيق
echo "🔨 بناء التطبيق..."
npm run build

if [ $? -ne 0 ]; then
    echo "❌ فشل بناء تطبيق العميل"
    cd ..
    exit 1
fi

# مزامنة مع Android
echo "🔄 مزامنة مع Android..."
npx cap sync android

if [ $? -ne 0 ]; then
    echo "❌ فشلت المزامنة مع Android"
    cd ..
    exit 1
fi

# بناء Release APK
echo "📦 بناء Release APK..."
cd android
./gradlew assembleRelease

if [ $? -ne 0 ]; then
    echo "❌ فشل بناء Release APK"
    cd ../..
    exit 1
fi

echo "✅ تم بناء تطبيق العميل بنجاح"
cd ../..

# إنشاء مجلد للملفات النهائية
echo ""
echo "📁 إنشاء مجلد للملفات النهائية..."
mkdir -p release-apks

# نسخ ملفات APK
cp frontend-worker-pwa/android/app/build/outputs/apk/release/app-release.apk release-apks/worktrack-worker-release.apk
cp frontend-admin-dashboard/android/app/build/outputs/apk/release/app-release.apk release-apks/worktrack-admin-release.apk
cp frontend-client-portal/android/app/build/outputs/apk/release/app-release.apk release-apks/worktrack-client-release.apk

echo ""
echo "=================================================="
echo "✅ تم بناء جميع التطبيقات بنجاح!"
echo "=================================================="
echo ""
echo "📦 ملفات APK:"
echo "   - Worker: release-apks/worktrack-worker-release.apk"
echo "   - Admin: release-apks/worktrack-admin-release.apk"
echo "   - Client: release-apks/worktrack-client-release.apk"
echo ""
echo "📋 الخطوات التالية:"
echo "1. اختبر ملفات APK على أجهزة حقيقية"
echo "2. ارفع التطبيقات إلى Google Play Console"
echo "3. أكمل إعداد App Listing و Metadata"
echo "4. أرسل للمراجعة"
echo ""

```

---

## 📄 frontend-admin-dashboard/build-obfuscated.js

```javascript
import JavaScriptObfuscator from 'javascript-obfuscator';
import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
import { join } from 'path';

// إعدادات التعتيم
const obfuscatorOptions = {
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.75,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.4,
  debugProtection: true,
  disableConsoleOutput: false, // Changed to false for debugging
  identifierNamesGenerator: 'hexadecimal',
  log: true, // Changed to true for debugging
  numbersToExpressions: true,
  renameGlobals: false,
  selfDefending: true,
  simplify: true,
  splitStrings: true,
  splitStringsChunkLength: 10,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayIndexShift: true,
  stringArrayWrappersCount: 2,
  stringArrayWrappersChainedCalls: true,
  stringArrayWrappersParametersMaxCount: 4,
  stringArrayWrappersType: 'function',
  stringArrayThreshold: 0.75,
  transformObjectKeys: true,
  unicodeEscapeSequence: false
};

// دالة للبحث عن ملفات JS بشكل متكرر
function findJsFiles(dir, fileList = []) {
  try {
    const files = readdirSync(dir);
    
    for (const file of files) {
      const filePath = join(dir, file);
      const stats = statSync(filePath);
      
      if (stats.isDirectory()) {
        findJsFiles(filePath, fileList);
      } else if (file.endsWith('.js')) {
        fileList.push(filePath);
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error.message);
  }
  
  return fileList;
}

// البحث عن جميع ملفات JS في مجلد dist
const jsFiles = findJsFiles('dist');

console.log(`Found ${jsFiles.length} JavaScript files to obfuscate`);

// تعتيم كل ملف
for (const file of jsFiles) {
  try {
    console.log(`Obfuscating: ${file}`);
    const code = readFileSync(file, 'utf8');
    const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, obfuscatorOptions).getObfuscatedCode();
    writeFileSync(file, obfuscatedCode, 'utf8');
    console.log(`✓ Obfuscated: ${file}`);
  } catch (error) {
    console.error(`✗ Error obfuscating ${file}:`, error.message);
  }
}

console.log('Obfuscation complete!');
```

---

## 📄 frontend-admin-dashboard/src/tests/setup.js

```javascript
import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// محاكاة window.matchMedia للاختبارات
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// إعدادات Vue Test Utils
config.global.mocks = {
  $t: (key) => key,
}

// محاكاة localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock
```

---

## 📄 frontend-client-portal/build-obfuscated.js

```javascript
import JavaScriptObfuscator from 'javascript-obfuscator';
import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
import { join } from 'path';

// إعدادات التعتيم
const obfuscatorOptions = {
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.75,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.4,
  debugProtection: true,
  disableConsoleOutput: false, // Changed to false for debugging
  identifierNamesGenerator: 'hexadecimal',
  log: true, // Changed to true for debugging
  numbersToExpressions: true,
  renameGlobals: false,
  selfDefending: true,
  simplify: true,
  splitStrings: true,
  splitStringsChunkLength: 10,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayIndexShift: true,
  stringArrayWrappersCount: 2,
  stringArrayWrappersChainedCalls: true,
  stringArrayWrappersParametersMaxCount: 4,
  stringArrayWrappersType: 'function',
  stringArrayThreshold: 0.75,
  transformObjectKeys: true,
  unicodeEscapeSequence: false
};

// دالة للبحث عن ملفات JS بشكل متكرر
function findJsFiles(dir, fileList = []) {
  try {
    const files = readdirSync(dir);
    
    for (const file of files) {
      const filePath = join(dir, file);
      const stats = statSync(filePath);
      
      if (stats.isDirectory()) {
        findJsFiles(filePath, fileList);
      } else if (file.endsWith('.js')) {
        fileList.push(filePath);
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error.message);
  }
  
  return fileList;
}

// البحث عن جميع ملفات JS في مجلد dist
const jsFiles = findJsFiles('dist');

console.log(`Found ${jsFiles.length} JavaScript files to obfuscate`);

// تعتيم كل ملف
for (const file of jsFiles) {
  try {
    console.log(`Obfuscating: ${file}`);
    const code = readFileSync(file, 'utf8');
    const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, obfuscatorOptions).getObfuscatedCode();
    writeFileSync(file, obfuscatedCode, 'utf8');
    console.log(`✓ Obfuscated: ${file}`);
  } catch (error) {
    console.error(`✗ Error obfuscating ${file}:`, error.message);
  }
}

console.log('Obfuscation complete!');
```

---

## 📄 frontend-client-portal/src/tests/setup.js

```javascript
import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// محاكاة window.matchMedia للاختبارات
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// إعدادات Vue Test Utils
config.global.mocks = {
  $t: (key) => key,
}

// محاكاة localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock
```

---

## 📄 frontend-worker-pwa/build-obfuscated.js

```javascript
import JavaScriptObfuscator from 'javascript-obfuscator';
import { readFileSync, writeFileSync, readdirSync, statSync } from 'fs';
import { join } from 'path';

// إعدادات التعتيم
const obfuscatorOptions = {
  compact: true,
  controlFlowFlattening: true,
  controlFlowFlatteningThreshold: 0.75,
  deadCodeInjection: true,
  deadCodeInjectionThreshold: 0.4,
  debugProtection: true,
  disableConsoleOutput: false, // Changed to false for debugging
  identifierNamesGenerator: 'hexadecimal',
  log: true, // Changed to true for debugging
  numbersToExpressions: true,
  renameGlobals: false,
  selfDefending: true,
  simplify: true,
  splitStrings: true,
  splitStringsChunkLength: 10,
  stringArray: true,
  stringArrayEncoding: ['base64'],
  stringArrayIndexShift: true,
  stringArrayWrappersCount: 2,
  stringArrayWrappersChainedCalls: true,
  stringArrayWrappersParametersMaxCount: 4,
  stringArrayWrappersType: 'function',
  stringArrayThreshold: 0.75,
  transformObjectKeys: true,
  unicodeEscapeSequence: false
};

// دالة للبحث عن ملفات JS بشكل متكرر
function findJsFiles(dir, fileList = []) {
  try {
    const files = readdirSync(dir);
    
    for (const file of files) {
      const filePath = join(dir, file);
      const stats = statSync(filePath);
      
      if (stats.isDirectory()) {
        findJsFiles(filePath, fileList);
      } else if (file.endsWith('.js')) {
        fileList.push(filePath);
      }
    }
  } catch (error) {
    console.error(`Error reading directory ${dir}:`, error.message);
  }
  
  return fileList;
}

// البحث عن جميع ملفات JS في مجلد dist
const jsFiles = findJsFiles('dist');

console.log(`Found ${jsFiles.length} JavaScript files to obfuscate`);

// تعتيم كل ملف
for (const file of jsFiles) {
  try {
    console.log(`Obfuscating: ${file}`);
    const code = readFileSync(file, 'utf8');
    const obfuscatedCode = JavaScriptObfuscator.obfuscate(code, obfuscatorOptions).getObfuscatedCode();
    writeFileSync(file, obfuscatedCode, 'utf8');
    console.log(`✓ Obfuscated: ${file}`);
  } catch (error) {
    console.error(`✗ Error obfuscating ${file}:`, error.message);
  }
}

console.log('Obfuscation complete!');
```

---

## 📄 frontend-worker-pwa/src/tests/setup.js

```javascript
import { vi } from 'vitest'
import { config } from '@vue/test-utils'

// محاكاة window.matchMedia للاختبارات
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// إعدادات Vue Test Utils
config.global.mocks = {
  $t: (key) => key,
}

// محاكاة localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock

// محاكاة navigator إذا لم يكن موجوداً
if (!global.navigator) {
  global.navigator = {}
}

// محاكاة Geolocation API
const geolocationMock = {
  getCurrentPosition: vi.fn((success) => {
    success({
      coords: {
        latitude: 24.7136,
        longitude: 46.6753,
        accuracy: 10,
      },
      timestamp: Date.now(),
    })
  }),
  watchPosition: vi.fn(),
  clearWatch: vi.fn(),
}
global.navigator.geolocation = geolocationMock
```

---

## 📄 setup-codemagic-api.sh

```bash
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

```

---

## 📄 setup-codemagic.sh

```bash
#!/bin/bash

# Script to setup GitHub Secrets for Codemagic integration
# Usage: ./setup-codemagic.sh

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

```

---


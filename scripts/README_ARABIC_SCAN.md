# Arabic Text Scanner Script

## Overview
The `scan_arabic_text.py` script scans source code files for Arabic text and checks if corresponding translations exist in the translation files. It helps identify hardcoded Arabic text that should be moved to translation files for proper internationalization.

With the `--auto-sync` option, it can automatically add missing translation keys and sync them across all supported languages.

## Features
- **Arabic Text Detection**: Uses Unicode patterns to find Arabic text in source code
- **Translation Coverage Check**: Compares found Arabic text with existing translations
- **Smart Filtering**: Excludes comments, console logs, and other non-user-facing text
- **Multiple Output Formats**: Console, JSON, and CSV output options
- **File Statistics**: Shows which files contain the most Arabic text
- **Context Display**: Shows the context around found Arabic text
- **Auto-Sync**: Automatically adds missing translation keys and syncs to all languages
- **Key Generation**: Generates unique translation keys from Arabic text using transliteration

## Why Use This Script?

### Problems It Solves
1. **Hardcoded Arabic Text**: Finds Arabic text directly in Vue/JS files instead of using translation keys
2. **Missing Translations**: Identifies Arabic text that doesn't have corresponding translations in other languages
3. **Inconsistent Internationalization**: Helps ensure all user-facing text goes through the translation system

### Benefits
- **Better Maintainability**: Centralized translations are easier to update
- **Proper Localization**: Ensures all languages have consistent translations
- **Code Quality**: Reduces hardcoded text and improves i18n practices

## Usage

### Basic Usage (All Frontend Directories)
```bash
python3 scripts/scan_arabic_text.py
```

### Scan Specific Directory
```bash
python3 scripts/scan_arabic_text.py --dir frontend-admin-dashboard/src
```

### Output as JSON
```bash
python3 scripts/scan_arabic_text.py --output json > report.json
```

### Output as CSV
```bash
python3 scripts/scan_arabic_text.py --output csv > report.csv
```

### Exclude Additional Patterns
```bash
python3 scripts/scan_arabic_text.py --exclude node_modules dist build custom_folder
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `--dir` | Specific directory to scan | All frontend src directories |
| `--exclude` | Patterns to exclude from scanning | node_modules, dist, build, .git, coverage |
| `--output` | Output format (console, json, csv) | console |

## Understanding the Results

### Report Sections

#### 1. Summary Statistics
```
Total Arabic text findings: 1106
Potentially translated: 251
Likely untranslated: 855
```

- **Total Arabic text findings**: All Arabic text found in source code
- **Potentially translated**: Text that exists in translation files
- **Likely untranslated**: Text that should be moved to translation files

#### 2. Likely Untranslated Arabic Text
Shows Arabic text that likely needs to be moved to translation files:

```
1. File: frontend-admin-dashboard/src/main.js
   Line: 52
   Arabic text: "يتوقف إذا كانت مفتوحة"
   Context: debugger; // يتوقف إذا كانت DevTools مفتوحة
   Similar translation keys found:
     - ar: and
```

#### 3. Potentially Translated Arabic Text
Shows Arabic text that already exists in translation files:

```
1. File: frontend-admin-dashboard/src/services/i18n.js
   Line: 6
   Arabic text: "تسجيل الدخول"
```

#### 4. Statistics by File
Shows which files contain the most Arabic text:

```
frontend-admin-dashboard/src/services/i18n.js: 509 Arabic text(s)
frontend-client-portal/src/plugins/i18n.js: 69 Arabic text(s)
```

## Example Workflow

### Step 1: Run the Scanner
```bash
python3 scripts/scan_arabic_text.py --output json > arabic_report.json
```

### Step 2: Review the Report
Look at the "Likely untranslated" section to identify text that needs translation keys.

### Step 3: Add Translation Keys
For each untranslated text, add a key to the base translation file:

```json
// frontend-admin-dashboard/src/i18n/en.json
{
  "devtools_open": "Stops if DevTools is open",
  "allow_right_click": "Allow right-click for copy/paste"
}
```

### Step 4: Translate to Arabic
Add the Arabic translation:

```json
// frontend-admin-dashboard/src/i18n/ar.json
{
  "devtools_open": "يتوقف إذا كانت DevTools مفتوحة",
  "allow_right_click": "السماح بالنقر الأيمن للنسخ واللصق"
}
```

### Step 5: Update Source Code
Replace hardcoded Arabic text with translation keys:

```javascript
// Before
console.log("يتوقف إذا كانت مفتوحة");

// After
console.log(t('devtools_open'));
```

### Step 6: Sync Translations
Use the translation sync script to propagate to all languages:

```bash
python3 scripts/sync_translations.py
```

## Patterns Automatically Excluded

The script automatically excludes:
- Single-line comments: `// comment`
- Multi-line comments: `/* comment */`
- Console logs: `console.log()`, `console.error()`, etc.
- TODO/FIXME/NOTE comments
- Very short text (less than 2 characters)

## Scanned File Types

The script scans these file extensions:
- `.vue` - Vue components
- `.js` - JavaScript files
- `.jsx` - React JSX files
- `.ts` - TypeScript files
- `.tsx` - TypeScript JSX files
- `.html` - HTML files

## Integration with CI/CD

Add to your CI/CD pipeline to prevent hardcoded Arabic text:

```yaml
# Example GitHub Actions
- name: Check for Hardcoded Arabic Text
  run: |
    python3 scripts/scan_arabic_text.py --output json > report.json
    untranslated_count=$(jq '.summary.untranslated' report.json)
    if [ $untranslated_count -gt 0 ]; then
      echo "Found $untranslated_count untranslated Arabic texts"
      exit 1
    fi
```

## Tips for Best Results

### 1. Focus on User-Facing Text
Prioritize fixing Arabic text in:
- Vue component templates
- User interface messages
- Form labels and buttons
- Error messages

### 2. Ignore Developer Comments
Text in comments (like `// TODO: fix this`) doesn't need translation.

### 3. Use Consistent Key Naming
Follow naming conventions:
```json
{
  "button_submit": "Submit",
  "error_connection": "Connection error",
  "label_email": "Email address"
}
```

### 4. Group Related Translations
Use nested structures for related content:
```json
{
  "login": {
    "title": "Login",
    "email": "Email",
    "password": "Password",
    "button": "Sign In"
  }
}
```

## Troubleshooting

### Too Many False Positives
If the script finds too many non-user-facing text:
1. Add more exclusion patterns to the script
2. Focus on specific directories with `--dir`
3. Manually review the JSON output to filter results

### Missing Translations
If text is marked as untranslated but you know it's translated:
1. The script checks for exact matches in Arabic translation files
2. Similar text might not match exactly
3. Consider standardizing the translation text

### Performance Issues
For large codebases:
1. Scan specific directories instead of all
2. Use `--exclude` to skip test files
3. Run during off-peak hours

## Related Scripts

- **sync_translations.py**: Synchronizes translation keys across all languages
- **Use together**: First fix hardcoded text, then sync translations

## Example Output Analysis

### High Priority Files
Files with many untranslated texts are high priority:
```
frontend-admin-dashboard/src/services/i18n.js: 509 Arabic text(s)
```
This file contains translation definitions, which is expected.

### Unexpected Locations
Arabic text in unexpected locations needs attention:
```
frontend-admin-dashboard/src/main.js: 3 Arabic text(s)
```
Main.js should not contain user-facing Arabic text.

## Best Practices

### 1. Prevent Hardcoded Text
- Never hardcode user-facing text in components
- Always use translation keys
- Create translation keys first, then implement

### 2. Regular Scanning
- Run the scanner before releases
- Add to pre-commit hooks
- Monitor for new hardcoded text

### 3. Code Review
- Check for hardcoded text in PRs
- Use the scanner as part of review process
- Educate team on i18n best practices

## Contributing

When adding new features:
1. Use translation keys for all user-facing text
2. Run the scanner to verify no hardcoded text
3. Sync translations to all languages
4. Document any new translation keys

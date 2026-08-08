# Translation Sync Script

## Overview
The `sync_translations.py` script automatically synchronizes translation keys across all language files in the project. It ensures that all language files have the same keys as the base language (English), copying missing keys from the base language while preserving existing translations.

## Features
- **Automatic Detection**: Finds all translation directories in the project
- **Deep Merge**: Preserves existing translations while adding missing keys
- **Nested Support**: Handles nested JSON structures
- **Dry Run Mode**: Preview changes before applying them
- **Multi-Directory**: Processes all frontend apps (admin, client, worker)

## Supported Languages
The script supports the following languages:
- `ar` - Arabic
- `he` - Hebrew  
- `en` - English (base language)
- `fr` - French
- `de` - German
- `es` - Spanish
- `it` - Italian
- `pt` - Portuguese
- `ru` - Russian
- `zh` - Chinese
- `ja` - Japanese
- `ko` - Korean
- `tr` - Turkish
- `nl` - Dutch
- `sv` - Swedish
- `pl` - Polish

## Usage

### Basic Usage
```bash
python3 scripts/sync_translations.py
```

### Dry Run (Preview Changes)
```bash
python3 scripts/sync_translations.py --dry-run
```

### Specific Language Directory
```bash
python3 scripts/sync_translations.py --lang-dir frontend-admin-dashboard/src/i18n
```

### Custom Base Language
```bash
python3 scripts/sync_translations.py --base-lang ar
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `--base-lang` | Base language to sync from | `en` |
| `--lang-dir` | Specific language directory to process | Auto-detects all |
| `--dry-run` | Show changes without modifying files | `False` |

## How It Works

1. **Load Base Language**: Reads the base language file (e.g., `en.json`)
2. **Scan Directories**: Finds all translation directories in the project
3. **Compare Keys**: For each language file, compares keys with the base
4. **Deep Merge**: Adds missing keys from base while keeping existing translations
5. **Save Files**: Updates the language files with merged content

## Example Output

```
🌍 Translation Sync Script
============================================================
Base language: en
Supported languages: ar, he, en, fr, de, es, it, pt, ru, zh, ja, ko, tr, nl, sv, pl
Dry run: False

Found 6 translation directory(s)

📁 Processing: /home/dev-bit/project/worktrack/frontend-admin-dashboard/src/i18n
============================================================
✅ ar: No missing keys
🔧 he: Found 1 missing keys
   - done
✅ he: Updated successfully
✅ fr: No missing keys
...
📊 Summary: 15/15 files updated

============================================================
✅ Sync complete! 90 files updated across 6 directories
```

## Workflow Recommendations

### Adding New Translation Keys
1. Add the new key to the base language file (`en.json`)
2. Run the sync script: `python3 scripts/sync_translations.py`
3. The script will automatically add the key to all other language files
4. Translate the new key in each language file manually

### Before Major Releases
1. Run dry-run first: `python3 scripts/sync_translations.py --dry-run`
2. Review the planned changes
3. Run the actual sync: `python3 scripts/sync_translations.py`
4. Test the application in different languages

### Continuous Integration
Add the script to your CI/CD pipeline to ensure translation keys are always synchronized:

```yaml
# Example GitHub Actions
- name: Sync Translations
  run: python3 scripts/sync_translations.py --dry-run
  
- name: Check for Translation Changes
  run: |
    if git diff --name-only | grep -q "i18n"; then
      echo "Translation files changed. Please sync translations."
      exit 1
    fi
```

## Translation Directories

The script automatically searches for translation directories in:
- `frontend-admin-dashboard/src/i18n/`
- `frontend-admin-dashboard/public/i18n/`
- `frontend-client-portal/src/i18n/`
- `frontend-client-portal/public/i18n/`
- `frontend-worker-pwa/src/i18n/`
- `frontend-worker-pwa/public/i18n/`

## Notes

- The script preserves existing translations - it only adds missing keys
- Nested JSON structures are supported (e.g., `"pwa": { "installTitle": "..." }`)
- The base language file must exist in each directory
- UTF-8 encoding is used for all files
- JSON files are formatted with 2-space indentation

## Troubleshooting

### "File not found" Error
Ensure the base language file exists in the translation directory.

### "Error parsing JSON"
Check that the JSON file is valid using a JSON validator.

### Too many missing keys
This is normal for newly added languages. The keys will be in English initially and can be translated later.

## Contributing

When adding new features that require translations:
1. Always add the keys to the base language file first
2. Run the sync script to propagate to all languages
3. Document any new keys in the translation guide
4. Consider using translation services for the new keys

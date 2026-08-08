#!/usr/bin/env python3
"""
Translation Keys Extractor Script
==================================
This script scans source code files for translation key usage (t('key')) and 
ensures all keys exist in translation files. It adds missing keys to all languages.

Usage:
    python scripts/extract_translation_keys.py

Options:
    --dir: Directory to scan (default: all frontend directories)
    --dry-run: Show changes without modifying files
"""

import re
import json
import os
import sys
import argparse
from pathlib import Path
from typing import Dict, List, Set
from collections import defaultdict


class TranslationKeysExtractor:
    def __init__(self, target_dir: str = None, dry_run: bool = False):
        self.target_dir = target_dir
        self.dry_run = dry_run
        
        # Pattern to match t('key') or t("key")
        # Supports nested keys like t('pwa.installTitle')
        self.t_pattern = re.compile(r"t\(['\"]([a-zA-Z_][a-zA-Z0-9_\.]*)['\"]\)")
        
        # Patterns to exclude (URLs, paths, special characters)
        self.exclude_patterns = [
            r'^[\/\#\,]',  # Starts with /, #, or ,
            r'^http',  # URLs
            r'^\/[a-z]',  # API paths
            r'^[0-9]',  # Starts with number
        ]
        
        # File extensions to scan
        self.scan_extensions = {'.vue', '.js', '.jsx', '.ts', '.tsx'}
        
        # Translation directories
        self.translation_dirs = [
            'frontend-admin-dashboard/src/i18n',
            'frontend-admin-dashboard/public/i18n',
            'frontend-client-portal/src/i18n',
            'frontend-client-portal/public/i18n',
            'frontend-worker-pwa/src/i18n',
            'frontend-worker-pwa/public/i18n',
        ]
        
        # Supported languages
        self.supported_langs = ['ar', 'he', 'en', 'fr', 'de', 'es', 'it', 'pt', 'ru', 
                                'zh', 'ja', 'ko', 'tr', 'nl', 'sv', 'pl']
        
        # Results
        self.found_keys = set()
        self.translations = defaultdict(dict)
    
    def extract_keys_from_file(self, file_path: Path) -> Set[str]:
        """Extract translation keys from a file"""
        keys = set()
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Find all t('key') patterns
            matches = self.t_pattern.findall(content)
            
            # Filter out excluded patterns
            for key in matches:
                should_exclude = False
                for pattern in self.exclude_patterns:
                    if re.match(pattern, key):
                        should_exclude = True
                        break
                
                if not should_exclude:
                    keys.add(key)
        
        except Exception as e:
            print(f"⚠️  Error reading {file_path}: {e}")
        
        return keys
    
    def scan_directory(self, directory: Path) -> Set[str]:
        """Scan a directory for translation keys"""
        keys = set()
        
        # Convert to absolute path if relative
        if not directory.is_absolute():
            directory = Path.cwd() / directory
        
        if not directory.exists():
            print(f"⚠️  Directory not found: {directory}")
            return keys
        
        print(f"📁 Scanning: {directory}")
        
        for file_path in directory.rglob('*'):
            if file_path.suffix in self.scan_extensions and file_path.is_file():
                file_keys = self.extract_keys_from_file(file_path)
                keys.update(file_keys)
        
        return keys
    
    def load_translations(self):
        """Load all translation files"""
        base_dir = Path.cwd()
        
        for trans_dir in self.translation_dirs:
            trans_path = base_dir / trans_dir
            if not trans_path.exists():
                continue
            
            print(f"📖 Loading translations from: {trans_dir}")
            
            for lang_file in trans_path.glob('*.json'):
                lang_code = lang_file.stem
                
                try:
                    with open(lang_file, 'r', encoding='utf-8') as f:
                        translations = json.load(f)
                        
                    # Extract all translation keys (including nested)
                    flat_translations = self.flatten_json(translations)
                    self.translations[lang_code].update(flat_translations)
                    
                except Exception as e:
                    print(f"⚠️  Error loading {lang_file}: {e}")
    
    def flatten_json(self, data: Dict, prefix: str = '') -> Dict:
        """Flatten nested JSON to dot-separated keys"""
        result = {}
        
        for key, value in data.items():
            new_key = f"{prefix}.{key}" if prefix else key
            
            if isinstance(value, dict):
                result.update(self.flatten_json(value, new_key))
            elif isinstance(value, str):
                result[new_key] = value
            # Skip non-string values
        
        return result
    
    def unflatten_json(self, flat_data: Dict) -> Dict:
        """Convert flat dot-separated keys back to nested structure"""
        result = {}
        
        for key, value in flat_data.items():
            parts = key.split('.')
            current = result
            
            for part in parts[:-1]:
                if part not in current:
                    current[part] = {}
                current = current[part]
            
            current[parts[-1]] = value
        
        return result
    
    def find_missing_keys(self) -> Dict[str, Set[str]]:
        """Find keys used in code but missing from translations"""
        missing_by_lang = defaultdict(set)
        
        # Check English as base
        en_keys = set(self.translations.get('en', {}).keys())
        
        for key in self.found_keys:
            if key not in en_keys:
                # Key is missing from English (and likely all languages)
                for lang in self.supported_langs:
                    missing_by_lang[lang].add(key)
        
        return missing_by_lang
    
    def add_missing_keys(self, missing_keys: Dict[str, Set[str]]) -> int:
        """Add missing keys to all translation files"""
        if not missing_keys or not any(missing_keys.values()):
            print("ℹ️  No missing keys to add")
            return 0
        
        total_added = 0
        base_dir = Path.cwd()
        
        for trans_dir in self.translation_dirs:
            trans_path = base_dir / trans_dir
            if not trans_path.exists():
                continue
            
            print(f"📝 Updating: {trans_dir}")
            
            for lang_code in self.supported_langs:
                if lang_code not in missing_keys or not missing_keys[lang_code]:
                    continue
                
                lang_file = trans_path / f'{lang_code}.json'
                if not lang_file.exists():
                    continue
                
                try:
                    # Load existing translations
                    with open(lang_file, 'r', encoding='utf-8') as f:
                        translations = json.load(f)
                    
                    # Add missing keys
                    added_count = 0
                    for key in missing_keys[lang_code]:
                        # Handle nested keys
                        if '.' in key:
                            parts = key.split('.')
                            current = translations
                            for part in parts[:-1]:
                                if part not in current:
                                    current[part] = {}
                                current = current[part]
                            
                            if parts[-1] not in current:
                                current[parts[-1]] = f"[TODO: Translate {key}]"
                                added_count += 1
                        else:
                            if key not in translations:
                                translations[key] = f"[TODO: Translate {key}]"
                                added_count += 1
                    
                    if added_count > 0:
                        if self.dry_run:
                            print(f"  [DRY RUN] Would add {added_count} keys to {lang_code}.json")
                        else:
                            # Save updated translations
                            with open(lang_file, 'w', encoding='utf-8') as f:
                                json.dump(translations, f, ensure_ascii=False, indent=2)
                            print(f"  ✅ Added {added_count} keys to {lang_code}.json")
                            total_added += added_count
                
                except Exception as e:
                    print(f"  ❌ Error updating {lang_file}: {e}")
        
        return total_added
    
    def generate_report(self, missing_keys: Dict[str, Set[str]]):
        """Generate report of findings"""
        print("\n" + "=" * 80)
        print("📊 TRANSLATION KEYS EXTRACTION REPORT")
        print("=" * 80)
        
        print(f"\nTotal keys found in code: {len(self.found_keys)}")
        print(f"Total keys in English translations: {len(self.translations.get('en', {}))}")
        
        # Find missing keys
        en_keys = set(self.translations.get('en', {}).keys())
        missing_from_en = self.found_keys - en_keys
        
        print(f"Missing from English: {len(missing_from_en)}")
        
        if missing_from_en:
            print(f"\n⚠️  MISSING KEYS (will be added to all languages):")
            for i, key in enumerate(sorted(missing_from_en)[:20], 1):
                print(f"  {i}. {key}")
            if len(missing_from_en) > 20:
                print(f"  ... and {len(missing_from_en) - 20} more")
        
        # Show keys in translations but not in code
        not_used_in_code = set(self.translations.get('en', {}).keys()) - self.found_keys
        print(f"\n📋 Keys in translations but not used in code: {len(not_used_in_code)}")
        
        if not_used_in_code:
            print(f"  (These might be unused keys)")
            for i, key in enumerate(sorted(not_used_in_code)[:10], 1):
                print(f"  {i}. {key}")
            if len(not_used_in_code) > 10:
                print(f"  ... and {len(not_used_in_code) - 10} more")
    
    def run(self):
        """Main execution"""
        print("🔍 Translation Keys Extractor")
        print("=" * 80)
        
        # Load translations first
        print("\n📖 Loading translation files...")
        self.load_translations()
        
        total_translations = sum(len(trans) for trans in self.translations.values())
        print(f"✅ Loaded {len(self.translations)} language files with {total_translations} total translations")
        
        # Determine directories to scan
        if self.target_dir:
            directories = [Path(self.target_dir)]
        else:
            base_dir = Path.cwd()
            directories = [
                base_dir / 'frontend-admin-dashboard' / 'src',
                base_dir / 'frontend-client-portal' / 'src',
                base_dir / 'frontend-worker-pwa' / 'src',
            ]
        
        # Scan directories
        all_keys = set()
        for directory in directories:
            if directory.exists():
                keys = self.scan_directory(directory)
                all_keys.update(keys)
        
        self.found_keys = all_keys
        print(f"\n✅ Found {len(all_keys)} unique translation keys in code")
        
        # Find missing keys
        missing_keys = self.find_missing_keys()
        
        # Generate report
        self.generate_report(missing_keys)
        
        # Add missing keys
        if missing_keys and any(missing_keys.values()):
            if self.dry_run:
                print("\n📋 [DRY RUN] No files will be modified")
            else:
                added = self.add_missing_keys(missing_keys)
                if added > 0:
                    print(f"\n✅ Successfully added {added} translation keys")
        else:
            print("\n✅ All translation keys are already present in translation files")
        
        return 0


def main():
    parser = argparse.ArgumentParser(
        description='Extract translation keys from source code and add missing ones to translation files'
    )
    parser.add_argument(
        '--dir',
        help='Specific directory to scan (default: all frontend src directories)'
    )
    parser.add_argument(
        '--dry-run',
        action='store_true',
        help='Show changes without modifying files'
    )
    
    args = parser.parse_args()
    
    extractor = TranslationKeysExtractor(
        target_dir=args.dir,
        dry_run=args.dry_run
    )
    
    sys.exit(extractor.run())


if __name__ == '__main__':
    main()

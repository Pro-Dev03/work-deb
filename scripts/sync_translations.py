#!/usr/bin/env python3
"""
Translation Sync Script
========================
This script synchronizes translation keys across all language files.
It reads the base language file (English) and ensures all other language files
have the same keys, copying missing keys from the base language.

Usage:
    python scripts/sync_translations.py

Options:
    --base-lang: Base language to sync from (default: en)
    --lang-dir: Directory containing language files (default: auto-detected)
    --dry-run: Show changes without modifying files
"""

import json
import os
import sys
import argparse
from pathlib import Path
from typing import Dict, Any
from copy import deepcopy


class TranslationSyncer:
    def __init__(self, base_lang: str = 'en', lang_dir: str = None, dry_run: bool = False):
        self.base_lang = base_lang
        self.lang_dir = lang_dir
        self.dry_run = dry_run
        self.supported_langs = [
            'ar', 'he', 'en', 'fr', 'de', 'es', 'it', 'pt', 'ru', 
            'zh', 'ja', 'ko', 'tr', 'nl', 'sv', 'pl'
        ]
        
    def find_translation_directories(self):
        """Find all translation directories in the project"""
        base_dir = Path.cwd()
        possible_dirs = [
            base_dir / 'frontend-admin-dashboard' / 'src' / 'i18n',
            base_dir / 'frontend-admin-dashboard' / 'public' / 'i18n',
            base_dir / 'frontend-client-portal' / 'src' / 'i18n',
            base_dir / 'frontend-client-portal' / 'public' / 'i18n',
            base_dir / 'frontend-worker-pwa' / 'src' / 'i18n',
            base_dir / 'frontend-worker-pwa' / 'public' / 'i18n',
        ]
        
        found_dirs = [d for d in possible_dirs if d.exists()]
        return found_dirs
    
    def load_json_file(self, file_path: Path) -> Dict[str, Any]:
        """Load JSON file safely"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                return json.load(f)
        except FileNotFoundError:
            print(f"⚠️  File not found: {file_path}")
            return {}
        except json.JSONDecodeError as e:
            print(f"❌ Error parsing {file_path}: {e}")
            return {}
    
    def save_json_file(self, file_path: Path, data: Dict[str, Any]):
        """Save JSON file with proper formatting"""
        try:
            with open(file_path, 'w', encoding='utf-8') as f:
                json.dump(data, f, ensure_ascii=False, indent=2)
            return True
        except Exception as e:
            print(f"❌ Error saving {file_path}: {e}")
            return False
    
    def deep_merge(self, base: Dict[str, Any], target: Dict[str, Any]) -> Dict[str, Any]:
        """
        Deep merge dictionaries: keep existing translations, add missing keys from base
        """
        result = deepcopy(target)
        
        for key, value in base.items():
            if key not in result:
                # Key missing in target, copy from base
                result[key] = deepcopy(value)
            elif isinstance(value, dict) and isinstance(result[key], dict):
                # Both are dicts, recurse
                result[key] = self.deep_merge(value, result[key])
            # else: key exists in target, keep existing translation
        
        return result
    
    def get_missing_keys(self, base: Dict[str, Any], target: Dict[str, Any], path: str = '') -> list:
        """Get list of missing keys in target compared to base"""
        missing = []
        
        for key, value in base.items():
            current_path = f"{path}.{key}" if path else key
            
            if key not in target:
                missing.append(current_path)
            elif isinstance(value, dict) and isinstance(target.get(key), dict):
                missing.extend(self.get_missing_keys(value, target[key], current_path))
        
        return missing
    
    def sync_language_file(self, lang_dir: Path, lang_code: str):
        """Sync a single language file"""
        base_file = lang_dir / f'{self.base_lang}.json'
        target_file = lang_dir / f'{lang_code}.json'
        
        if not base_file.exists():
            print(f"⚠️  Base file not found: {base_file}")
            return False
        
        # Load base translations
        base_data = self.load_json_file(base_file)
        if not base_data:
            print(f"❌ Failed to load base file: {base_file}")
            return False
        
        # Load target translations (or create empty if doesn't exist)
        target_data = self.load_json_file(target_file) if target_file.exists() else {}
        
        # Find missing keys
        missing_keys = self.get_missing_keys(base_data, target_data)
        
        if not missing_keys:
            print(f"✅ {lang_code}: No missing keys")
            return True
        
        print(f"🔧 {lang_code}: Found {len(missing_keys)} missing keys")
        for key in missing_keys[:5]:  # Show first 5
            print(f"   - {key}")
        if len(missing_keys) > 5:
            print(f"   ... and {len(missing_keys) - 5} more")
        
        # Merge translations
        merged_data = self.deep_merge(base_data, target_data)
        
        if self.dry_run:
            print(f"📋 [DRY RUN] Would update: {target_file}")
            return True
        
        # Save merged data
        if self.save_json_file(target_file, merged_data):
            print(f"✅ {lang_code}: Updated successfully")
            return True
        else:
            print(f"❌ {lang_code}: Failed to update")
            return False
    
    def sync_directory(self, lang_dir: Path):
        """Sync all language files in a directory"""
        print(f"\n📁 Processing: {lang_dir}")
        print("=" * 60)
        
        success_count = 0
        for lang_code in self.supported_langs:
            if lang_code == self.base_lang:
                continue  # Skip base language
            
            if self.sync_language_file(lang_dir, lang_code):
                success_count += 1
        
        print(f"\n📊 Summary: {success_count}/{len(self.supported_langs)-1} files updated")
        return success_count
    
    def run(self):
        """Main execution"""
        print("🌍 Translation Sync Script")
        print("=" * 60)
        print(f"Base language: {self.base_lang}")
        print(f"Supported languages: {', '.join(self.supported_langs)}")
        print(f"Dry run: {self.dry_run}")
        print()
        
        # Find translation directories
        if self.lang_dir:
            lang_dirs = [Path(self.lang_dir)]
        else:
            lang_dirs = self.find_translation_directories()
        
        if not lang_dirs:
            print("❌ No translation directories found")
            return 1
        
        print(f"Found {len(lang_dirs)} translation directory(s)")
        
        # Process each directory
        total_success = 0
        for lang_dir in lang_dirs:
            total_success += self.sync_directory(lang_dir)
        
        print("\n" + "=" * 60)
        print(f"✅ Sync complete! {total_success} files updated across {len(lang_dirs)} directories")
        
        return 0


def main():
    parser = argparse.ArgumentParser(
        description='Sync translation keys across all language files'
    )
    parser.add_argument(
        '--base-lang',
        default='en',
        help='Base language to sync from (default: en)'
    )
    parser.add_argument(
        '--lang-dir',
        help='Specific language directory to process (default: all found directories)'
    )
    parser.add_argument(
        '--dry-run',
        action='store_true',
        help='Show changes without modifying files'
    )
    
    args = parser.parse_args()
    
    syncer = TranslationSyncer(
        base_lang=args.base_lang,
        lang_dir=args.lang_dir,
        dry_run=args.dry_run
    )
    
    sys.exit(syncer.run())


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
"""
Arabic Text Scanner Script
===========================
This script scans source code files for Arabic text and checks if corresponding
translations exist in the translation files. It helps identify hardcoded Arabic
text that should be moved to translation files.

With --auto-sync option, it can automatically add missing translation keys
and sync them across all languages.

Usage:
    python scripts/scan_arabic_text.py

Options:
    --dir: Directory to scan (default: all frontend directories)
    --exclude: Patterns to exclude (default: node_modules, dist, build)
    --output: Output format (console, json, csv) (default: console)
    --auto-sync: Automatically add missing translation keys and sync
    --key-prefix: Prefix for auto-generated translation keys (default: "auto_")
"""

import re
import json
import os
import sys
import argparse
from pathlib import Path
from typing import Dict, List, Set, Tuple
from collections import defaultdict


class ArabicTextScanner:
    def __init__(self, target_dir: str = None, exclude_patterns: List[str] = None, 
                 output_format: str = 'console', auto_sync: bool = False, key_prefix: str = 'auto_'):
        self.target_dir = target_dir
        self.exclude_patterns = exclude_patterns or ['node_modules', 'dist', 'build', '.git', 'coverage']
        self.output_format = output_format
        self.auto_sync = auto_sync
        self.key_prefix = key_prefix
        
        # Arabic Unicode range
        self.arabic_pattern = re.compile(r'[\u0600-\u06FF\u0750-\u077F\u08A0-\u08FF\uFB50-\uFDFF\uFE70-\uFEFF]+')
        
        # Patterns to exclude (code comments, strings that are likely not user-facing)
        self.exclude_string_patterns = [
            r'^\/\/.*',  # Single line comments
            r'^\/\*.*\*\/$',  # Multi-line comments
            r'^\s*\/\/.*',  # Indented comments
            r'^\s*\/\*.*\*\/\s*$',  # Indented multi-line comments
            r'console\.',  # Console logs
            r'\/\/\s*TODO',  # TODO comments
            r'\/\/\s*FIXME',  # FIXME comments
            r'\/\/\s*NOTE',  # NOTE comments
        ]
        
        # File extensions to scan
        self.scan_extensions = {'.vue', '.js', '.jsx', '.ts', '.tsx', '.html'}
        
        # Translation directories
        self.translation_dirs = [
            'frontend-admin-dashboard/src/i18n',
            'frontend-admin-dashboard/public/i18n',
            'frontend-client-portal/src/i18n',
            'frontend-client-portal/public/i18n',
            'frontend-worker-pwa/src/i18n',
            'frontend-worker-pwa/public/i18n',
        ]
        
        # Results storage
        self.scan_results = defaultdict(list)
        self.translations = defaultdict(dict)
        
    def is_excluded(self, path: Path) -> bool:
        """Check if path should be excluded"""
        for pattern in self.exclude_patterns:
            if pattern in str(path):
                return True
        return False
    
    def extract_arabic_from_file(self, file_path: Path) -> List[Dict]:
        """Extract Arabic text from a file"""
        results = []
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
                lines = content.split('\n')
                
            for line_num, line in enumerate(lines, 1):
                # Skip lines that match exclusion patterns
                if any(re.match(pattern, line.strip()) for pattern in self.exclude_string_patterns):
                    continue
                
                # Find Arabic text in the line
                arabic_matches = self.arabic_pattern.findall(line)
                if arabic_matches:
                    arabic_text = ' '.join(arabic_matches)
                    
                    # Clean up the text
                    arabic_text = arabic_text.strip()
                    if len(arabic_text) < 2:  # Skip very short matches
                        continue
                    
                    # Get relative path safely
                    try:
                        rel_path = file_path.relative_to(Path.cwd())
                    except ValueError:
                        rel_path = file_path
                    
                    results.append({
                        'file': str(rel_path),
                        'line': line_num,
                        'text': arabic_text,
                        'context': line.strip()[:100]  # First 100 chars for context
                    })
        
        except Exception as e:
            print(f"⚠️  Error reading {file_path}: {e}")
        
        return results
    
    def scan_directory(self, directory: Path) -> List[Dict]:
        """Scan a directory for Arabic text"""
        results = []
        
        # Convert to absolute path if relative
        if not directory.is_absolute():
            directory = Path.cwd() / directory
        
        if not directory.exists():
            print(f"⚠️  Directory not found: {directory}")
            return results
        
        print(f"📁 Scanning: {directory}")
        
        for file_path in directory.rglob('*'):
            if self.is_excluded(file_path):
                continue
            
            if file_path.suffix in self.scan_extensions and file_path.is_file():
                file_results = self.extract_arabic_from_file(file_path)
                results.extend(file_results)
        
        return results
    
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
                        
                    # Extract all translation values (including nested)
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
    
    def generate_translation_key(self, arabic_text: str, existing_keys: Set[str]) -> str:
        """Generate a unique translation key from Arabic text"""
        # Remove special characters and spaces
        clean_text = re.sub(r'[^\w\s]', '', arabic_text)
        clean_text = re.sub(r'\s+', '_', clean_text.strip())
        
        # Transliterate Arabic to Latin characters for key
        transliteration_map = {
            'ا': 'a', 'أ': 'a', 'إ': 'i', 'آ': 'aa',
            'ب': 'b', 'ت': 't', 'ث': 'th', 'ج': 'j',
            'ح': 'h', 'خ': 'kh', 'د': 'd', 'ذ': 'dh',
            'ر': 'r', 'ز': 'z', 'س': 's', 'ش': 'sh',
            'ص': 's', 'ض': 'd', 'ط': 't', 'ظ': 'z',
            'ع': 'a', 'غ': 'gh', 'ف': 'f', 'ق': 'q',
            'ك': 'k', 'ل': 'l', 'م': 'm', 'ن': 'n',
            'ه': 'h', 'و': 'w', 'ي': 'y', 'ة': 'a',
        }
        
        transliterated = ''
        for char in clean_text:
            if char in transliteration_map:
                transliterated += transliteration_map[char]
            elif char.isalnum():
                transliterated += char
            elif char == '_':
                transliterated += '_'
        
        # Limit length
        if len(transliterated) > 50:
            transliterated = transliterated[:50]
        
        # Generate unique key
        base_key = f"{self.key_prefix}{transliterated}"
        counter = 1
        final_key = base_key
        
        while final_key in existing_keys:
            final_key = f"{base_key}_{counter}"
            counter += 1
        
        return final_key
    
    def add_missing_translations(self, untranslated_texts: List[Dict]) -> int:
        """Add missing translation keys to all language files"""
        if not untranslated_texts:
            print("ℹ️  No untranslated texts to add")
            return 0
        
        print(f"\n🔧 Auto-sync: Adding {len(untranslated_texts)} missing translation keys...")
        
        # Collect all existing keys across all languages
        all_existing_keys = set()
        for lang_translations in self.translations.values():
            all_existing_keys.update(lang_translations.keys())
        
        # Generate keys for untranslated texts
        new_entries = defaultdict(dict)
        
        for item in untranslated_texts:
            arabic_text = item['text']
            
            # Skip if this exact text already exists as a value in any language
            already_exists = any(
                arabic_text in translations.values() 
                for translations in self.translations.values()
            )
            if already_exists:
                continue
            
            # Generate unique key
            key = self.generate_translation_key(arabic_text, all_existing_keys)
            all_existing_keys.add(key)
            
            # Add to Arabic with the original text
            new_entries['ar'][key] = arabic_text
            
            # Add to English with placeholder (will be translated later)
            new_entries['en'][key] = f"[TODO: Translate '{arabic_text}']"
        
        if not new_entries['ar']:
            print("ℹ️  All texts already exist in translations")
            return 0
        
        # Update translation files (only ar and en)
        updated_count = 0
        base_dir = Path.cwd()
        
        for trans_dir in self.translation_dirs:
            trans_path = base_dir / trans_dir
            if not trans_path.exists():
                continue
            
            print(f"📝 Updating: {trans_dir}")
            
            for lang_code in ['ar', 'en']:
                lang_file = trans_path / f'{lang_code}.json'
                if not lang_file.exists():
                    continue
                
                try:
                    # Load existing translations
                    with open(lang_file, 'r', encoding='utf-8') as f:
                        translations = json.load(f)
                    
                    # Add new entries
                    if lang_code in new_entries:
                        for key, value in new_entries[lang_code].items():
                            # Check if key already exists (nested check)
                            if '.' not in key and key not in translations:
                                translations[key] = value
                            elif '.' in key:
                                # Handle nested keys
                                parts = key.split('.')
                                current = translations
                                for part in parts[:-1]:
                                    if part not in current:
                                        current[part] = {}
                                    current = current[part]
                                current[parts[-1]] = value
                            elif key not in translations:
                                translations[key] = value
                        
                        # Save updated translations
                        with open(lang_file, 'w', encoding='utf-8') as f:
                            json.dump(translations, f, ensure_ascii=False, indent=2)
                        
                        updated_count += 1
                        print(f"  ✅ Added {len(new_entries[lang_code])} keys to {lang_code}.json")
                
                except Exception as e:
                    print(f"  ❌ Error updating {lang_file}: {e}")
        
        # Now sync to all other languages using the sync script
        if updated_count > 0:
            print(f"\n🔄 Syncing to all other languages...")
            self.sync_all_languages()
        
        return updated_count
    
    def sync_all_languages(self):
        """Sync translations across all languages using sync_translations logic"""
        try:
            import subprocess
            result = subprocess.run(
                ['python3', 'scripts/sync_translations.py'],
                cwd=Path.cwd(),
                capture_output=True,
                text=True
            )
            
            if result.returncode == 0:
                print("✅ Successfully synced to all languages")
            else:
                print(f"⚠️  Sync had issues: {result.stderr}")
        except Exception as e:
            print(f"⚠️  Could not run sync script: {e}")
    
    def find_similar_translations(self, arabic_text: str, lang_translations: Dict) -> List[str]:
        """Find similar translations in other languages"""
        similar = []
        
        for key, translation in lang_translations.items():
            # Check if the Arabic text is very similar to a translation value
            # This helps identify if the text might be a translation key
            if arabic_text.lower() in translation.lower() or translation.lower() in arabic_text.lower():
                similar.append(key)
        
        return similar
    
    def check_translation_coverage(self, arabic_text: str) -> Dict:
        """Check if Arabic text has translations in other languages"""
        coverage = {
            'arabic_text': arabic_text,
            'has_translation': False,
            'languages': {},
            'similar_keys': {}
        }
        
        # Check if the exact Arabic text exists in Arabic translation files
        for lang_code, translations in self.translations.items():
            if lang_code == 'ar':
                if arabic_text in translations.values():
                    coverage['has_translation'] = True
            
            # Check for similar keys in other languages
            similar_keys = self.find_similar_translations(arabic_text, translations)
            if similar_keys:
                coverage['similar_keys'][lang_code] = similar_keys
        
        return coverage
    
    def analyze_results(self, scan_results: List[Dict]) -> List[Dict]:
        """Analyze scan results against translations"""
        analyzed = []
        
        print(f"\n🔍 Analyzing {len(scan_results)} Arabic text findings...")
        
        for result in scan_results:
            arabic_text = result['text']
            coverage = self.check_translation_coverage(arabic_text)
            
            result['analysis'] = coverage
            analyzed.append(result)
        
        return analyzed
    
    def generate_report(self, analyzed_results: List[Dict]):
        """Generate and display the report"""
        total_arabic = len(analyzed_results)
        untranslated = [r for r in analyzed_results if not r['analysis']['has_translation']]
        potentially_translated = [r for r in analyzed_results if r['analysis']['has_translation']]
        
        print("\n" + "=" * 80)
        print("📊 ARABIC TEXT SCAN REPORT")
        print("=" * 80)
        print(f"\nTotal Arabic text findings: {total_arabic}")
        print(f"Potentially translated: {len(potentially_translated)}")
        print(f"Likely untranslated: {len(untranslated)}")
        
        if self.output_format == 'console':
            self.print_console_report(analyzed_results, untranslated, potentially_translated)
        elif self.output_format == 'json':
            self.print_json_report(analyzed_results)
        elif self.output_format == 'csv':
            self.print_csv_report(analyzed_results)
    
    def print_console_report(self, analyzed_results: List[Dict], untranslated: List[Dict], potentially_translated: List[Dict]):
        """Print report in console format"""
        
        if untranslated:
            print("\n" + "⚠️  LIKELY UNTRANSLATED ARABIC TEXT " + "⚠️")
            print("-" * 80)
            
            for i, result in enumerate(untranslated[:20], 1):  # Show first 20
                print(f"\n{i}. File: {result['file']}")
                print(f"   Line: {result['line']}")
                print(f"   Arabic text: \"{result['text']}\"")
                print(f"   Context: {result['context']}")
                
                if result['analysis']['similar_keys']:
                    print(f"   Similar translation keys found:")
                    for lang, keys in result['analysis']['similar_keys'].items():
                        print(f"     - {lang}: {', '.join(keys[:3])}")
            
            if len(untranslated) > 20:
                print(f"\n... and {len(untranslated) - 20} more untranslated texts")
        
        if potentially_translated:
            print("\n" + "✅ POTENTIALLY TRANSLATED ARABIC TEXT " + "✅")
            print("-" * 80)
            
            for i, result in enumerate(potentially_translated[:10], 1):  # Show first 10
                print(f"\n{i}. File: {result['file']}")
                print(f"   Line: {result['line']}")
                print(f"   Arabic text: \"{result['text']}\"")
        
        # Statistics by file
        print("\n" + "📈 STATISTICS BY FILE " + "📈")
        print("-" * 80)
        
        file_stats = defaultdict(int)
        for result in analyzed_results:
            file_stats[result['file']] += 1
        
        sorted_files = sorted(file_stats.items(), key=lambda x: x[1], reverse=True)
        
        for file_path, count in sorted_files[:10]:
            print(f"{file_path}: {count} Arabic text(s)")
    
    def print_json_report(self, analyzed_results: List[Dict]):
        """Print report in JSON format"""
        report = {
            'summary': {
                'total_findings': len(analyzed_results),
                'untranslated': len([r for r in analyzed_results if not r['analysis']['has_translation']]),
                'translated': len([r for r in analyzed_results if r['analysis']['has_translation']])
            },
            'findings': analyzed_results
        }
        
        print(json.dumps(report, ensure_ascii=False, indent=2))
    
    def print_csv_report(self, analyzed_results: List[Dict]):
        """Print report in CSV format"""
        print("file,line,arabic_text,context,has_translation,similar_keys")
        
        for result in analyzed_results:
            similar_keys = json.dumps(result['analysis']['similar_keys']) if result['analysis']['similar_keys'] else ''
            print(f"{result['file']},{result['line']},\"{result['text']}\",\"{result['context']}\","
                  f"{result['analysis']['has_translation']},\"{similar_keys}\"")
    
    def run(self):
        """Main execution"""
        print("🔍 Arabic Text Scanner")
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
        all_results = []
        for directory in directories:
            if directory.exists():
                results = self.scan_directory(directory)
                all_results.extend(results)
        
        print(f"\n✅ Found {len(all_results)} Arabic text occurrences")
        
        # Analyze results
        analyzed = self.analyze_results(all_results)
        
        # Generate report
        self.generate_report(analyzed)
        
        # Auto-sync if enabled
        if self.auto_sync:
            untranslated = [r for r in analyzed if not r['analysis']['has_translation']]
            if untranslated:
                self.add_missing_translations(untranslated)
            else:
                print("ℹ️  No untranslated texts found for auto-sync")
        
        return 0


def main():
    parser = argparse.ArgumentParser(
        description='Scan source code for Arabic text and check translation coverage'
    )
    parser.add_argument(
        '--dir',
        help='Specific directory to scan (default: all frontend src directories)'
    )
    parser.add_argument(
        '--exclude',
        nargs='+',
        default=['node_modules', 'dist', 'build', '.git', 'coverage'],
        help='Patterns to exclude from scanning'
    )
    parser.add_argument(
        '--output',
        choices=['console', 'json', 'csv'],
        default='console',
        help='Output format (default: console)'
    )
    parser.add_argument(
        '--auto-sync',
        action='store_true',
        help='Automatically add missing translation keys and sync to all languages'
    )
    parser.add_argument(
        '--key-prefix',
        default='auto_',
        help='Prefix for auto-generated translation keys (default: "auto_")'
    )
    
    args = parser.parse_args()
    
    scanner = ArabicTextScanner(
        target_dir=args.dir,
        exclude_patterns=args.exclude,
        output_format=args.output,
        auto_sync=args.auto_sync,
        key_prefix=args.key_prefix
    )
    
    sys.exit(scanner.run())


if __name__ == '__main__':
    main()

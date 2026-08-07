#!/usr/bin/env python3
"""
Code Extractor Tool
Extracts and analyzes code from WorkTrack PWA project
"""

import os
import re
from pathlib import Path
from typing import Dict, List, Optional
import json

class CodeExtractor:
    def __init__(self, project_path: str):
        self.project_path = Path(project_path)
        self.vue_files = []
        self.js_files = []
        self.css_files = []
        
    def scan_project(self):
        """Scan project for Vue, JS, and CSS files"""
        self.vue_files = list(self.project_path.rglob("*.vue"))
        self.js_files = list(self.project_path.rglob("*.js"))
        self.css_files = list(self.project_path.rglob("*.css"))
        
        print(f"Found {len(self.vue_files)} Vue files")
        print(f"Found {len(self.js_files)} JS files")
        print(f"Found {len(self.css_files)} CSS files")
        
    def extract_vue_code(self, file_path: Path) -> Dict:
        """Extract template, script, and style from Vue file"""
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Extract template
        template_match = re.search(r'<template>(.*?)</template>', content, re.DOTALL)
        template = template_match.group(1).strip() if template_match else ""
        
        # Extract script
        script_match = re.search(r'<script[^>]*>(.*?)</script>', content, re.DOTALL)
        script = script_match.group(1).strip() if script_match else ""
        
        # Extract style
        style_match = re.search(r'<style[^>]*>(.*?)</style>', content, re.DOTALL)
        style = style_match.group(1).strip() if style_match else ""
        
        # Get relative path if possible
        try:
            rel_path = str(file_path.relative_to(self.project_path))
        except ValueError:
            rel_path = str(file_path)
        
        return {
            'file': rel_path,
            'template': template,
            'script': script,
            'style': style,
            'template_lines': len(template.split('\n')),
            'script_lines': len(script.split('\n')),
            'style_lines': len(style.split('\n'))
        }
    
    def extract_js_code(self, file_path: Path) -> Dict:
        """Extract code from JS file"""
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        try:
            rel_path = str(file_path.relative_to(self.project_path))
        except ValueError:
            rel_path = str(file_path)
        
        return {
            'file': rel_path,
            'code': content,
            'lines': len(content.split('\n'))
        }
    
    def extract_css_code(self, file_path: Path) -> Dict:
        """Extract code from CSS file"""
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        try:
            rel_path = str(file_path.relative_to(self.project_path))
        except ValueError:
            rel_path = str(file_path)
        
        return {
            'file': rel_path,
            'code': content,
            'lines': len(content.split('\n'))
        }
    
    def extract_all_vue(self) -> List[Dict]:
        """Extract all Vue files"""
        results = []
        for vue_file in self.vue_files:
            try:
                result = self.extract_vue_code(vue_file)
                results.append(result)
            except Exception as e:
                print(f"Error processing {vue_file}: {e}")
        return results
    
    def extract_all_js(self) -> List[Dict]:
        """Extract all JS files"""
        results = []
        for js_file in self.js_files:
            try:
                result = self.extract_js_code(js_file)
                results.append(result)
            except Exception as e:
                print(f"Error processing {js_file}: {e}")
        return results
    
    def extract_all_css(self) -> List[Dict]:
        """Extract all CSS files"""
        results = []
        for css_file in self.css_files:
            try:
                result = self.extract_css_code(css_file)
                results.append(result)
            except Exception as e:
                print(f"Error processing {css_file}: {e}")
        return results
    
    def save_to_json(self, output_file: str):
        """Save all extracted code to JSON file"""
        data = {
            'vue_files': self.extract_all_vue(),
            'js_files': self.extract_all_js(),
            'css_files': self.extract_all_css()
        }
        
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        
        print(f"Saved to {output_file}")
    
    def print_summary(self):
        """Print summary of extracted code"""
        print("\n" + "="*60)
        print("CODE EXTRACTION SUMMARY")
        print("="*60)
        
        print("\nVue Files:")
        for vue_file in self.vue_files:
            print(f"  - {vue_file.relative_to(self.project_path)}")
        
        print("\nJS Files:")
        for js_file in self.js_files:
            print(f"  - {js_file.relative_to(self.project_path)}")
        
        print("\nCSS Files:")
        for css_file in self.css_files:
            print(f"  - {css_file.relative_to(self.project_path)}")
    
    def extract_specific_file(self, file_path: str) -> Dict:
        """Extract code from a specific file"""
        file_path = Path(file_path)
        if not file_path.exists():
            file_path = self.project_path / file_path
        
        if not file_path.exists():
            return {'error': f'File not found: {file_path}'}
        
        if file_path.suffix == '.vue':
            return self.extract_vue_code(file_path)
        elif file_path.suffix == '.js':
            return self.extract_js_code(file_path)
        elif file_path.suffix == '.css':
            return self.extract_css_code(file_path)
        else:
            return {'error': 'Unsupported file type'}


def main():
    project_path = "/home/dev-bit/worktrack/frontend-worker-pwa"
    
    extractor = CodeExtractor(project_path)
    extractor.scan_project()
    extractor.print_summary()
    
    # Save to JSON
    output_file = "/home/dev-bit/worktrack/frontend-worker-pwa/extracted_code.json"
    extractor.save_to_json(output_file)
    
    # Extract specific files
    print("\n" + "="*60)
    print("SPECIFIC FILE EXTRACTION")
    print("="*60)
    
    files_to_extract = [
        "src/views/AttendanceView.vue",
        "src/views/ProfileView.vue",
        "src/views/NotesView.vue",
        "src/App.vue",
        "src/styles/tokens.css"
    ]
    
    for file_path in files_to_extract:
        result = extractor.extract_specific_file(file_path)
        if 'error' not in result:
            print(f"\n{file_path}:")
            if 'template_lines' in result:
                print(f"  Template: {result.get('template_lines', 0)} lines")
                print(f"  Script: {result.get('script_lines', 0)} lines")
                print(f"  Style: {result.get('style_lines', 0)} lines")
            else:
                print(f"  Lines: {result.get('lines', 0)}")
        else:
            print(f"\n{file_path}: {result['error']}")


if __name__ == "__main__":
    main()

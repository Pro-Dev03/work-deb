#!/usr/bin/env python3
"""
HTML to PDF Converter
Converts HTML files to PDF using available system tools
"""

import os
import subprocess
import time
from pathlib import Path
from datetime import datetime

class HTMLToPDFConverter:
    def __init__(self):
        self.project_root = Path(".")
        self.html_files = {
            "hebrew": {
                "html": "מדריך_המשתמש_המלא.html",
                "pdf": "מדריך_המשתמש_המלא.pdf",
                "name": "Hebrew User Guide"
            },
            "english": {
                "html": "Complete_User_Guide.html",
                "pdf": "Complete_User_Guide.pdf", 
                "name": "English User Guide"
            },
            "arabic": {
                "html": "دليل_المستخدم_الشامل.html",
                "pdf": "دليل_المستخدم_الشامل.pdf",
                "name": "Arabic User Guide"
            }
        }
    
    def check_firefox(self):
        """Check if Firefox is available"""
        try:
            result = subprocess.run(['which', 'firefox'], capture_output=True, text=True)
            return result.returncode == 0
        except:
            return False
    
    def convert_with_firefox(self, html_file, pdf_file):
        """Convert HTML to PDF using Firefox headless mode"""
        try:
            # Get absolute paths
            html_path = os.path.abspath(html_file)
            pdf_path = os.path.abspath(pdf_file)
            
            # Use Firefox to print to PDF
            cmd = [
                'firefox',
                '--headless',
                '--print-to-pdf=' + pdf_path,
                html_path
            ]
            
            print(f"Running: {' '.join(cmd)}")
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
            
            if result.returncode == 0 and os.path.exists(pdf_file):
                print(f"✅ Successfully created: {pdf_file}")
                return True
            else:
                print(f"❌ Firefox conversion failed for {html_file}")
                print(f"Error: {result.stderr}")
                return False
                
        except subprocess.TimeoutExpired:
            print(f"⏱️ Firefox conversion timed out for {html_file}")
            return False
        except Exception as e:
            print(f"❌ Error converting {html_file}: {e}")
            return False
    
    def convert_manual_instructions(self, language, html_file, pdf_file):
        """Provide manual conversion instructions"""
        print(f"\n{'='*60}")
        print(f"MANUAL CONVERSION REQUIRED FOR: {language.upper()}")
        print(f"{'='*60}")
        print(f"HTML File: {html_file}")
        print(f"PDF File: {pdf_file}")
        print(f"\nSteps to convert manually:")
        print(f"1. Open {html_file} in your web browser")
        print(f"2. Press Ctrl+P (or Cmd+P on Mac)")
        print(f"3. Select 'Save as PDF' as the printer")
        print(f"4. Click 'Save' and choose the name: {pdf_file}")
        print(f"{'='*60}\n")
    
    def convert_all(self):
        """Convert all HTML files to PDF"""
        print("="*60)
        print("HTML TO PDF CONVERSION")
        print("="*60)
        print(f"Started: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print()
        
        # Check for Firefox
        has_firefox = self.check_firefox()
        print(f"Firefox available: {has_firefox}")
        print()
        
        results = {}
        
        for lang, files in self.html_files.items():
            html_file = files['html']
            pdf_file = files['pdf']
            name = files['name']
            
            print(f"Processing {name} ({lang})...")
            print(f"  HTML: {html_file}")
            print(f"  PDF:  {pdf_file}")
            
            # Check if HTML file exists
            if not os.path.exists(html_file):
                print(f"  ❌ HTML file not found: {html_file}")
                results[lang] = False
                continue
            
            # Check if PDF already exists
            if os.path.exists(pdf_file):
                print(f"  ⚠️ PDF already exists: {pdf_file}")
                results[lang] = True
                continue
            
            # Try automatic conversion if Firefox is available
            if has_firefox:
                success = self.convert_with_firefox(html_file, pdf_file)
                results[lang] = success
            else:
                print(f"  ℹ️ No automatic conversion tool available")
                self.convert_manual_instructions(lang, html_file, pdf_file)
                results[lang] = False
            
            print()
        
        # Summary
        print("="*60)
        print("CONVERSION SUMMARY")
        print("="*60)
        
        successful = sum(1 for result in results.values() if result)
        total = len(results)
        
        print(f"Successfully converted: {successful}/{total}")
        print()
        
        for lang, success in results.items():
            status = "✅" if success else "❌"
            name = self.html_files[lang]['name']
            pdf_file = self.html_files[lang]['pdf']
            print(f"{status} {name}: {pdf_file}")
        
        print()
        print(f"Completed: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print("="*60)
        
        return results

if __name__ == "__main__":
    converter = HTMLToPDFConverter()
    converter.convert_all()
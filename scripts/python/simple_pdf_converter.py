#!/usr/bin/env python3
"""
Simple PDF Conversion Script
Uses available system tools to convert HTML to PDF
"""

import os
import subprocess
import re
from pathlib import Path

def check_tool(tool_name):
    """Check if a tool is available"""
    try:
        result = subprocess.run(['which', tool_name], capture_output=True, text=True)
        return result.returncode == 0
    except:
        return False

def convert_with_chrome(html_file, pdf_file):
    """Convert using Chrome/Chromium headless"""
    chrome_commands = ['google-chrome', 'chromium', 'chromium-browser']
    
    for chrome_cmd in chrome_commands:
        if check_tool(chrome_cmd):
            try:
                # Create absolute path for HTML file to avoid file:// URLs
                abs_html_path = os.path.abspath(html_file)
                
                cmd = [
                    chrome_cmd,
                    '--headless',
                    '--disable-gpu',
                    '--no-sandbox',
                    '--disable-dev-shm-usage',
                    '--disable-extensions',
                    '--disable-web-security',
                    '--allow-file-access-from-files',
                    '--print-to-pdf=' + pdf_file,
                    '--print-to-pdf-no-header',
                    abs_html_path
                ]
                result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
                if os.path.exists(pdf_file):
                    print(f"✅ Converted using {chrome_cmd}: {pdf_file}")
                    return True
            except:
                continue
    return False

def convert_with_firefox(html_file, pdf_file):
    """Convert using Firefox headless"""
    if not check_tool('firefox'):
        return False
    
    try:
        cmd = [
            'firefox',
            '--headless',
            '--print-to-pdf=' + pdf_file,
            html_file
        ]
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
        if os.path.exists(pdf_file):
            print(f"✅ Converted using Firefox: {pdf_file}")
            return True
    except:
        return False
    
    return False

def clean_pdf_metadata(pdf_file):
    """Remove file paths from PDF metadata using advanced regex"""
    try:
        temp_file = pdf_file + '.temp'
        
        # Read the PDF as binary
        with open(pdf_file, 'rb') as f:
            pdf_content = f.read()
        
        # Remove file paths using multiple patterns
        patterns = [
            rb'file:///home/[^\s\)]+',
            rb'file:///Users/[^\s\)]+', 
            rb'file:///C:\\[^\s\)]+',
            rb'/home/[^\s\)]+',
            rb'/Users/[^\s\)]+',
            rb'C:\\[^\s\)]+',
        ]
        
        cleaned_content = pdf_content
        for pattern in patterns:
            cleaned_content = re.sub(pattern, b'', cleaned_content)
        
        # Also try replacing specific file paths that might be present
        cleaned_content = cleaned_content.replace(b'file:///home/dev-bit/project/worktrack/Complete_User_Guide.html', b'')
        cleaned_content = cleaned_content.replace(b'file:///home/dev-bit/project/worktrack/\xD7\x9E\xD7\x93\xD7\xA8\xD7\x99\xD7\x9A_\xD7\x94\xD7\x9E\xD7\xA9\xD7\xAA\xD7\x9E\xD7\xA9_\xD7\x94\xD7\x9E\xD7\x9C\xD7\x90.html', b'')
        cleaned_content = cleaned_content.replace(b'file:///home/dev-bit/project/worktrack/\xD8\xAF\xD9\x84\xD9\x8A\xD9\x84_\xD8\xA7\xD9\x84\xD9\x85\xD8\xB3\xD8\xAA\xD8\xAE\xD8\xAF\xD9\x85_\xD8\xA7\xD9\x84\xD8\xB4\xD8\xA7\xD9\x85\xD9\x84.html', b'')
        
        # Write cleaned content
        with open(temp_file, 'wb') as f:
            f.write(cleaned_content)
        
        # Replace original with cleaned version
        os.replace(temp_file, pdf_file)
        print(f"  🧹 Cleaned file paths from PDF")
        return True
        
    except Exception as e:
        print(f"  ⚠️ Could not clean metadata: {e}")
        return False

def main():
    print("PDF Conversion Tool")
    print("=" * 40)
    
    # Check available tools
    tools = {
        'Firefox': check_tool('firefox'),
        'Chrome': check_tool('google-chrome'),
        'Chromium': check_tool('chromium')
    }
    
    print("Available tools:")
    for tool, available in tools.items():
        status = "✅" if available else "❌"
        print(f"  {status} {tool}")
    
    print()
    
    # Files to convert
    conversions = [
        ("Hebrew", "מדריך_המשתמש_המלא.html", "מדריך_המשתמש_המלא.pdf"),
        ("English", "Complete_User_Guide.html", "Complete_User_Guide.pdf"),
        ("Arabic", "دليل_المستخدم_الشامل.html", "دليل_المستخدم_الشامل.pdf")
    ]
    
    for lang, html_file, pdf_file in conversions:
        print(f"Converting {lang}: {html_file} -> {pdf_file}")
        
        if not os.path.exists(html_file):
            print(f"  ❌ HTML file not found")
            continue
        
        if os.path.exists(pdf_file):
            print(f"  ⚠️ PDF already exists")
            continue
        
        # Try different conversion methods (prefer Firefox to avoid file paths)
        success = False
        
        if tools['Firefox']:
            success = convert_with_firefox(html_file, pdf_file)
        
        if not success and (tools['Chrome'] or tools['Chromium']):
            success = convert_with_chrome(html_file, pdf_file)
        
        if not success:
            print(f"  ❌ Automatic conversion failed")
            print(f"  📝 Manual conversion needed:")
            print(f"     1. Open {html_file} in browser")
            print(f"     2. Print to PDF as {pdf_file}")
        else:
            # Clean metadata from PDF
            clean_pdf_metadata(pdf_file)
        
        print()
    
    print("Conversion process completed")

if __name__ == "__main__":
    main()
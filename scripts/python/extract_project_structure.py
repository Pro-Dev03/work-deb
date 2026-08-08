#!/usr/bin/env python3
"""
Project Structure Extractor
Extracts backend and frontend code from WorkTrack project and organizes them into separate .txt files
"""

import os
import glob
from pathlib import Path
from datetime import datetime

class ProjectExtractor:
    def __init__(self, project_root="."):
        self.project_root = Path(project_root)
        self.output_dir = self.project_root / "extracted_code"
        self.output_dir.mkdir(exist_ok=True)
        
        # Define directories to exclude
        self.exclude_dirs = {
            'node_modules', '.git', 'dist', 'build', '__pycache__', 
            '.next', '.nuxt', 'coverage', '.pytest_cache', 'venv', 'env'
        }
        
        # Define file extensions for each category
        self.backend_extensions = {'.go'}
        self.frontend_extensions = {'.vue', '.js', '.jsx', '.ts', '.tsx', '.css', '.html', '.scss'}
        self.config_extensions = {'.json', '.yaml', '.yml', '.toml', '.ini', '.env', '.md', '.sql'}
        
    def should_exclude(self, path):
        """Check if path should be excluded"""
        parts = path.parts
        return any(part in self.exclude_dirs for part in parts)
    
    def get_file_content(self, file_path):
        """Read file content with proper encoding"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                return f.read()
        except UnicodeDecodeError:
            try:
                with open(file_path, 'r', encoding='latin-1') as f:
                    return f.read()
            except Exception as e:
                return f"[Error reading file: {e}]"
        except Exception as e:
            return f"[Error reading file: {e}]"
    
    def extract_backend_files(self):
        """Extract all backend Go files"""
        print("Extracting backend files...")
        backend_dir = self.project_root / "backend"
        if not backend_dir.exists():
            print("Backend directory not found")
            return
        
        output_file = self.output_dir / "01_BACKEND_GO_FILES.txt"
        
        with open(output_file, 'w', encoding='utf-8') as out:
            out.write("=" * 80 + "\n")
            out.write("BACKEND - GO FILES\n")
            out.write(f"Extracted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            out.write("=" * 80 + "\n\n")
            
            go_files = list(backend_dir.rglob("*.go"))
            total_files = len(go_files)
            out.write(f"Total files found: {total_files}\n\n")
            out.write("=" * 80 + "\n\n")
            
            for i, file_path in enumerate(go_files, 1):
                if self.should_exclude(file_path):
                    continue
                    
                relative_path = file_path.relative_to(self.project_root)
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE {i}/{total_files}: {relative_path}\n")
                out.write(f"{'=' * 80}\n\n")
                
                content = self.get_file_content(file_path)
                out.write(content)
                out.write("\n\n")
        
        print(f"Backend files extracted to: {output_file}")
    
    def extract_admin_dashboard(self):
        """Extract admin dashboard frontend files"""
        print("Extracting admin dashboard files...")
        admin_dir = self.project_root / "frontend-admin-dashboard"
        if not admin_dir.exists():
            print("Admin dashboard directory not found")
            return
        
        output_file = self.output_dir / "02_FRONTEND_ADMIN_DASHBOARD.txt"
        
        with open(output_file, 'w', encoding='utf-8') as out:
            out.write("=" * 80 + "\n")
            out.write("FRONTEND - ADMIN DASHBOARD\n")
            out.write(f"Extracted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            out.write("=" * 80 + "\n\n")
            
            # Get all relevant frontend files
            frontend_files = []
            for ext in self.frontend_extensions:
                frontend_files.extend(admin_dir.rglob(f"*{ext}"))
            
            # Filter out excluded directories
            frontend_files = [f for f in frontend_files if not self.should_exclude(f)]
            total_files = len(frontend_files)
            
            out.write(f"Total files found: {total_files}\n\n")
            out.write("=" * 80 + "\n\n")
            
            for i, file_path in enumerate(frontend_files, 1):
                relative_path = file_path.relative_to(self.project_root)
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE {i}/{total_files}: {relative_path}\n")
                out.write(f"{'=' * 80}\n\n")
                
                content = self.get_file_content(file_path)
                out.write(content)
                out.write("\n\n")
        
        print(f"Admin dashboard files extracted to: {output_file}")
    
    def extract_client_portal(self):
        """Extract client portal frontend files"""
        print("Extracting client portal files...")
        client_dir = self.project_root / "frontend-client-portal"
        if not client_dir.exists():
            print("Client portal directory not found")
            return
        
        output_file = self.output_dir / "03_FRONTEND_CLIENT_PORTAL.txt"
        
        with open(output_file, 'w', encoding='utf-8') as out:
            out.write("=" * 80 + "\n")
            out.write("FRONTEND - CLIENT PORTAL\n")
            out.write(f"Extracted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            out.write("=" * 80 + "\n\n")
            
            # Get all relevant frontend files
            frontend_files = []
            for ext in self.frontend_extensions:
                frontend_files.extend(client_dir.rglob(f"*{ext}"))
            
            # Filter out excluded directories
            frontend_files = [f for f in frontend_files if not self.should_exclude(f)]
            total_files = len(frontend_files)
            
            out.write(f"Total files found: {total_files}\n\n")
            out.write("=" * 80 + "\n\n")
            
            for i, file_path in enumerate(frontend_files, 1):
                relative_path = file_path.relative_to(self.project_root)
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE {i}/{total_files}: {relative_path}\n")
                out.write(f"{'=' * 80}\n\n")
                
                content = self.get_file_content(file_path)
                out.write(content)
                out.write("\n\n")
        
        print(f"Client portal files extracted to: {output_file}")
    
    def extract_worker_pwa(self):
        """Extract worker PWA frontend files"""
        print("Extracting worker PWA files...")
        worker_dir = self.project_root / "frontend-worker-pwa"
        if not worker_dir.exists():
            print("Worker PWA directory not found")
            return
        
        output_file = self.output_dir / "04_FRONTEND_WORKER_PWA.txt"
        
        with open(output_file, 'w', encoding='utf-8') as out:
            out.write("=" * 80 + "\n")
            out.write("FRONTEND - WORKER PWA\n")
            out.write(f"Extracted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            out.write("=" * 80 + "\n\n")
            
            # Get all relevant frontend files
            frontend_files = []
            for ext in self.frontend_extensions:
                frontend_files.extend(worker_dir.rglob(f"*{ext}"))
            
            # Filter out excluded directories
            frontend_files = [f for f in frontend_files if not self.should_exclude(f)]
            total_files = len(frontend_files)
            
            out.write(f"Total files found: {total_files}\n\n")
            out.write("=" * 80 + "\n\n")
            
            for i, file_path in enumerate(frontend_files, 1):
                relative_path = file_path.relative_to(self.project_root)
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE {i}/{total_files}: {relative_path}\n")
                out.write(f"{'=' * 80}\n\n")
                
                content = self.get_file_content(file_path)
                out.write(content)
                out.write("\n\n")
        
        print(f"Worker PWA files extracted to: {output_file}")
    
    def extract_config_files(self):
        """Extract configuration and setup files"""
        print("Extracting configuration files...")
        output_file = self.output_dir / "05_CONFIG_AND_SETUP.txt"
        
        with open(output_file, 'w', encoding='utf-8') as out:
            out.write("=" * 80 + "\n")
            out.write("CONFIGURATION AND SETUP FILES\n")
            out.write(f"Extracted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            out.write("=" * 80 + "\n\n")
            
            config_files = []
            
            # Get config files from root directory
            for ext in self.config_extensions:
                config_files.extend(self.project_root.glob(f"*{ext}"))
            
            # Get config files from subdirectories (but not node_modules)
            for ext in self.config_extensions:
                for file_path in self.project_root.rglob(f"*{ext}"):
                    if not self.should_exclude(file_path):
                        config_files.append(file_path)
            
            # Remove duplicates and sort
            config_files = list(set(config_files))
            config_files.sort()
            total_files = len(config_files)
            
            out.write(f"Total files found: {total_files}\n\n")
            out.write("=" * 80 + "\n\n")
            
            for i, file_path in enumerate(config_files, 1):
                relative_path = file_path.relative_to(self.project_root)
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE {i}/{total_files}: {relative_path}\n")
                out.write(f"{'=' * 80}\n\n")
                
                content = self.get_file_content(file_path)
                out.write(content)
                out.write("\n\n")
        
        print(f"Configuration files extracted to: {output_file}")
    
    def extract_database_files(self):
        """Extract database schema and migration files"""
        print("Extracting database files...")
        output_file = self.output_dir / "06_DATABASE_FILES.txt"
        
        with open(output_file, 'w', encoding='utf-8') as out:
            out.write("=" * 80 + "\n")
            out.write("DATABASE SCHEMA AND MIGRATION FILES\n")
            out.write(f"Extracted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            out.write("=" * 80 + "\n\n")
            
            # Get SQL files
            sql_files = list(self.project_root.rglob("*.sql"))
            sql_files = [f for f in sql_files if not self.should_exclude(f)]
            sql_files.sort()
            total_files = len(sql_files)
            
            out.write(f"Total SQL files found: {total_files}\n\n")
            out.write("=" * 80 + "\n\n")
            
            for i, file_path in enumerate(sql_files, 1):
                relative_path = file_path.relative_to(self.project_root)
                out.write(f"\n{'=' * 80}\n")
                out.write(f"FILE {i}/{total_files}: {relative_path}\n")
                out.write(f"{'=' * 80}\n\n")
                
                content = self.get_file_content(file_path)
                out.write(content)
                out.write("\n\n")
        
        print(f"Database files extracted to: {output_file}")
    
    def generate_structure_summary(self):
        """Generate a summary of the project structure"""
        print("Generating structure summary...")
        output_file = self.output_dir / "00_PROJECT_STRUCTURE_SUMMARY.txt"
        
        with open(output_file, 'w', encoding='utf-8') as out:
            out.write("=" * 80 + "\n")
            out.write("WORKTRACK PROJECT STRUCTURE SUMMARY\n")
            out.write(f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            out.write("=" * 80 + "\n\n")
            
            # Backend structure
            out.write("BACKEND (Go)\n")
            out.write("-" * 40 + "\n")
            backend_dir = self.project_root / "backend"
            if backend_dir.exists():
                self._print_directory_structure(backend_dir, out, exclude=self.exclude_dirs)
            
            # Frontend Admin Dashboard
            out.write("\n\nFRONTEND - ADMIN DASHBOARD (Vue.js)\n")
            out.write("-" * 40 + "\n")
            admin_dir = self.project_root / "frontend-admin-dashboard"
            if admin_dir.exists():
                self._print_directory_structure(admin_dir, out, exclude=self.exclude_dirs)
            
            # Frontend Client Portal
            out.write("\n\nFRONTEND - CLIENT PORTAL (Vue.js)\n")
            out.write("-" * 40 + "\n")
            client_dir = self.project_root / "frontend-client-portal"
            if client_dir.exists():
                self._print_directory_structure(client_dir, out, exclude=self.exclude_dirs)
            
            # Frontend Worker PWA
            out.write("\n\nFRONTEND - WORKER PWA (Vue.js)\n")
            out.write("-" * 40 + "\n")
            worker_dir = self.project_root / "frontend-worker-pwa"
            if worker_dir.exists():
                self._print_directory_structure(worker_dir, out, exclude=self.exclude_dirs)
            
            # Root files
            out.write("\n\nROOT CONFIGURATION FILES\n")
            out.write("-" * 40 + "\n")
            for item in sorted(self.project_root.iterdir()):
                if item.is_file() and not item.name.startswith('.'):
                    out.write(f"  {item.name}\n")
        
        print(f"Structure summary generated: {output_file}")
    
    def _print_directory_structure(self, directory, out, exclude=None, prefix="", max_depth=3):
        """Helper function to print directory structure"""
        if exclude is None:
            exclude = set()
        
        if max_depth <= 0:
            return
        
        try:
            items = sorted(directory.iterdir())
        except PermissionError:
            return
        
        for i, item in enumerate(items):
            if item.name in exclude:
                continue
            
            is_last = i == len(items) - 1
            current_prefix = "└── " if is_last else "├── "
            out.write(f"{prefix}{current_prefix}{item.name}\n")
            
            if item.is_dir() and max_depth > 1:
                next_prefix = prefix + ("    " if is_last else "│   ")
                self._print_directory_structure(item, out, exclude, next_prefix, max_depth - 1)
    
    def run_extraction(self):
        """Run the complete extraction process"""
        print("=" * 80)
        print("WORKTRACK PROJECT EXTRACTION")
        print("=" * 80)
        print(f"Project root: {self.project_root}")
        print(f"Output directory: {self.output_dir}")
        print("=" * 80)
        print()
        
        # Generate structure summary
        self.generate_structure_summary()
        
        # Extract all components
        self.extract_backend_files()
        self.extract_admin_dashboard()
        self.extract_client_portal()
        self.extract_worker_pwa()
        self.extract_config_files()
        self.extract_database_files()
        
        print()
        print("=" * 80)
        print("EXTRACTION COMPLETED SUCCESSFULLY")
        print("=" * 80)
        print(f"All files extracted to: {self.output_dir}")
        print()
        print("Generated files:")
        for file_path in sorted(self.output_dir.glob("*.txt")):
            file_size = file_path.stat().st_size
            print(f"  {file_path.name} ({file_size:,} bytes)")

if __name__ == "__main__":
    # Run extraction from current directory
    extractor = ProjectExtractor(".")
    extractor.run_extraction()
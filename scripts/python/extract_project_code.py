#!/usr/bin/env python3
"""
WorkTrack Project Code Extractor
يقوم هذا السكربت باستخراج الأكواد المهمة من المشروع وتقسيمها إلى 5 ملفات
مع دعم توليد ملفات iOS و Android

📝 وظيفة هذا السكريبت:
- استخراج الأكواد المهمة من المشروع بشكل منظم
- تقسيم الأكواد إلى 5 تصنيفات رئيسية:
  1. Frontend Code (Vue/JS)
  2. Backend Code
  3. Configuration Files
  4. Build & Deployment Scripts
  5. Security & Obfuscation
- دعم توليد ملفات iOS و Android إضافية
- استبعاد المجلدات غير المرغوبة (node_modules, dist, build, .git)
- تنظيم الملفات بتنسيق Markdown مع تظليل اللغة

🚀 الاستخدام:
  python3 scripts/python/extract_project_code.py --path /path/to/project
  python3 scripts/python/extract_project_code.py --path /path/to/project --mobile
"""

import os
import shutil
from pathlib import Path
from datetime import datetime

class ProjectCodeExtractor:
    def __init__(self, project_root):
        self.project_root = Path(project_root)
        self.output_dir = self.project_root / "extracted_code"
        self.timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        
        # المجلدات المستهدفة
        self.frontend_dirs = [
            "frontend-worker-pwa",
            "frontend-admin-dashboard", 
            "frontend-client-portal"
        ]
        self.backend_dir = "backend"
        
        # أنواع الملفات لكل تصنيف
        self.file_categories = {
            'frontend': {
                'extensions': ['.vue', '.js', '.ts', '.jsx', '.tsx'],
                'dirs': self.frontend_dirs,
                'exclude_dirs': ['node_modules', 'dist', 'build', '.git', 'android', 'ios', 'electron']
            },
            'backend': {
                'extensions': ['.py', '.go', '.js', '.ts', '.sql'],
                'dirs': [self.backend_dir, 'supabase_migrations'],
                'exclude_dirs': ['node_modules', '__pycache__', '.git', 'venv']
            },
            'config': {
                'extensions': ['.json', '.yaml', '.yml', '.toml', '.ini', '.env', '.xml', '.gradle', '.properties'],
                'dirs': ['.', *self.frontend_dirs, self.backend_dir],
                'exclude_dirs': ['node_modules', 'dist', 'build', '.git', 'android', 'ios']
            },
            'build': {
                'extensions': ['.sh', '.bat', '.js', '.py'],
                'names': ['build', 'deploy', 'setup', 'install', 'script'],
                'dirs': ['.'],
                'exclude_dirs': ['node_modules', 'dist', 'build', '.git']
            },
            'security': {
                'extensions': ['.js', '.md', '.pro'],
                'names': ['obfuscator', 'security', 'keystore', 'proguard', 'privacy'],
                'dirs': ['.'],
                'exclude_dirs': ['node_modules', 'dist', 'build', '.git']
            }
        }
        
    def create_output_directory(self):
        """إنشاء مجلد الإخراج"""
        self.output_dir.mkdir(exist_ok=True)
        print(f"✅ تم إنشاء مجلد الإخراج: {self.output_dir}")
        
    def should_include_file(self, file_path, category):
        """تحديد ما إذا كان الملف يجب تضمينه في التصنيف"""
        category_config = self.file_categories[category]
        
        # التحقق من الامتداد
        if 'extensions' in category_config:
            if file_path.suffix.lower() not in category_config['extensions']:
                return False
        
        # التحقق من الاسم (للملفات الخاصة)
        if 'names' in category_config:
            name_match = any(name in file_path.name.lower() for name in category_config['names'])
            if not name_match:
                return False
        
        # التحقق من المجلدات المستثناة
        for exclude_dir in category_config['exclude_dirs']:
            if exclude_dir in file_path.parts:
                return False
        
        # التحقق من المجلدات المسموحة
        if 'dirs' in category_config:
            in_allowed_dir = False
            for allowed_dir in category_config['dirs']:
                if allowed_dir == '.' or str(allowed_dir) in file_path.parts:
                    in_allowed_dir = True
                    break
            if not in_allowed_dir:
                return False
        
        return True
    
    def extract_file_content(self, file_path):
        """استخراج محتوى الملف"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            return content
        except Exception as e:
            return f"# Error reading file: {e}\n"
    
    def get_files_for_category(self, category):
        """الحصول على جميع الملفات لتصنيف معين"""
        files = []
        
        if category == 'frontend':
            # البحث في مجلدات Frontend
            for frontend_dir in self.frontend_dirs:
                dir_path = self.project_root / frontend_dir
                if dir_path.exists():
                    for file_path in dir_path.rglob('*'):
                        if file_path.is_file() and self.should_include_file(file_path, category):
                            files.append(file_path)
                            
        elif category == 'backend':
            # البحث في مجلد Backend
            for backend_dir in [self.backend_dir, 'supabase_migrations']:
                dir_path = self.project_root / backend_dir
                if dir_path.exists():
                    for file_path in dir_path.rglob('*'):
                        if file_path.is_file() and self.should_include_file(file_path, category):
                            files.append(file_path)
                            
        elif category in ['config', 'build', 'security']:
            # البحث في المشروع بالكامل
            for file_path in self.project_root.rglob('*'):
                if file_path.is_file() and self.should_include_file(file_path, category):
                    files.append(file_path)
        
        return sorted(files)
    
    def create_category_file(self, category, platform=None):
        """إنشاء ملف لتصنيف معين"""
        category_names = {
            'frontend': 'Frontend Code (Vue/JS)',
            'backend': 'Backend Code',
            'config': 'Configuration Files',
            'build': 'Build & Deployment Scripts',
            'security': 'Security & Obfuscation'
        }
        
        platform_suffix = f" - {platform}" if platform else ""
        filename = f"{category_names[category].lower().replace(' ', '_').replace('/', '_')}{platform_suffix}.md"
        output_path = self.output_dir / filename
        
        files = self.get_files_for_category(category)
        
        if platform:
            # تصفية الملفات حسب المنصة
            files = [f for f in files if platform.lower() in str(f).lower()]
        
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(f"# {category_names[category]}{platform_suffix}\n\n")
            f.write(f"## تم الاستخراج في: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n")
            f.write(f"## عدد الملفات: {len(files)}\n\n")
            f.write("---\n\n")
            
            for file_path in files:
                relative_path = file_path.relative_to(self.project_root)
                f.write(f"## 📄 {relative_path}\n\n")
                f.write(f"```{self.get_language(file_path)}\n")
                f.write(self.extract_file_content(file_path))
                f.write("\n```\n\n")
                f.write("---\n\n")
        
        print(f"✅ تم إنشاء ملف: {filename} ({len(files)} ملف)")
        return output_path
    
    def get_language(self, file_path):
        """تحديد لغة البرمجة للتظليل"""
        ext = file_path.suffix.lower()
        lang_map = {
            '.vue': 'vue',
            '.js': 'javascript',
            '.ts': 'typescript',
            '.jsx': 'javascript',
            '.tsx': 'typescript',
            '.py': 'python',
            '.go': 'go',
            '.sql': 'sql',
            '.json': 'json',
            '.yaml': 'yaml',
            '.yml': 'yaml',
            '.xml': 'xml',
            '.gradle': 'groovy',
            '.properties': 'properties',
            '.sh': 'bash',
            '.bat': 'batch',
            '.md': 'markdown',
            '.pro': 'prolog'
        }
        return lang_map.get(ext, 'text')
    
    def generate_mobile_files(self, platform):
        """توليد ملفات خاصة لمنصة معينة (iOS/Android)"""
        print(f"\n📱 توليد ملفات {platform}...")
        
        platform_files = {
            'ios': ['ios', '.podspec', 'Podfile', '.plist', '.swift', '.m', '.h'],
            'android': ['android', '.gradle', 'AndroidManifest.xml', '.java', '.kt']
        }
        
        # المجلدات المستثناة للمنصات المحمولة
        mobile_exclude_dirs = [
            'build', 'gradle', '.gradle', 'captures', 
            '.idea', '.externalNativeBuild', 'cxx',
            'intermediates', 'generated', 'outputs'
        ]
        
        mobile_content = f"# Mobile Platform Code - {platform}\n\n"
        mobile_content += f"## تم الاستخراج في: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n"
        mobile_content += "---\n\n"
        
        files_found = []
        
        for frontend_dir in self.frontend_dirs:
            platform_dir = self.project_root / frontend_dir / platform.lower()
            if platform_dir.exists():
                for file_path in platform_dir.rglob('*'):
                    if file_path.is_file():
                        # استبعاد المجلدات غير المرغوبة
                        if any(exclude in str(file_path) for exclude in mobile_exclude_dirs):
                            continue
                        
                        # التحقق من relevance للمنصة
                        is_relevant = False
                        for keyword in platform_files[platform.lower()]:
                            if keyword.lower() in str(file_path).lower():
                                is_relevant = True
                                break
                        
                        if is_relevant:
                            files_found.append(file_path)
        
        for file_path in files_found:
            relative_path = file_path.relative_to(self.project_root)
            mobile_content += f"## 📄 {relative_path}\n\n"
            mobile_content += f"```{self.get_language(file_path)}\n"
            mobile_content += self.extract_file_content(file_path)
            mobile_content += "\n```\n\n"
            mobile_content += "---\n\n"
        
        filename = f"mobile_{platform.lower()}_code.md"
        output_path = self.output_dir / filename
        
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(mobile_content)
        
        print(f"✅ تم إنشاء ملف: {filename} ({len(files_found)} ملف)")
        return output_path
    
    def extract_all(self, include_mobile=False):
        """استخراج جميع التصنيفات"""
        print("🚀 بدء استخراج أكواد المشروع...")
        print(f"📁 مسار المشروع: {self.project_root}")
        
        self.create_output_directory()
        
        # استخراج التصنيفات الخمسة الرئيسية
        categories = ['frontend', 'backend', 'config', 'build', 'security']
        
        for category in categories:
            print(f"\n📦 استخراج {category}...")
            self.create_category_file(category)
        
        # توليد ملفات Mobile إذا طُلب
        if include_mobile:
            print("\n📱 توليد ملفات المنصات المحمولة...")
            self.generate_mobile_files('ios')
            self.generate_mobile_files('android')
        
        print(f"\n✅ تم الاستخراج بنجاح!")
        print(f"📁 الملفات محفوظة في: {self.output_dir}")
        
        # إنشاء ملف ملخص
        self.create_summary_file(include_mobile)
    
    def create_summary_file(self, include_mobile):
        """إنشاء ملف ملخص"""
        summary_path = self.output_dir / "EXTRACTION_SUMMARY.md"
        
        with open(summary_path, 'w', encoding='utf-8') as f:
            f.write("# Project Code Extraction Summary\n\n")
            f.write(f"**تاريخ الاستخراج:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n")
            f.write(f"**مسار المشروع:** {self.project_root}\n\n")
            f.write("---\n\n")
            f.write("## 📁 الملفات المُستخرجة\n\n")
            f.write("### التصنيفات الرئيسية:\n\n")
            f.write("1. **Frontend Code (Vue/JS)** - أكواد واجهات المستخدم\n")
            f.write("2. **Backend Code** - أكواد الخادم وقاعدة البيانات\n")
            f.write("3. **Configuration Files** - ملفات التكوين والإعدادات\n")
            f.write("4. **Build & Deployment Scripts** - سكريبتات البناء والنشر\n")
            f.write("5. **Security & Obfuscation** - إعدادات الأمان والتعتيم\n\n")
            
            if include_mobile:
                f.write("### ملفات المنصات المحمولة:\n\n")
                f.write("6. **Mobile iOS Code** - أكواد خاصة بـ iOS\n")
                f.write("7. **Mobile Android Code** - أكواد خاصة بـ Android\n\n")
            
            f.write("---\n\n")
            f.write("## 📊 إحصائيات\n\n")
            f.write(f"- **عدد التطبيقات Frontend:** {len(self.frontend_dirs)}\n")
            f.write(f"- **مجلد Backend:** {self.backend_dir}\n")
            f.write(f"- **تاريخ الاستخراج:** {self.timestamp}\n\n")
            f.write("---\n\n")
            f.write("## 🎯 ملاحظات\n\n")
            f.write("- تم استبعاد مجلدات node_modules و dist و build\n")
            f.write("- تم استبعاد مجلد .git\n")
            f.write("- الملفات مرتبة حسب المسار النسبي\n")
            f.write("- كل ملف يحتوي على مساره ومحتواه\n\n")
        
        print(f"✅ تم إنشاء ملف الملخص: EXTRACTION_SUMMARY.md")


def main():
    """الوظيفة الرئيسية"""
    import sys
    import argparse
    
    parser = argparse.ArgumentParser(description='WorkTrack Project Code Extractor')
    parser.add_argument('--path', type=str, default=str(Path(__file__).parent),
                       help='مسار المشروع (افتراضي: مسار السكربت)')
    parser.add_argument('--mobile', action='store_true',
                       help='تضمين ملفات iOS و Android')
    parser.add_argument('--auto', action='store_true',
                       help='تشغيل تلقائي بدون أسئلة')
    
    args = parser.parse_args()
    
    project_root = Path(args.path)
    
    print("🔧 WorkTrack Project Code Extractor")
    print("=" * 50)
    
    # التحقق من وجود المشروع
    if not project_root.exists():
        print(f"❌ خطأ: المسار {project_root} غير موجود")
        sys.exit(1)
    
    # إنشاء المستخرج وتشغيله
    extractor = ProjectCodeExtractor(project_root)
    extractor.extract_all(include_mobile=args.mobile)
    
    print("\n🎉 تم الاستخراج بنجاح!")
    print(f"📁 يمكنك العثور على الملفات في: {extractor.output_dir}")


if __name__ == "__main__":
    main()

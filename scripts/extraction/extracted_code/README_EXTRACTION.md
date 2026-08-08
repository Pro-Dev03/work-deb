# WorkTrack Project Extraction - README

## Overview
This directory contains the complete extraction of the WorkTrack project code, organized into separate text files for easy reference and documentation.

## Extraction Date
July 28, 2026 - 14:47:51

## Extracted Files

### 1. **00_PROJECT_STRUCTURE_SUMMARY.txt** (11KB)
Complete directory structure of the entire WorkTrack project including:
- Backend (Go) structure
- Frontend Admin Dashboard (Vue.js) structure  
- Frontend Client Portal (Vue.js) structure
- Frontend Worker PWA (Vue.js) structure
- Root configuration files

### 2. **01_BACKEND_GO_FILES.txt** (247KB)
All backend Go source files including:
- Main application entry point (`cmd/api/main.go`)
- Configuration (`internal/config/`)
- Database layer (`internal/database/`)
- API handlers (`internal/handlers/`)
- Internationalization (`internal/i18n/`)
- Middleware (`internal/middleware/`)
- Data models (`internal/models/`)
- Router (`internal/router/`)
- Business services (`internal/services/`)
- Utility functions (`pkg/utils/`)
- Test files (`tests/`)

### 3. **02_FRONTEND_ADMIN_DASHBOARD.txt** (470KB)
Complete Admin Dashboard frontend code including:
- Vue.js components (`src/components/`)
- Views and pages (`src/views/`)
- API services (`src/services/`)
- State management (`src/store/`)
- Routing (`src/router/`)
- Internationalization files (`src/i18n/`, `src/locales/`)
- Styling (`src/styles/`)
- Configuration files (`package.json`, `vite.config.js`)
- PWA manifest and service worker

### 4. **03_FRONTEND_CLIENT_PORTAL.txt** (157KB)
Complete Client Portal frontend code including:
- Vue.js components for client functionality
- Service request forms
- Rating system
- Service history tracking
- Reports generation
- API integration
- Multi-language support

### 5. **04_FRONTEND_WORKER_PWA.txt** (231KB)
Complete Worker PWA frontend code including:
- Attendance tracking system
- Task management interface
- Geolocation features
- Photo upload functionality
- Offline PWA capabilities
- Real-time notifications
- Multi-language support

### 6. **05_CONFIG_AND_SETUP.txt** (1.6MB)
All configuration and setup files including:
- Package.json files (all dependencies)
- Environment configuration files
- Docker configuration
- Render deployment configs
- Markdown documentation files
- Translation files
- Build and deployment scripts

### 7. **06_DATABASE_FILES.txt** (133KB)
Database schema and migration files including:
- Complete database schema
- Migration scripts
- Fix scripts for data corrections
- SQL setup files

## Features of the Extraction Script

The extraction script (`extract_project_structure.py`) provides:

### Smart Filtering
- Excludes common development directories (`node_modules`, `.git`, `dist`, etc.)
- Only includes relevant source files
- Handles different file encodings (UTF-8, Latin-1)

### Organized Output
- Files are grouped by component type
- Each file is clearly labeled with its path
- File contents are separated by visual dividers
- Includes file count and size information

### Error Handling
- Gracefully handles files that can't be read
- Reports errors without stopping the extraction
- Supports multiple text encodings

## How to Use the Extraction Script

### Run Extraction
```bash
python3 extract_project_structure.py
```

### Customization
You can modify the script to:
- Add or exclude specific file types
- Change the output directory
- Modify the filtering criteria
- Add custom formatting

### Script Location
The main script is located at: `/home/dev-bit/project/worktrack/extract_project_structure.py`

## Project Structure Overview

### Backend (Go)
- **Framework**: Custom Go HTTP server
- **Database**: PostgreSQL with Redis caching
- **Architecture**: Clean architecture with handlers, services, and models
- **Features**: 
  - RESTful API
  - WebSocket support for real-time updates
  - JWT authentication
  - Role-based access control
  - Geofence and location tracking
  - Multi-language support

### Frontend Applications
All three frontend applications use:
- **Framework**: Vue.js 3 with Composition API
- **Build Tool**: Vite
- **State Management**: Pinia (Vue Store)
- **Routing**: Vue Router
- **Internationalization**: Vue I18n
- **Styling**: CSS with custom design tokens
- **PWA**: Service workers and manifests for offline capability

#### Admin Dashboard
- Employee management
- Work site management  
- Service request handling
- Real-time monitoring with maps
- Reporting and analytics
- Task management

#### Client Portal
- Service request submission
- Request tracking
- Service rating
- History and reports
- Profile management

#### Worker PWA
- Attendance tracking (Check In/Out)
- Task management
- Location tracking
- Photo uploads
- Offline functionality
- Real-time notifications

## File Organization in Text Files

Each text file follows this structure:
```
================================================================================
COMPONENT NAME
Extracted: [Date]
================================================================================

Total files found: [Number]

================================================================================

FILE 1/[Total]: [Relative Path]
================================================================================

[File Content]

FILE 2/[Total]: [Relative Path]
================================================================================

[File Content]

... and so on
```

## Benefits of This Extraction

1. **Easy Code Review**: All code in one place for quick reference
2. **Documentation**: Complete project documentation in text format
3. **Backup**: Text-based backup of source code
4. **Search**: Easy to search across all files
5. **Sharing**: Simple format to share with team members
6. **Analysis**: Can be used for code analysis and metrics

## Notes

- Binary files (images, compiled executables) are not included
- Node modules are excluded to reduce file size
- Build artifacts (dist, build directories) are excluded
- File sizes are approximate and may vary
- Some very large files may be truncated in display

## Technical Details

- **Script Language**: Python 3
- **Encoding**: UTF-8 with fallback to Latin-1
- **Line Endings**: Unix-style (\n)
- **File Paths**: Relative to project root
- **Timestamps**: ISO 8601 format

## Support

For questions about the extraction or the WorkTrack project structure, refer to:
- Project README.md
- Individual component documentation
- Technical guides in the root directory

---

**Extraction completed successfully on July 28, 2026**
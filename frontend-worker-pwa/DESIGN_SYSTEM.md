# WorkTrack Worker PWA - Design System (2024-2026 Standards)

## 📋 Table of Contents
1. [Design Philosophy](#design-philosophy)
2. [Color System](#color-system)
3. [Typography](#typography)
4. [Icons](#icons)
5. [Spacing & Layout](#spacing--layout)
6. [Components](#components)
7. [Page-by-Page Design](#page-by-page-design)
8. [Animations & Interactions](#animations--interactions)
9. [Responsive Design](#responsive-design)
10. [Dark Mode](#dark-mode)

---

## 🎨 Design Philosophy

### Core Principles (Based on 2024-2026 Standards)

1. **Professional Simplicity**
   - Clean, uncluttered interfaces
   - Purposeful use of color and decoration
   - High information density without overwhelm

2. **Visual Hierarchy**
   - Clear distinction between primary, secondary, and tertiary elements
   - Typography-driven design
   - Strategic use of size, weight, and color

3. **Modern Aesthetics**
   - Subtle gradients (not overwhelming)
   - Soft shadows with depth
   - Rounded corners (modern, not cartoonish)
   - Glassmorphism effects where appropriate

4. **Performance-First**
   - Minimal animations (60fps)
   - Optimized images and icons
   - CSS-based effects over JavaScript

5. **Accessibility**
   - WCAG 2.1 AA compliant
   - High contrast ratios (4.5:1 minimum)
   - Keyboard navigation support
   - Screen reader friendly

---

## 🎨 Color System

### Primary Palette (Professional & Trustworthy)

```css
/* Primary - Deep Blue-Indigo */
--primary-50: #EEF2FF
--primary-100: #E0E7FF
--primary-200: #C7D2FE
--primary-300: #A5B4FC
--primary-400: #818CF8
--primary-500: #6366F1  /* Main Brand Color */
--primary-600: #4F46E5
--primary-700: #4338CA
--primary-800: #3730A3
--primary-900: #312E81
```

### Semantic Colors

```css
/* Success - Professional Green */
--success-50: #ECFDF5
--success-100: #D1FAE5
--success-500: #10B981
--success-600: #059669
--success-700: #047857

/* Warning - Amber (Not Orange) */
--warning-50: #FFFBEB
--warning-100: #FEF3C7
--warning-500: #F59E0B
--warning-600: #D97706
--warning-700: #B45309

/* Error - Professional Red */
--error-50: #FEF2F2
--error-100: #FEE2E2
--error-500: #EF4444
--error-600: #DC2626
--error-700: #B91C1C

/* Info - Sky Blue */
--info-50: #F0F9FF
--info-100: #E0F2FE
--info-500: #0EA5E9
--info-600: #0284C7
--info-700: #0369A1
```

### Neutral Scale (Warm Gray - More Human)

```css
--gray-50: #FAFAF9
--gray-100: #F5F5F4
--gray-200: #E7E5E4
--gray-300: #D6D3D1
--gray-400: #A8A29E
--gray-500: #78716C
--gray-600: #57534E
--gray-700: #44403C
--gray-800: #292524
--gray-900: #1C1917
```

### Surface Colors

```css
--background: #FAFAF9      /* Light gray background */
--surface: #FFFFFF         /* Card background */
--surface-elevated: #F5F5F4  /* Elevated card */
--border: #E7E5E4         /* Subtle borders */
--border-strong: #D6D3D1   /* Strong borders */
```

### Text Colors

```css
--text-primary: #1C1917    /* Headings, important text */
--text-secondary: #57534E  /* Body text */
--text-tertiary: #A8A29E   /* Captions, hints */
--text-inverse: #FFFFFF     /* On dark backgrounds */
```

---

## 🔤 Typography

### Font Stack

```css
/* Primary - System Fonts */
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;

/* Arabic Support */
font-family: 'Cairo', 'Segoe UI', Roboto, sans-serif;
```

### Type Scale

```css
--text-xs: 0.75rem      /* 12px - Captions */
--text-sm: 0.875rem     /* 14px - Body small */
--text-base: 1rem       /* 16px - Body */
--text-lg: 1.125rem     /* 18px - Subheadings */
--text-xl: 1.25rem      /* 20px - Headings small */
--text-2xl: 1.5rem      /* 24px - Headings */
--text-3xl: 1.875rem    /* 30px - Section titles */
--text-4xl: 2.25rem     /* 36px - Page titles */
```

### Font Weights

```css
--font-normal: 400
--font-medium: 500
--font-semibold: 600
--font-bold: 700
```

### Line Heights

```css
--leading-tight: 1.25
--leading-normal: 1.5
--leading-relaxed: 1.75
```

---

## 🎯 Icons

### Icon System

**Use SVG Icons Only** - NO Emojis

**Recommended Libraries:**
- **Lucide Icons** (Modern, clean, consistent)
- **Heroicons** (Professional, from Tailwind team)
- **Feather Icons** (Simple, elegant)

**Icon Guidelines:**
- Size: 16px, 20px, 24px, 32px
- Stroke width: 2px (consistent)
- Color: Inherit or semantic colors
- No background decorations unless needed

**Icon Usage Examples:**

```html
<!-- Header Icon -->
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
  <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
</svg>

<!-- Status Icon (Success) -->
<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
  <polyline points="22 4 12 14.01 9 11.01"/>
</svg>
```

---

## 📐 Spacing & Layout

### Spacing Scale

```css
--space-0: 0
--space-1: 0.25rem    /* 4px */
--space-2: 0.5rem     /* 8px */
--space-3: 0.75rem    /* 12px */
--space-4: 1rem       /* 16px */
--space-5: 1.25rem    /* 20px */
--space-6: 1.5rem     /* 24px */
--space-8: 2rem       /* 32px */
--space-10: 2.5rem    /* 40px */
--space-12: 3rem      /* 48px */
--space-16: 4rem      /* 64px */
```

### Border Radius

```css
--radius-sm: 0.375rem   /* 6px - Small elements */
--radius-md: 0.5rem     /* 8px - Cards, buttons */
--radius-lg: 0.75rem    /* 12px - Large cards */
--radius-xl: 1rem       /* 16px - Modals */
--radius-2xl: 1.5rem     /* 24px - Hero cards */
--radius-full: 9999px   /* Pill shapes */
```

### Shadows

```css
--shadow-xs: 0 1px 2px rgba(0, 0, 0, 0.05)
--shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06)
--shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)
--shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)
--shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.04)
```

---

## 🧩 Components

### Buttons

#### Primary Button
```css
.btn-primary {
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-600) 100%);
  color: white;
  padding: 12px 24px;
  border-radius: var(--radius-md);
  font-weight: var(--font-semibold);
  box-shadow: var(--shadow-md);
  transition: all 0.2s ease;
}

.btn-primary:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-1px);
}
```

#### Secondary Button
```css
.btn-secondary {
  background: var(--surface);
  color: var(--text-primary);
  border: 1px solid var(--border);
  padding: 12px 24px;
  border-radius: var(--radius-md);
  font-weight: var(--font-semibold);
  transition: all 0.2s ease;
}

.btn-secondary:hover {
  border-color: var(--primary-500);
  color: var(--primary-500);
}
```

### Cards

```css
.card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-sm);
  transition: all 0.2s ease;
}

.card:hover {
  box-shadow: var(--shadow-md);
  border-color: var(--border-strong);
}
```

### Badges

```css
.badge {
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
}

.badge-success {
  background: var(--success-100);
  color: var(--success-700);
}

.badge-warning {
  background: var(--warning-100);
  color: var(--warning-700);
}

.badge-error {
  background: var(--error-100);
  color: var(--error-700);
}
```

---

## 📄 Page-by-Page Design

### 1. Login Page (`LoginView.vue`)

#### Design Specifications

**Layout:**
- Centered card on screen
- Max width: 400px
- Full height: 100vh
- Background: Gradient from primary-50 to background

**Header:**
- Logo icon (SVG): 64x64px
- App name: text-3xl, font-bold
- Subtitle: text-base, text-secondary

**Form:**
- Phone input: Large (18px), centered
- Language selector: Pill buttons
- Submit button: Full width, primary

**Visual Style:**
- Glassmorphism card effect
- Subtle gradient background
- Professional, trustworthy feel

**Code Structure:**
```vue
<template>
  <div class="login-page">
    <div class="login-card">
      <!-- Logo -->
      <div class="login-header">
        <div class="logo-icon">
          <svg>...</svg>
        </div>
        <h1>WorkTrack</h1>
        <p>Employee Portal</p>
      </div>

      <!-- Language Selector -->
      <div class="language-selector">
        <button class="lang-btn active">🇸🇦 العربية</button>
        <button class="lang-btn">🇮🇱 עברית</button>
        <button class="lang-btn">🇬🇧 English</button>
      </div>

      <!-- Form -->
      <form class="login-form">
        <div class="form-group">
          <label class="form-label">
            <svg>...</svg>
            Phone Number
          </label>
          <input type="tel" class="form-input" placeholder="05xxxxxxxx" />
        </div>
        <button class="btn btn-primary btn-full">Login</button>
      </form>
    </div>
  </div>
</template>
```

---

### 2. Main App Structure (`App.vue`)

#### Header
- Brand icon (SVG): 40x40px
- Brand name + badge
- Avatar: 40x40px, circle
- Sticky position
- Subtle shadow

#### Bottom Navigation
- 3 tabs: Attendance, Notes, Profile
- Icon + label for each
- Active state: primary background
- Fixed position at bottom
- Safe area for iPhone notch

#### Main Content
- Padding: 24px
- Max width: 100% (responsive)
- No 480px limit

---

### 3. Attendance Page (`AttendanceView.vue`)

#### Page Header
```vue
<header class="page-header">
  <div class="header-icon">
    <svg>...</svg>
  </div>
  <div class="header-text">
    <h1>Attendance</h1>
    <p>Track your work hours</p>
  </div>
</header>
```

#### Hours Summary Cards
- 3 cards in horizontal layout
- Icon + value + label
- Different colors for today/week/month
- Responsive: stack on mobile

#### Timer Display (When Working)
- Large time display (text-4xl)
- Gradient background
- Pulsing animation
- Worksite name below

#### Worksite Selection
- Grid layout (2 columns on tablet, 1 on mobile)
- Card style for each worksite
- Active state: border + checkmark
- Hover: lift effect

#### Location Card
- Current location info
- Distance to worksite
- Status indicator (in/out range)
- GPS coordinates

---

### 4. Profile Page (`ProfileView.vue`)

#### Profile Header
- Large avatar: 72x72px
- Name: text-xl, bold
- Role: text-sm, secondary
- Gradient background

#### Language Selector
- Card with title
- Button-based selection
- Active state highlighted

#### Menu Items
- Icon + label
- Chevron for expandable
- Hover: background change
- Active: primary color

#### WhatsApp Button
- Floating action button
- Bottom-right corner
- Green gradient
- Pulse animation

---

### 5. Notes Page (`NotesView.vue`)

#### Page Header
- Title + subtitle
- Icon in title

#### Note Cards
- Avatar of sender
- Name + role
- Content text
- Worksite badge (if applicable)
- Read/unread indicator
- Date + updated time
- Mark as read button

#### Visual Hierarchy
- Unread: Left border + gradient background
- Read: Simple card
- Hover: Shadow + border highlight

---

### 6. Tasks Page (`TasksListView.vue`)

#### Page Header
- Title + badge count
- Breadcrumb style

#### Task Cards
- Status stripe (left edge)
- Title + badge
- Location icon + text
- Time icon + text
- Click to view details

#### Status Colors
- Pending: Gray
- In Progress: Amber
- Completed: Green
- Late: Red

---

## ✨ Animations & Interactions

### Animations

```css
/* Fade In */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Slide Up */
@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Pulse */
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Spin */
@keyframes spin {
  to { transform: rotate(360deg); }
}
```

### Interactions

**Hover Effects:**
- Cards: Lift (+2px) + shadow increase
- Buttons: Color shift + slight scale
- Links: Underline + color change

**Active States:**
- Buttons: Scale down (0.97)
- Cards: Border highlight

**Loading States:**
- Spinner (SVG)
- Skeleton screens
- Progress indicators

---

## 📱 Responsive Design

### Breakpoints

```css
/* Mobile First */
--mobile: 375px;
--mobile-large: 414px;
--tablet: 768px;
--tablet-large: 1024px;
--desktop: 1280px;
--desktop-large: 1440px;
```

### Mobile (< 768px)
- Single column layouts
- Touch-friendly targets (44px min)
- Simplified navigation
- Reduced padding

### Tablet (768px - 1024px)
- 2-column grids
- Enhanced spacing
- Larger touch targets

### Desktop (> 1024px)
- Multi-column layouts
- Hover interactions
- Max-width containers
- Side-by-side components

---

## 🌙 Dark Mode

### Color Overrides

```css
[data-theme="dark"] {
  --background: #1C1917;
  --surface: #292524;
  --surface-elevated: #44403C;
  --border: #44403C;
  --border-strong: #57534E;
  
  --text-primary: #FAFAF9;
  --text-secondary: #A8A29E;
  --text-tertiary: #78716C;
  
  --primary-500: #818CF8;  /* Lighter for dark mode */
  --primary-600: #A5B4FC;
}
```

### Adjustments
- Increase contrast ratios
- Adjust shadows (lighter)
- Modify gradients
- Ensure readability

---

## 🎯 Implementation Checklist

### Must-Have
- [ ] SVG icons (no emojis)
- [ ] Professional color palette
- [ ] Clear visual hierarchy
- [ ] Responsive design (mobile-first)
- [ ] Accessible (WCAG AA)
- [ ] Dark mode support
- [ ] Smooth animations (60fps)
- [ ] Loading states
- [ ] Error states
- [ ] Empty states

### Nice-to-Have
- [ ] Micro-interactions
- [ ] Haptic feedback (mobile)
- [ ] Skeleton screens
- [ ] Progressive enhancement
- [ ] Offline support indicators
- [ ] Toast notifications
- [ ] Confetti animations (success)

---

## 📚 References & Inspiration

### Design Systems to Study
1. **Linear** - Modern, dark mode, gradients
2. **Vercel** - Clean, professional, monochrome
3. **Stripe** - Trustworthy, subtle animations
4. **Notion** - Functional, emoji-as-icons (but professional)
5. **Figma** - Professional tools UI
6. **Apple** - Typography, spacing, simplicity

### Icon Libraries
1. **Lucide Icons** - Modern, consistent
2. **Heroicons** - Professional, Tailwind-based
3. **Feather Icons** - Simple, elegant

### Color Palettes
1. **Tailwind CSS** - Professional defaults
2. **Radix UI** - Accessible colors
3. **shadcn/ui** - Modern components

---

## 🚀 Next Steps

1. **Implement Design System**
   - Create `tokens.css` with all variables
   - Build component library
   - Document usage patterns

2. **Redesign Pages**
   - Start with Login (most important first impression)
   - Then App structure (navigation, header)
   - Then Attendance (main functionality)
   - Then remaining pages

3. **Test & Iterate**
   - Mobile testing (critical)
   - Accessibility testing
   - Performance testing
   - User feedback

4. **Polish**
   - Micro-interactions
   - Loading states
   - Error handling
   - Edge cases

---

*Last Updated: 2024*
*Based on 2024-2026 Design Standards*

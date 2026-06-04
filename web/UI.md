# Urja Design System — Material You (MD3) Light Theme

## Seed Color
`#6750A4` (Purple)

## Color Tokens

| Token | Value | Usage |
|:------|:------|:------|
| `md-primary` | `#6750A4` | Primary buttons, FABs, active states |
| `md-on-primary` | `#FFFFFF` | Text/icons on primary |
| `md-primary-container` | `#EADDFF` | Tonal buttons, chips, selected states |
| `md-on-primary-container` | `#21005D` | Text/icons on primary container |
| `md-secondary` | `#625B71` | Secondary actions |
| `md-on-secondary` | `#FFFFFF` | Text on secondary |
| `md-secondary-container` | `#E8DEF8` | Secondary tonal fills |
| `md-on-secondary-container` | `#1D192B` | Text on secondary container |
| `md-tertiary` | `#7D5260` | Tertiary accent |
| `md-on-tertiary` | `#FFFFFF` | Text on tertiary |
| `md-tertiary-container` | `#FFD8E4` | Tertiary tonal fills |
| `md-on-tertiary-container` | `#31111D` | Text on tertiary container |
| `md-error` | `#B3261E` | Error states |
| `md-on-error` | `#FFFFFF` | Text on error |
| `md-error-container` | `#F9DEDC` | Error container |
| `md-on-error-container` | `#410E0B` | Text on error container |
| `md-background` | `#FEF7FF` | Page background |
| `md-on-background` | `#1C1B1F` | Text on background |
| `md-surface` | `#FEF7FF` | Surface base |
| `md-on-surface` | `#1C1B1F` | Text on surface |
| `md-on-surface-variant` | `#49454F` | Secondary text, icons |
| `md-outline` | `#79747E` | Borders, dividers |
| `md-outline-variant` | `#CAC4D0` | Subtle borders |
| `md-surface-container-lowest` | `#FFFFFF` | Lowest surface |
| `md-surface-container-low` | `#F7F2FA` | Low surface |
| `md-surface-container` | `#F3EDF7` | Default container |
| `md-surface-container-high` | `#ECE6F0` | High container |
| `md-surface-container-highest` | `#E6E0E9` | Highest container |

## Typography

**Font:** Roboto (400, 500, 700) via `next/font/google`

| Level | Size | Weight | Usage |
|:------|:-----|:-------|:------|
| Display Large | 57px | 400 | Hero headlines |
| Display Medium | 45px | 400 | Section headers |
| Headline Large | 32px | 400 | Page titles |
| Headline Medium | 28px | 400 | Card titles |
| Title Large | 22px | 500 | Nav titles |
| Title Medium | 16px | 500 | Subtitles |
| Body Large | 16px | 400 | Primary body |
| Body Medium | 14px | 400 | Secondary body |
| Label Large | 14px | 500 | Buttons, tabs |
| Label Medium | 12px | 500 | Badges, captions |

## Elevation (Shadow System)

| Level | Shadow | Usage |
|:------|:-------|:------|
| `md-1` | `0 1px 3px 1px rgba(0,0,0,0.15), 0 1px 2px rgba(0,0,0,0.3)` | Cards, buttons |
| `md-2` | `0 2px 6px 2px rgba(0,0,0,0.15), 0 1px 2px rgba(0,0,0,0.3)` | Hover cards, dropdowns |
| `md-3` | `0 4px 8px 3px rgba(0,0,0,0.15), 0 1px 3px rgba(0,0,0,0.3)` | FABs, modals |
| `md-4` | `0 6px 10px 4px rgba(0,0,0,0.15), 0 2px 3px rgba(0,0,0,0.3)` | Navigation drawers |

## Shape (Border Radius)

| Element | Radius |
|:--------|:-------|
| Buttons | `rounded-full` (pill) |
| Cards | `rounded-3xl` (24px) |
| Containers | `rounded-2xl` (16px) |
| Inputs | `rounded-xl` (12px) |
| Badges/chips | `rounded-full` |
| Icons | `rounded-xl` (12px) |

## Component Patterns

### Buttons
- **Filled:** `bg-md-primary text-md-on-primary rounded-full px-6 py-2.5`
- **Tonal:** `bg-md-secondary-container text-md-on-secondary-container rounded-full`
- **Outlined:** `border border-md-outline text-md-primary rounded-full`
- **Text:** `text-md-primary rounded-full`

### Cards
- Background: `bg-md-surface-container`
- Radius: `rounded-3xl` (24px)
- Shadow: `shadow-md-1` default, `shadow-md-2` on hover
- No visible border by default

### Inputs
- Background: transparent
- Border: `border-md-outline`
- Focus: `border-md-primary ring-1 ring-md-primary`
- Label floats above on focus (MD3 pattern)
- Radius: `rounded-xl`

### State Layers
- Hover: 8% opacity overlay of content color
- Focus: 12% opacity overlay
- Pressed: 12% opacity overlay
- Dragged: 16% opacity overlay

## Landing Page Elements

### Organic Blur Shapes
Decorative background blobs using the palette:
```
absolute w-[600px] h-[600px] rounded-full blur-[150px] opacity-20
bg-md-primary-container (or bg-md-tertiary-container)
```

### Social Proof
Gym count badges, customer logos, trust indicators

### Pricing Cards
Three tiers — Basic / Pro / Enterprise
Featured tier uses `bg-md-primary text-md-on-primary` with `shadow-md-3`

# Tech Stack and Dependencies

## Core Technologies

This project is a **Go web application** using modern web technologies with a focus on server-side rendering and component-based UI development.

### Backend
- **Go 1.24.5** - Primary programming language
- **net/http** - Standard HTTP server (no external web framework)
- **slog** - Structured logging (Go standard library)
- **Module**: `github.com/JacobSchroder/jup`

### Frontend & UI
- **Templ** (`v0.3.906`) - Go templating engine for type-safe HTML components
- **Tailwind CSS** (`v4.0.9`) - Utility-first CSS framework with custom theming
- **TemplUI** - Component library system for Go/Templ (configured via `.templui.json`)
- **tailwind-merge-go** (`v0.2.1`) - Utility for merging Tailwind CSS classes

### Development Tools
- **Air** - Live reload for Go applications (configured via `.air.toml`)
- **Make** - Build automation and task management
- **Docker & Docker Compose** - Containerization and local development

## Project Structure

```
jup/
├── cmd/                    # Application entry points
├── server/                 # Server setup, handlers, and routing
│   ├── handlers/          # HTTP request handlers
│   └── routes/            # Route definitions
├── pages/                 # Page-level Templ templates
├── templates/             # Shared/layout Templ templates
├── components/            # Reusable UI components (TemplUI)
├── utils/                 # Shared utility functions
├── assets/                # Static assets
│   ├── css/              # Stylesheets (Tailwind)
│   └── js/               # JavaScript files
├── dev/                  # Development configuration
└── tmp/                  # Temporary build files
```

## Component System

- **TemplUI Components**: Pre-built, customizable UI components in `/components/`
- **Component Types**: Button, Card, Modal, Form inputs, Navigation, Charts, etc.
- **Styling**: Components use Tailwind CSS with custom design tokens
- **Type Safety**: Full Go type safety for component props and rendering

## Build System

### Development
- `make dev` - Start development server with Air live reload
- `make templ-watch` - Watch and regenerate Templ templates
- `make tailwind-watch` - Watch and rebuild CSS

### Production
- `make build` - Full production build (CSS minification + Templ generation)
- `make docker-build` - Docker container build

## Configuration Files

- `go.mod/go.sum` - Go module dependencies
- `.templui.json` - TemplUI component library configuration
- `.air.toml` - Live reload configuration
- `Makefile` - Build and development commands
- `assets/css/input.css` - Tailwind CSS configuration with custom theme

## Key Features

- **Server-Side Rendering**: Full SSR with Go templates
- **Component-Based**: Reusable UI components with TemplUI
- **Type Safety**: End-to-end type safety from Go to HTML
- **Modern CSS**: Tailwind v4 with custom design system
- **Live Development**: Hot reload for Go, Templ, and CSS
- **Production Ready**: Optimized builds with minification

## Dependencies

### Direct Dependencies
- `github.com/a-h/templ v0.3.906` - Template engine
- `github.com/Oudwins/tailwind-merge-go v0.2.1` - CSS class merging

### Development Dependencies
- `air` - Live reload
- `templ` - Template compiler
- `tailwindcss` - CSS framework (standalone binary)
- `staticcheck` - Go static analysis

## Server Configuration

- **Host**: localhost
- **Port**: 8082
- **Graceful Shutdown**: Signal-based context cancellation
- **Logging**: Structured JSON logging with slog

When working with this project, always consider the type-safe component system and the integration between Go, Templ, and Tailwind CSS.

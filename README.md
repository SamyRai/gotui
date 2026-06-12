# GoTUI - Advanced Terminal User Interface Framework

A comprehensive Go-based Terminal User Interface (TUI) framework built on Bubble Tea and Bubbles libraries. Features modern TUI components, responsive layouts, advanced theming, and extensive customization options for building sophisticated terminal applications.

## Overview

GoTUI provides:

- **🎨 Modern TUI Components**: Beautiful, responsive terminal UI components
- **📱 Responsive Design**: Adaptive layouts for different terminal sizes
- **🎭 Advanced Theming**: Comprehensive theming and styling system
- **⚡ High Performance**: Optimized for smooth terminal interactions
- **🔧 Extensible**: Plugin-based architecture for custom components
- **📊 Rich Components**: Charts, tables, forms, and interactive elements

## Features

### Core TUI Features
- **Modern Components**: Buttons, inputs, lists, tables, charts, and more
- **Responsive Layouts**: Adaptive layouts that work on any terminal size
- **Advanced Theming**: Comprehensive theming with dark/light modes
- **Interactive Elements**: Mouse support, keyboard shortcuts, and gestures
- **Animation Support**: Smooth animations and transitions
- **Accessibility**: Screen reader support and keyboard navigation

### Advanced Components
- **Data Visualization**: Charts, graphs, and statistical displays
- **Form Handling**: Complex forms with validation and error handling
- **File Management**: File browsers, directory trees, and file operations
- **Text Editing**: Rich text editing with syntax highlighting
- **Progress Tracking**: Progress bars, spinners, and status indicators
- **Modal Dialogs**: Popup dialogs, confirmations, and input prompts

### Developer Experience
- **Hot Reloading**: Development-time hot reloading for rapid iteration
- **Component Library**: Extensive library of pre-built components
- **Documentation**: Comprehensive documentation and examples
- **Testing Tools**: Built-in testing utilities and mock components
- **CLI Tools**: Command-line tools for component generation and management
- **IDE Integration**: VS Code extensions and IDE support

## Architecture

### System Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                    Terminal Applications                   │
├─────────────────────────────────────────────────────────────┤
│  CLI Tools  │  Interactive Apps  │  Dashboards  │  Games   │
├─────────────────────────────────────────────────────────────┤
│                    GoTUI Framework                        │
├─────────────────────────────────────────────────────────────┤
│  Component Engine  │  Layout Engine  │  Theme Engine  │  Event │
├─────────────────────────────────────────────────────────────┤
│                    Core Libraries                          │
├─────────────────────────────────────────────────────────────┤
│  Bubble Tea  │  Bubbles  │  Lip Gloss  │  Termenv  │  Tcell │
├─────────────────────────────────────────────────────────────┤
│                    Terminal Layer                          │
├─────────────────────────────────────────────────────────────┤
│  Terminal Emulator  │  Screen Buffer  │  Input Events  │  Colors │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

#### Core Framework
- **Go 1.21+**: Primary development language
- **Bubble Tea**: TUI framework and state management
- **Bubbles**: Component library
- **Lip Gloss**: Styling and layout
- **Termenv**: Terminal environment detection
- **Tcell**: Low-level terminal control

#### Additional Libraries
- **Charm**: Charm ecosystem integration
- **Viper**: Configuration management
- **Cobra**: CLI framework
- **Testify**: Testing framework
- **Go Modules**: Dependency management

## Installation

### Prerequisites
- **Go 1.21+**
- **Terminal Emulator** (iTerm2, Terminal.app, etc.)
- **Git** (for version control)

### Quick Setup

```bash
# Clone repository
git clone <repository-url>
cd gotui

# Install dependencies
go mod download

# Set up development environment
go run cmd/setup/main.go

# Start development server
go run cmd/dev/main.go
```

### Environment Configuration

```bash
# Terminal Configuration
TERM=xterm-256color
COLORTERM=truecolor

# Development Configuration
GOTUI_DEV_MODE=true
GOTUI_HOT_RELOAD=true
GOTUI_LOG_LEVEL=debug

# Theme Configuration
GOTUI_THEME=default
GOTUI_COLOR_SCHEME=auto
GOTUI_FONT_SIZE=medium
```

## Usage

### Basic TUI Application

#### Simple Application
```go
package main

import (
    "fmt"
    "github.com/gotui/gotui"
    "github.com/gotui/gotui/components"
)

func main() {
    // Create new TUI application
    app := gotui.New()
    
    // Add components
    app.AddComponent(components.NewTitle("Welcome to GoTUI"))
    app.AddComponent(components.NewButton("Click Me", func() {
        fmt.Println("Button clicked!")
    }))
    app.AddComponent(components.NewInput("Enter text:", func(text string) {
        fmt.Printf("You entered: %s\n", text)
    }))
    
    // Start application
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

#### Advanced Application
```go
package main

import (
    "github.com/gotui/gotui"
    "github.com/gotui/gotui/components"
    "github.com/gotui/gotui/layouts"
    "github.com/gotui/gotui/themes"
)

func main() {
    // Create application with theme
    app := gotui.NewWithTheme(themes.Dark())
    
    // Create layout
    layout := layouts.NewGrid(2, 2)
    
    // Add components to layout
    layout.AddComponent(0, 0, components.NewChart(components.LineChart))
    layout.AddComponent(0, 1, components.NewTable(data))
    layout.AddComponent(1, 0, components.NewForm(formFields))
    layout.AddComponent(1, 1, components.NewList(items))
    
    // Set layout
    app.SetLayout(layout)
    
    // Add event handlers
    app.OnKeyPress("q", func() {
        app.Quit()
    })
    
    // Start application
    app.Run()
}
```

### Advanced Components

#### Data Visualization
```go
import "github.com/gotui/gotui/components"

// Create chart
chart := components.NewChart(components.BarChart)
chart.SetData(chartData)
chart.SetTitle("Sales Data")
chart.SetColors([]string{"#FF6B6B", "#4ECDC4", "#45B7D1"})

// Create table
table := components.NewTable()
table.SetHeaders([]string{"Name", "Age", "City"})
table.SetData(tableData)
table.SetSortable(true)
table.SetSelectable(true)

// Create progress bar
progress := components.NewProgressBar()
progress.SetValue(0.75)
progress.SetLabel("Loading...")
progress.SetColor("#4ECDC4")
```

#### Form Handling
```go
// Create form
form := components.NewForm()

// Add form fields
form.AddField(components.NewInputField("Name", "Enter your name", true))
form.AddField(components.NewSelectField("Country", countries, "Select country"))
form.AddField(components.NewCheckboxField("Subscribe", false))
form.AddField(components.NewTextAreaField("Message", "Enter your message", 5))

// Set form validation
form.SetValidator(func(data map[string]interface{}) error {
    if data["name"].(string) == "" {
        return fmt.Errorf("Name is required")
    }
    return nil
})

// Set form submit handler
form.OnSubmit(func(data map[string]interface{}) {
    fmt.Printf("Form submitted: %+v\n", data)
})
```

#### Interactive Elements
```go
// Create list with custom items
list := components.NewList()
list.AddItem("Option 1", "Description 1", "icon1")
list.AddItem("Option 2", "Description 2", "icon2")
list.AddItem("Option 3", "Description 3", "icon3")

// Set list behavior
list.SetSelectable(true)
list.SetMultiSelect(false)
list.SetSearchable(true)

// Set list event handlers
list.OnSelect(func(item components.ListItem) {
    fmt.Printf("Selected: %s\n", item.Text)
})

// Create modal dialog
modal := components.NewModal()
modal.SetTitle("Confirmation")
modal.SetContent("Are you sure you want to continue?")
modal.AddButton("Yes", func() {
    // Handle confirmation
    modal.Close()
})
modal.AddButton("No", func() {
    modal.Close()
})
```

### Theming and Styling

#### Custom Theme
```go
import "github.com/gotui/gotui/themes"

// Create custom theme
theme := themes.NewTheme()
theme.SetBackgroundColor("#1E1E1E")
theme.SetForegroundColor("#FFFFFF")
theme.SetAccentColor("#4ECDC4")
theme.SetBorderColor("#333333")
theme.SetErrorColor("#FF6B6B")
theme.SetSuccessColor("#4ECDC4")
theme.SetWarningColor("#FFA500")

// Apply theme to application
app := gotui.NewWithTheme(theme)
```

#### Component Styling
```go
// Style individual components
button := components.NewButton("Click Me", nil)
button.SetStyle(components.Style{
    Background: "#4ECDC4",
    Foreground: "#FFFFFF",
    Border:     "#45B7D1",
    Padding:    components.Padding{Top: 1, Bottom: 1, Left: 2, Right: 2},
})

// Style with CSS-like properties
input := components.NewInput("Enter text:", nil)
input.SetStyle(components.Style{
    Background: "#2D2D2D",
    Foreground: "#FFFFFF",
    Border:     "#4ECDC4",
    BorderStyle: components.BorderRounded,
})
```

### Layout Management

#### Grid Layout
```go
import "github.com/gotui/gotui/layouts"

// Create grid layout
grid := layouts.NewGrid(3, 2)
grid.SetGap(1)

// Add components to grid
grid.AddComponent(0, 0, header)
grid.AddComponent(0, 1, sidebar)
grid.AddComponent(1, 0, mainContent)
grid.AddComponent(1, 1, sidebar)
grid.AddComponent(2, 0, footer)
grid.AddComponent(2, 1, footer)
```

#### Flex Layout
```go
// Create flex layout
flex := layouts.NewFlex()
flex.SetDirection(layouts.Row)
flex.SetJustifyContent(layouts.SpaceBetween)
flex.SetAlignItems(layouts.Center)

// Add components to flex
flex.AddComponent(logo)
flex.AddComponent(navigation)
flex.AddComponent(userMenu)
```

#### Stack Layout
```go
// Create stack layout
stack := layouts.NewStack()
stack.AddComponent(background)
stack.AddComponent(modal)
stack.AddComponent(tooltip)
```

## API Reference

### Core Components

#### Buttons
```go
// Basic button
button := components.NewButton("Click Me", func() {
    // Handle click
})

// Styled button
button := components.NewButton("Styled", handler)
button.SetStyle(components.Style{
    Background: "#4ECDC4",
    Foreground: "#FFFFFF",
})

// Button with icon
button := components.NewButtonWithIcon("Save", "💾", handler)
```

#### Inputs
```go
// Text input
input := components.NewInput("Enter text:", func(text string) {
    // Handle input
})

// Password input
password := components.NewPasswordInput("Password:", func(password string) {
    // Handle password
})

// Number input
number := components.NewNumberInput("Age:", func(value int) {
    // Handle number
})
```

#### Lists
```go
// Basic list
list := components.NewList()
list.AddItem("Item 1", "Description 1", "icon1")
list.AddItem("Item 2", "Description 2", "icon2")

// List with search
list.SetSearchable(true)
list.OnSearch(func(query string) {
    // Handle search
})
```

#### Tables
```go
// Create table
table := components.NewTable()
table.SetHeaders([]string{"Name", "Age", "City"})
table.SetData([][]string{
    {"John", "25", "New York"},
    {"Jane", "30", "London"},
})

// Table with sorting
table.SetSortable(true)
table.OnSort(func(column int, ascending bool) {
    // Handle sort
})
```

### Layout Components

#### Grid Layout
```go
// Create grid
grid := layouts.NewGrid(2, 3)
grid.SetGap(1)
grid.SetPadding(components.Padding{All: 1})

// Add components
grid.AddComponent(0, 0, component1)
grid.AddComponent(0, 1, component2)
grid.AddComponent(1, 0, component3)
```

#### Flex Layout
```go
// Create flex
flex := layouts.NewFlex()
flex.SetDirection(layouts.Row)
flex.SetJustifyContent(layouts.SpaceBetween)
flex.SetAlignItems(layouts.Center)

// Add components
flex.AddComponent(component1)
flex.AddComponent(component2)
```

### Event Handling

#### Keyboard Events
```go
// Handle key press
app.OnKeyPress("q", func() {
    app.Quit()
})

// Handle key combination
app.OnKeyPress("ctrl+c", func() {
    app.Quit()
})

// Handle any key
app.OnKeyPress("", func(key string) {
    fmt.Printf("Key pressed: %s\n", key)
})
```

#### Mouse Events
```go
// Handle mouse click
app.OnMouseClick(func(x, y int, button string) {
    fmt.Printf("Clicked at (%d, %d) with %s button\n", x, y, button)
})

// Handle mouse move
app.OnMouseMove(func(x, y int) {
    // Handle mouse movement
})
```

## Configuration

### Application Configuration
```yaml
# config/app.yml
application:
  name: "GoTUI App"
  version: "1.0.0"
  theme: "dark"
  hot_reload: true
  log_level: "info"
  
  window:
    width: 80
    height: 24
    resizable: true
    fullscreen: false
  
  components:
    default_style:
      background: "#1E1E1E"
      foreground: "#FFFFFF"
      border: "#333333"
      padding: 1
    
    animations:
      enabled: true
      duration: 200
      easing: "ease-in-out"
```

### Theme Configuration
```yaml
# config/themes.yml
themes:
  dark:
    background: "#1E1E1E"
    foreground: "#FFFFFF"
    accent: "#4ECDC4"
    border: "#333333"
    error: "#FF6B6B"
    success: "#4ECDC4"
    warning: "#FFA500"
  
  light:
    background: "#FFFFFF"
    foreground: "#000000"
    accent: "#007ACC"
    border: "#CCCCCC"
    error: "#DC3545"
    success: "#28A745"
    warning: "#FFC107"
```

## Advanced Features

### Custom Components
```go
// Create custom component
type CustomComponent struct {
    *components.BaseComponent
    data string
}

func NewCustomComponent(data string) *CustomComponent {
    return &CustomComponent{
        BaseComponent: components.NewBaseComponent(),
        data: data,
    }
}

func (c *CustomComponent) Render() string {
    return fmt.Sprintf("Custom: %s", c.data)
}

// Use custom component
custom := NewCustomComponent("Hello World")
app.AddComponent(custom)
```

### Plugin System
```go
// Create plugin
type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "MyPlugin"
}

func (p *MyPlugin) Initialize(app *gotui.App) error {
    // Plugin initialization
    return nil
}

func (p *MyPlugin) Components() []components.Component {
    return []components.Component{
        NewCustomComponent("Plugin Component"),
    }
}

// Register plugin
app.RegisterPlugin(&MyPlugin{})
```

### Animation System
```go
// Create animation
animation := components.NewAnimation()
animation.SetDuration(500)
animation.SetEasing(components.EaseInOut)
animation.SetProperty("opacity", 0.0, 1.0)

// Apply animation to component
component.SetAnimation(animation)
```

## Development

### Project Structure
```
gotui/
├── cmd/                   # Command-line applications
│   ├── dev/              # Development server
│   ├── build/            # Build tool
│   └── generate/         # Component generator
├── internal/             # Internal packages
│   ├── components/       # TUI components
│   ├── layouts/          # Layout engines
│   ├── themes/           # Theming system
│   ├── events/           # Event handling
│   └── renderer/         # Rendering engine
├── pkg/                  # Public packages
├── examples/             # Example applications
├── docs/                 # Documentation
└── tests/                # Test files
```

### Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run component tests
go test -tags=components ./...

# Run integration tests
go test -tags=integration ./...
```

## Contributing

### Development Setup
```bash
# Fork and clone repository
git clone <your-fork-url>
cd gotui

# Install dependencies
go mod download

# Set up development environment
go run cmd/dev/main.go

# Start development server
go run cmd/dev/main.go
```

### Code Standards
- Follow Go best practices and idioms
- Write comprehensive tests
- Document all public APIs
- Follow TUI design principles
- Update documentation for changes

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions, issues, or contributions:

- **Documentation**: Check the docs/ directory
- **Issues**: Report bugs and feature requests on GitHub
- **Discussions**: Join community discussions
- **Email**: Contact support@gotui.com

---

*Building beautiful terminal user interfaces with Go* 🎨📱⚡
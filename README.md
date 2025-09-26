# td

A modern, efficient terminal-based task management application written in Go. td provides a clean and intuitive interface for managing your todo list with features like task prioritization, undo/redo operations, and data persistence.

## Features

- **Terminal UI**: Clean, responsive interface using Bubble Tea
- **Task Prioritization**: Organize tasks by priority (None, Low, Medium, High)
- **Undo/Redo**: Full undo/redo support for all operations
- **Data Persistence**: Automatic saving to JSON file
- **Keyboard Shortcuts**: Vim-inspired navigation
- **Filtering**: Filter tasks by priority level
- **Cross-platform**: Works on macOS, Linux, and Windows

## Installation

### macOS
```bash
brew tap voioo/homebrew-tap
brew install td-tui
```

<details>
<summary>Manual Installation</summary>

```bash
# For Apple Silicon Macs:
curl -LO https://github.com/voioo/td/releases/latest/download/td_darwin_arm64.tar.gz
sudo tar xf td_darwin_arm64.tar.gz -C /usr/local/bin td

# For Intel Macs:
curl -LO https://github.com/voioo/td/releases/latest/download/td_darwin_amd64.tar.gz
sudo tar xf td_darwin_amd64.tar.gz -C /usr/local/bin td
```
</details>

### Arch Linux
```bash
yay -S td-tui
```
or
```bash
paru -S td-tui
```
or
```bash
git clone https://aur.archlinux.org/td-tui.git
cd td-tui
makepkg -si
```

## Ubuntu/Debian
```bash
curl -LO https://github.com/voioo/td/releases/latest/download/td_linux_amd64.tar.gz
sudo tar xf td_linux_amd64.tar.gz -C /usr/local/bin td
```

### RHEL/Fedora/CentOS
```bash
curl -LO https://github.com/voioo/td/releases/latest/download/td_linux_amd64.tar.gz
sudo tar xf td_linux_amd64.tar.gz -C /usr/local/bin td
```

### Windows
Open PowerShell as Administrator and run:
```powershell
irm https://raw.githubusercontent.com/voioo/td/main/install.ps1 | iex
```

<details>
<summary>Manual Installation</summary>

```powershell
# For AMD64 systems:
Invoke-WebRequest -Uri https://github.com/voioo/td/releases/latest/download/td_windows_amd64.zip -OutFile td.zip
Expand-Archive td.zip -DestinationPath "$env:LOCALAPPDATA\Programs\td"
$env:Path += ";$env:LOCALAPPDATA\Programs\td"

# For ARM64 systems:
Invoke-WebRequest -Uri https://github.com/voioo/td/releases/latest/download/td_windows_arm64.zip -OutFile td.zip
Expand-Archive td.zip -DestinationPath "$env:LOCALAPPDATA\Programs\td"
$env:Path += ";$env:LOCALAPPDATA\Programs\td"
```
</details>

You can also check the releases page on Github and download the one you need.

## Usage

### Basic Operations

- `a` - Add new task
- `d` - Delete selected task
- `enter` - Mark task as complete/incomplete
- `→` or `l` - Edit selected task
- `p` - Cycle task priority
- `1-4` - Set priority directly (1=none, 2=low, 3=medium, 4=high)
- `f` - Filter tasks by priority
- `t` - Toggle between active/completed tasks

### Navigation

- `↑` or `k` - Move up
- `↓` or `j` - Move down
- `←` or `h` - Move left
- `→` or `l` - Move right (or edit task)
- `home` or `g` - Go to top
- `end` or `G` - Go to bottom

### Other

- `?` - Show/hide help
- `C` - Clear all completed tasks
- `ctrl+u` - Undo last action
- `ctrl+r` - Redo last action
- `q` or `ctrl+c` - Quit

### Priority Levels

Tasks are automatically sorted by priority (high to low) and then by creation time. The priority indicators are:

- ○ - No priority
- ● (gray) - Low priority
- ● (yellow) - Medium priority
- ● (red) - High priority

### Configuration

td can be configured via a JSON or YAML config file at `~/.config/td/config.json` or `~/.config/td/config.yaml`:

**JSON format:**
```json
{
  "data_file": "~/.td.json",
  "theme": {
    "primary_color": "#FF75B7",
    "high_priority_color": "#FF0000",
    "medium_priority_color": "#FFFF00",
    "low_priority_color": "#00FF00"
  },
  "keymap": {
    "add": "a",
    "delete": "d",
    "enter": "enter",
    "quit": "q"
  }
}
```

**YAML format:**
```yaml
data_file: ~/.td.json
theme:
  primary_color: "#FF75B7"
  high_priority_color: "#FF0000"
  medium_priority_color: "#FFFF00"
  low_priority_color: "#00FF00"
keymap:
  add: "a"
  delete: "d"
  enter: "enter"
  quit: "q"
```

## Acknowledgements

This project is a derivative of [todo-cli](https://github.com/yuzuy/todo-cli), which is developed by [Ren Ogaki (yuzuy)](https://github.com/yuzuy) for the purposes of learning the Go language. The original code is licensed under the MIT License.

## License

This project is released under the BSD Zero Clause License (0BSD). For more details, please refer to the [LICENSE](LICENSE) file.

---
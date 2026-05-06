# journal

CLI journal application for tracking daily entries with timestamps.
Inspired by [Zed](https://github.com/zed-industries/zed)'s **journal** feature.

## Installation

### build from source

```bash
git clone https://github.com/halqme/journal.git
cd journal/cmd/journal
go install .
```

## Usage

```bash
# Direct input
journal "今日の振り返り"

# Open editor ($EDITOR)
journal

# Open editor (explicit)
journal -e

# Project mode (create .journal/ in current directory)
journal -p "プロジェクト用メモ"

# Print entries directory (for fzf etc)
journal --root
```

### fzf Integration

```bash
# Browse and open entries
vim "$(find $(journal --root) -name '*.md' | fzf)"

# With project mode
vim "$(find $(journal -p --root) -name '*.md' | fzf)"
```

## Settings

`~/.journal/settings.json`

```json
{
  "hour_format": "24",
  "path": "~/.journal/entries",
  "editor": "",
  "template": "# {time}\n",
  "default_project": false,
  "file_name_format": "02.md",
  "auto_create_settings": true
}
```

### Options

| Key | Description | Default |
|-----|-------------|---------|
| `hour_format` | Time format for headings | `"24"` (or `"12"`) |
| `path` | Base directory for entries | `"~/.journal/entries"` |
| `editor` | Override `$EDITOR` env var | `""` (use env) |
| `template` | Entry header template. `{time}` is replaced with the time value | `"# {time}\n"` |
| `default_project` | Create entries in current directory by default | `false` |
| `file_name_format` | Entry filename format (Go time format) | `"02.md"` |
| `auto_create_settings` | Auto-generate settings file on first run | `true` |

### Template Examples

```json
// Default
"template": "# {time}\n"
// → # 09:55
```

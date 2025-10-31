# 🧱 jot - Terminal Note-Taking CLI

A lightning-fast terminal-based note-taking and journaling CLI that feels like git and fzf had a baby.

**Cross-platform support:** ✅ Linux ✅ macOS ✅ Windows

## Features

- **Speed over style** — open → dump thought → close
- **Text-first** — every note is a plain .md file
- **Git-like IDs** — short hash identifiers with partial matching support
- **Full-text search** — SQLite FTS5 powered search across all content
- **SQLite backend** — fast, reliable database with rich querying capabilities
- **Searchable & local-first** — your notes are yours, no network dependency
- **Cross-platform** — works on Linux, macOS, and Windows
- **Beautiful UI** — colored tags, styled output, progress bars
- **Smart editor detection** — automatically finds best available editor
- **Expandable** — designed for future additions like Git sync, encryption, tags, TUI

## Installation

### Windows

#### Option 1: PowerShell (Recommended)
```powershell
# Run in PowerShell as Administrator
powershell -ExecutionPolicy Bypass -Command "iwr -UseBasicParsing https://raw.githubusercontent.com/sk25469/jot/main/install-windows.ps1 | iex"
```

#### Option 2: Manual Windows Installation
```cmd
# Install Go first from https://golang.org/dl/
go install github.com/sk25469/jot@latest

# Add to PATH if needed: %USERPROFILE%\go\bin
```

#### Option 3: Download Script
Download `install-windows.bat` or `install-windows.ps1` from releases and run.

### Linux/macOS

#### Quick Install
```bash
go install github.com/sk25469/jot@latest
```

#### From Source
```bash
git clone https://github.com/sk25469/jot.git
cd jot
go build -o jot
sudo mv jot /usr/local/bin/  # Linux
# or
mv jot /usr/local/bin/       # macOS
```

### Verify Installation
```bash
jot --help
jot list  # Should work after installation
```

## Usage

### Create a new note
```bash
jot new "Fix offset reset" --tag kafka --tag debugging --mode dev
```

### List all notes
```bash
jot list
# ID       DATE         TITLE                    TAGS
# f4f1c39  2025-10-31   Refactoring complete     refactoring, go
# 5f3f8ed  2025-10-31   Daily reflection         journal

jot list --tag kafka      # Filter by tag
jot list --mode journal   # Filter by mode
```

### Search notes
```bash
# Basic search
jot search "kafka offset"
jot search "debugging"

# Search in titles, content, and tags
jot search "fts"          # Matches tags
jot search "refactoring"  # Matches titles
jot search "implementation" # Matches content

# Full-text search powered by SQLite FTS5
jot search "search terms" # Fast indexed search
```

### Open a note
```bash
# Open by git-like short hash ID
jot open f4f1c39

# Open by partial ID (like git commits)
jot open f4f

# Open by title matching (still works)
jot open "Fix offset reset"
jot open "daily"
```

### View statistics
```bash
jot stats
```

## Note IDs

jot uses **git-like short hash IDs** for each note:

- **Unique 7-character IDs** generated from filename (e.g., `f4f1c39`)
- **Partial matching** - use just the first few characters (e.g., `f4f` instead of `f4f1c39`)
- **Consistent** - same file always gets the same ID
- **No more fake sequential numbers** - IDs actually work with `jot open`

### ID Examples
```bash
jot list              # Shows real usable IDs
jot open f4f1c39      # Open by full ID
jot open f4f          # Open by partial ID
jot open xyz123       # Error: note not found
```

## Configuration

Configuration varies by platform:

### Linux/macOS
Configuration stored in `~/.jot/config.yaml`:

### Windows  
Configuration stored in `%APPDATA%\jot\config.yaml`:

```yaml
editor: "code"              # Windows: code, notepad.exe, notepad++.exe
                           # macOS: code, vim, nano, "open -t"  
                           # Linux: code, vim, nano, gedit
default_mode: "dev"        # Default mode for new notes
storage_path: "auto"       # Platform-specific default location
```

### Editor Detection

jot automatically detects the best available editor:

**Windows Priority:**
1. VS Code (`code`)
2. Notepad++ (`notepad++.exe`) 
3. Vim (`vim.exe`)
4. Notepad (`notepad.exe`) - always available fallback

**macOS Priority:**
1. VS Code (`code`)
2. Vim (`vim`)
3. Nano (`nano`)
4. TextEdit (`open -t`)

**Linux Priority:**
1. VS Code (`code`)
2. Vim (`vim`) 
3. Nano (`nano`)
4. Gedit (`gedit`)

### Storage Locations

**Windows:** `%APPDATA%\jot\` (e.g., `C:\Users\Username\AppData\Roaming\jot\`)  
**macOS/Linux:** `~/.jot/` (e.g., `/home/username/.jot/`)

## Database & Performance

jot uses **SQLite with FTS5** for lightning-fast operations:

- **Full-text search** - Search across titles, content, and tags instantly
- **Indexed queries** - Fast filtering by date, mode, and tags  
- **Rich statistics** - Advanced analytics on your note-taking patterns
- **Automatic sync** - File system and database stay in perfect sync
- **Content tracking** - Detects changes and maintains search index
- **Efficient storage** - Normalized tags, content hashing, word counts

### Database Structure

**Windows:**
```
%APPDATA%\jot\
├── config.yaml         # User configuration
├── jot.db             # SQLite database with FTS index
└── notes\             # Markdown files (source of truth)
    ├── 2025-11-01T01-10-05Z-fix-offset-reset.md
    └── ...
```

**Linux/macOS:**
```
~/.jot/
├── config.yaml         # User configuration  
├── jot.db             # SQLite database with FTS index
└── notes/             # Markdown files (source of truth)
    ├── 2025-11-01T01-10-05Z-fix-offset-reset.md
    └── ...
```

## Note Format

Each note is a markdown file with YAML frontmatter:

```markdown
---
title: Fix offset reset
tags: [kafka, debugging]
mode: dev
date: 2025-11-01T01:10:05Z
---

Your note content goes here...
```

## Future Ideas

- [ ] Git sync (auto-commit every edit)
- [ ] Encrypted mode for personal journaling  
- [ ] Backlinks ([[related note]])
- [ ] `jot web` — minimal read-only web view
- [ ] AI-assisted recall
- [ ] Multi-device sync (Dropbox/GitHub)
- [ ] Interactive TUI mode
- [ ] Note templates
- [ ] Daily/weekly note automation

## Contributing

Pull requests welcome! This is an MVP focused on core functionality.

## License

MIT
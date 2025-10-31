# 🧱 jot - Terminal Note-Taking CLI

A lightning-fast terminal-based note-taking and journaling CLI that feels like git and fzf had a baby.

## Features

- **Speed over style** — open → dump thought → close
- **Text-first** — every note is a plain .md file
- **Searchable & local-first** — your notes are yours, no network dependency
- **Expandable** — designed for future additions like Git sync, encryption, tags, TUI

## Installation

### From Source

```bash
git clone https://github.com/sahilsarwar/jot.git
cd jot
go build -o jot
sudo mv jot /usr/local/bin/
```

### Quick Install Script

```bash
go install github.com/sahilsarwar/jot@latest
```

## Usage

### Create a new note
```bash
jot new "Fix offset reset" --tag kafka --tag debugging --mode dev
```

### List all notes
```bash
jot list
jot list --tag kafka
jot list --mode journal
```

### Search notes
```bash
jot search "kafka offset"
jot search "debugging"
```

### Open a note
```bash
jot open "Fix offset reset"
jot open 2025-11-01T01-10-05Z
```

### View statistics
```bash
jot stats
```

## Configuration

Configuration is stored in `~/.jot/config.yaml`:

```yaml
editor: "nvim"           # Your preferred editor
default_mode: "dev"      # Default mode for new notes
storage_path: "~/.jot/notes"  # Where to store notes
```

## Directory Structure

```
~/.jot/
├── config.yaml
├── notes/
│   ├── 2025-11-01T01-10-05Z-fix-offset-reset.md
│   ├── 2025-10-30T23-14-42Z-night-reflection.md
│   └── ...
└── index.db  # optional metadata cache (future)
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
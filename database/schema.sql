-- jot SQLite Database Schema
-- This schema supports the current functionality while enabling future features

-- Notes table - core note information
CREATE TABLE notes (
    id TEXT PRIMARY KEY,           -- 7-char hash ID (e.g., 'f4f1c39')
    title TEXT NOT NULL,           -- Note title
    mode TEXT NOT NULL DEFAULT 'dev',  -- Note mode (dev, journal, etc.)
    file_path TEXT NOT NULL UNIQUE,    -- Full path to the .md file
    file_name TEXT NOT NULL,       -- Just the filename for easy reference
    content_hash TEXT,             -- Hash of file content for change detection
    created_at DATETIME NOT NULL,  -- When note was created
    updated_at DATETIME NOT NULL,  -- When note was last modified
    content_preview TEXT,          -- First 200 chars of content for quick display
    word_count INTEGER DEFAULT 0   -- Number of words in the note
);

-- Tags table - normalized tag storage
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,     -- Tag name (e.g., 'kafka', 'debugging')
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    usage_count INTEGER DEFAULT 0  -- How many notes use this tag
);

-- Note-Tag junction table - many-to-many relationship
CREATE TABLE note_tags (
    note_id TEXT NOT NULL,         -- References notes.id
    tag_id INTEGER NOT NULL,       -- References tags.id
    PRIMARY KEY (note_id, tag_id),
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Full-text search virtual table for content searching
CREATE VIRTUAL TABLE notes_fts USING fts5(
    note_id UNINDEXED,             -- References notes.id (not indexed in FTS)
    title,                         -- Note title (searchable)
    content,                       -- Full note content (searchable)  
    tags                           -- Space-separated tag names (searchable)
);

-- Triggers to maintain FTS index
-- CREATE TRIGGER notes_fts_insert AFTER INSERT ON notes
-- BEGIN
--     INSERT INTO notes_fts(note_id, title, content, tags)
--     SELECT NEW.id, NEW.title, 
--            (SELECT content FROM note_files WHERE note_id = NEW.id),
--            (SELECT GROUP_CONCAT(t.name, ' ') FROM tags t 
--             JOIN note_tags nt ON t.id = nt.tag_id 
--             WHERE nt.note_id = NEW.id);
-- END;

-- CREATE TRIGGER notes_fts_update AFTER UPDATE ON notes
-- BEGIN
--     UPDATE notes_fts SET 
--         title = NEW.title,
--         content = (SELECT content FROM note_files WHERE note_id = NEW.id),
--         tags = (SELECT GROUP_CONCAT(t.name, ' ') FROM tags t 
--                 JOIN note_tags nt ON t.id = nt.tag_id 
--                 WHERE nt.note_id = NEW.id)
--     WHERE note_id = NEW.id;
-- END;

-- CREATE TRIGGER notes_fts_delete AFTER DELETE ON notes
-- BEGIN
--     DELETE FROM notes_fts WHERE note_id = OLD.id;
-- END;

-- Configuration table for app settings
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert default config values
INSERT INTO config (key, value) VALUES 
    ('editor', 'vim'),
    ('default_mode', 'dev'),
    ('storage_path', '~/.jot/notes'),
    ('db_version', '1.0');

-- Views for common queries

-- View for notes with tag information
CREATE VIEW notes_with_tags AS
SELECT 
    n.id,
    n.title,
    n.mode,
    n.file_path,
    n.created_at,
    n.updated_at,
    n.content_preview,
    n.word_count,
    GROUP_CONCAT(t.name, ', ') as tags
FROM notes n
LEFT JOIN note_tags nt ON n.id = nt.note_id
LEFT JOIN tags t ON nt.tag_id = t.id
GROUP BY n.id, n.title, n.mode, n.file_path, n.created_at, n.updated_at;

-- View for recent notes (last 7 days)
CREATE VIEW recent_notes AS
SELECT * FROM notes_with_tags 
WHERE created_at >= datetime('now', '-7 days')
ORDER BY created_at DESC;

-- View for tag statistics
CREATE VIEW tag_stats AS
SELECT 
    t.name,
    t.usage_count,
    COUNT(nt.note_id) as actual_usage,
    t.created_at
FROM tags t
LEFT JOIN note_tags nt ON t.id = nt.tag_id
GROUP BY t.id, t.name, t.usage_count, t.created_at
ORDER BY actual_usage DESC;

-- Example indexes for performance
CREATE INDEX idx_notes_created_desc ON notes(created_at DESC);
CREATE INDEX idx_notes_updated_desc ON notes(updated_at DESC);
CREATE INDEX idx_notes_mode ON notes(mode);
CREATE INDEX idx_notes_title ON notes(title);
CREATE INDEX idx_notes_mode_created ON notes(mode, created_at DESC);
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_tags_usage ON tags(usage_count DESC);
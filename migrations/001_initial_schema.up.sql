-- Create todo_lists table
CREATE TABLE todo_lists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create tasks table
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    list_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    priority INTEGER NOT NULL DEFAULT 0, -- 0=Low, 1=Medium, 2=High, 3=Critical
    deadline DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE
);

-- Create settings table
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default settings
INSERT INTO settings (key, value) VALUES
    ('reminder_minutes', '60'),
    ('show_completed', 'true'),
    ('auto_save', 'true');

-- Create indexes for better performance
CREATE INDEX idx_tasks_list_id ON tasks(list_id);
CREATE INDEX idx_tasks_completed ON tasks(completed);
CREATE INDEX idx_tasks_deadline ON tasks(deadline);
CREATE INDEX idx_tasks_priority ON tasks(priority);

-- Create trigger to update updated_at timestamps
CREATE TRIGGER update_todo_lists_timestamp 
    AFTER UPDATE ON todo_lists
    FOR EACH ROW
    BEGIN
        UPDATE todo_lists SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_tasks_timestamp 
    AFTER UPDATE ON tasks
    FOR EACH ROW
    BEGIN
        UPDATE tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_settings_timestamp 
    AFTER UPDATE ON settings
    FOR EACH ROW
    BEGIN
        UPDATE settings SET updated_at = CURRENT_TIMESTAMP WHERE key = NEW.key;
    END; 
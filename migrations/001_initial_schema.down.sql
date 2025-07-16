-- Drop triggers
DROP TRIGGER IF EXISTS update_settings_timestamp;
DROP TRIGGER IF EXISTS update_tasks_timestamp;
DROP TRIGGER IF EXISTS update_todo_lists_timestamp;

-- Drop indexes
DROP INDEX IF EXISTS idx_tasks_priority;
DROP INDEX IF EXISTS idx_tasks_deadline;
DROP INDEX IF EXISTS idx_tasks_completed;
DROP INDEX IF EXISTS idx_tasks_list_id;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS todo_lists; 
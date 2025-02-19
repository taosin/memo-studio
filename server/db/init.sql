-- server/db/init.sql
-- CREATE TABLE IF NOT EXISTS notes (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     title TEXT,
--     content TEXT NOT NULL,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP
-- );


-- sqlite3 db/notes.db < db/init.sql

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

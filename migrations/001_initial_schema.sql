-- Initial database schema

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created DATETIME NOT NULL,
    edited DATETIME,
    deleted DATETIME
);

-- SSH Keys table
CREATE TABLE IF NOT EXISTS keys (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type INTEGER NOT NULL,
    data BLOB NOT NULL,
    created DATETIME NOT NULL,
    deleted DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Repositories table (named "permitions" in the codebase)
CREATE TABLE IF NOT EXISTS permitions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created DATETIME NOT NULL,
    edited DATETIME,
    deleted DATETIME
);

-- Roles/Access permissions table
CREATE TABLE IF NOT EXISTS roles (
    role_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    rep_id TEXT NOT NULL,
    branch TEXT,
    created DATETIME NOT NULL,
    deleted DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (rep_id) REFERENCES permitions(id)
); 
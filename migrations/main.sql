PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) NOT NULL PRIMARY KEY, 
    name VARCHAR(256) NOT NULL, 
    email VARCHAR(256) UNIQUE NOT NULL, 
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_user_deleted ON users(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS keys (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(256) NOT NULL,
    
    type INTEGER NOT NULL,
    key BLOB NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_keys_user_id ON keys(user_id);
CREATE INDEX idx_key_body ON keys(key);
CREATE INDEX idx_key_full ON keys(type, key);
CREATE INDEX idx_keys_deleted ON keys(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS repos (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(256) UNIQUE NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_repos_name ON repos(name);
CREATE INDEX idx_repos_deleted ON repos(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS roles (
    user_id VARCHAR(36) NOT NULL,
    repo_id VARCHAR(36) NOT NULL,
    branch VARCHAR(256) NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    PRIMARY KEY (user_id, repo_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (repo_id) REFERENCES repos(id)
);

CREATE INDEX idx_roles_repo_id ON roles(repo_id);
CREATE INDEX idx_roles_deleted ON roles(deleted_at) WHERE deleted_at IS NOT NULL;

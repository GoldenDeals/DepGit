PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXIST users (
    id VARCHAR(36) NOT NULL PRIMARY KEY, 
    name VARCHAR(256) NOT NULL, 
    email VARCHAR(256) NOT NULL, 

    created TIME, 
    edited TIME,
    deleted TIME
);
CREATE INDEX idx_user_email ON users(email);

CREATE TABLE IF NOT EXIST keys (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(256) NOT NULL,

    type INT NOT NULL,
    key BLOB NOT NULL,

    created TIME,
    deleted TIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_key_body ON keys(key);
CREATE INDEX idx_key_full ON keys(type, key);

CREATE TABLE IF NOT EXIST repos (
    id VARCHAR(36) PRIMARY KEY,

    name VARCHAR(256) PRIMARY KEY,

    created TIME, 
    edited TIME,
    deleted TIME
);

CREATE INDEX idx_repos_name ON repos(name);

CREATE TABLE IF NOT EXIST  roles (
    PRIMARY KEY ( user_id, rep_id ),
    user_id VARCHAR(36) NOT NULL,
    rep_id VARCHAR(36) NOT NULL,

    branch VARCHAR(256) NOT NULL,

    created TIME,
    deleted TIME,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (rep_id) REFERENCES repos(id)
);

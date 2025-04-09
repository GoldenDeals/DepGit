-- Add indexes to improve query performance

-- Add index for users' email lookup
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Add index for keys by user_id for faster filtering
CREATE INDEX IF NOT EXISTS idx_keys_user_id ON keys(user_id);

-- Add indexes for role lookups
CREATE INDEX IF NOT EXISTS idx_roles_user_id ON roles(user_id);
CREATE INDEX IF NOT EXISTS idx_roles_rep_id ON roles(rep_id); 
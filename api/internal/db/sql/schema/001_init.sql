-- Accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(32) NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(127) NOT NULL,
    created_by UUID NOT NULL REFERENCES accounts(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id),
    author_id UUID NOT NULL REFERENCES accounts(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_accounts_username ON accounts(username);
CREATE INDEX idx_rooms_created_by ON rooms(created_by);
CREATE INDEX idx_messages_room_id ON messages(room_id);
CREATE INDEX idx_messages_author ON messages(author_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- Sample data
INSERT INTO accounts (id, username, password_hash) VALUES
    ('5e305dee-d8d8-49b3-ad6c-73037e58601a', 'testuser', '$2a$10$ZQ9Z9Z9Z9Z9Z9Z9Z9Z9Z9eKGX7JQ7J7J7J7J7J7J7J7J7J7J7J7J7');

INSERT INTO rooms (id, name, created_by) VALUES
    ('8481027d-d6f6-402f-ae6d-98571e8f6496', 'General', '5e305dee-d8d8-49b3-ad6c-73037e58601a');

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    username VARCHAR(30) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

create table if not exists chats (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_by UUID REFERENCES users(id),
    type VARCHAR(20) NOT NULL CHECK ("type" IN ('private', 'group'))
);

CREATE TABLE IF NOT EXISTS chat_members (
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'member')), -- 'admin', 'member'
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    file_url VARCHAR(255),
    user_id INT NOT NULL REFERENCES users(id),
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE group_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_sender_fk_users INTEGER REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    guild_fk_guilds INTEGER REFERENCES guilds(id) ON DELETE CASCADE,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
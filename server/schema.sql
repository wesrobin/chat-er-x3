CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    content TEXT NOT NULL,

    sent_at DATETIME NOT NULL,
    sent_by INTEGER NOT NULL,
    room_id INTEGER NOT NULL,

    FOREIGN KEY (sent_by) REFERENCES users(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);
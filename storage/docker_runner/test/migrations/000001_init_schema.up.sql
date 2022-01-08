CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE,
    password VARCHAR(256) NOT NULL,
    email VARCHAR(256) NOT NULL
);

CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    raw VARCHAR(1000) NOT NULL,
    shortened VARCHAR(256) NOT NULL,
    user_id INTEGER NOT NULL,
    month INTEGER DEFAULT 0,
    week INTEGER DEFAULT 0,
    day INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
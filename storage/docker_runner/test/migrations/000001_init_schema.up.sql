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
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE redirects (
    id SERIAL PRIMARY KEY,
    url_id INTEGER NOT NULL,
    date_of_usage DATE,
    FOREIGN KEY (url_id) REFERENCES urls(id)
);

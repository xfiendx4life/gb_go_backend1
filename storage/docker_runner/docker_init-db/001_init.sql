CREATE USER xfiendx4life
WITH PASSWORD '123456';

CREATE DATABASE shortener
    WITH OWNER xfiendx4life
    TEMPLATE = 'template0'
    ENCODING = 'utf-8'
    LC_COLLATE = 'C.UTF-8'
    LC_CTYPE = 'C.UTF-8';

\connect shortener

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


GRANT ALL PRIVILEGES ON DATABASE shortener TO xfiendx4life;
GRANT ALL ON users to xfiendx4life;
GRANT ALL ON redirects to xfiendx4life;
GRANT ALL ON redirects_id_seq to xfiendx4life;
GRANT ALL ON urls to xfiendx4life;
GRANT ALL ON urls_id_seq to xfiendx4life;
GRANT ALL ON users to xfiendx4life;
GRANT ALL ON users_id_seq to xfiendx4life;

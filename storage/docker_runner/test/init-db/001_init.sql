CREATE USER xfiendx4life
WITH PASSWORD '123456';

CREATE DATABASE shortener
    WITH OWNER xfiendx4life
    TEMPLATE = 'template0'
    ENCODING = 'utf-8'
    LC_COLLATE = 'C.UTF-8'
    LC_CTYPE = 'C.UTF-8';

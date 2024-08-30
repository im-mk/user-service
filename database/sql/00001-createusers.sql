
--liquibase formatted sql

--changeset user:00001
--comment: users table contains the general user information and login credentials.
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

--rollback DROP TABLE IF EXISTS users;
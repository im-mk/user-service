
--liquibase formatted sql

--changeset user:00001
--comment: users table contains the general user information and login credentials.
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name VARCHAR(100),
    middle_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_users_username ON users(username);

--rollback DROP INDEX IF EXISTS idx_users_username; 
--rollback DROP TABLE IF EXISTS users;
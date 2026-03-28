# user-service

A Go microservice demonstrating secure user management with:

-   RSA-signed JWT authentication
-   Refresh token rotation
-   PostgreSQL persistence
-   Liquibase schema migrations
-   JWKS endpoint for public key discovery
-   Swagger API documentation
-   Full unit test coverage

This project is designed as a clean, production-style microservice
showcasing proper layering (controllers → services → repositories),
token handling, and containerised infrastructure.

------------------------------------------------------------------------

# 🚀 Features

-   Create, bootstrap, and retrieve user records\
-   Secure login endpoint issuing RSA-signed JWT access tokens\
-   Refresh endpoint issuing new access tokens using stored refresh
    tokens\
-   Configurable token issuer, audience, expiry, and key identifiers\
-   `/jwks.json` endpoint exposing public keys for downstream
    verification\
-   Liquibase-driven schema migrations (executed in a dedicated
    container)\
-   Swagger API documentation auto-generated at startup\
-   Configuration via JSON file and environment variables (Viper)\
-   Comprehensive unit tests for controllers, services, and repositories

------------------------------------------------------------------------

# 🏗 Architecture Overview

Client → API (Gin) → Services → Repositories (SQLx) → PostgreSQL

Authentication flow:

1.  User logs in → receives Access Token + Refresh Token\
2.  Access token expires\
3.  Client calls refresh endpoint\
4.  New access token issued if refresh token is valid

------------------------------------------------------------------------

# 📦 Tech Stack

-   Go 1.25+
-   Gin
-   SQLx
-   PostgreSQL
-   Liquibase
-   RSA (JWT signing)
-   Docker & Docker Compose
-   Swagger (swaggo)

------------------------------------------------------------------------

# 🛠 Prerequisites

-   Go 1.25+
-   Docker
-   Docker Compose

------------------------------------------------------------------------

# 🔑 Initial Setup

Generate the RSA keypair used for signing JWT tokens:

    make create-keys

------------------------------------------------------------------------

# ⚙️ Environment Variables

Required variables (see docker-compose.yaml):

-   USER_SERVICE_DB_USER
-   USER_SERVICE_DB_PASSWORD
-   USER_SERVICE_DB_NAME
-   PGADMIN_DEFAULT_EMAIL
-   PGADMIN_DEFAULT_PASSWORD
-   PGADMIN_PORT
-   USER_SERVICE_PORT

------------------------------------------------------------------------

# ▶️ Running the Service

Start the full stack:

    make start

Swagger UI: http://localhost:8040/swagger/index.html

Swagger JSON: http://localhost:8040/swagger/doc.json

Start only database + migrations:

    make start-user-db

Stop everything:

    make down

------------------------------------------------------------------------

# 🧪 Development

All Go source code lives in:

    src/

Build:

    make

Run tests:

    make test-user-api

------------------------------------------------------------------------

# 🔐 Authentication

JWT includes:

-   sub (user ID)
-   iss (issuer)
-   aud (audience)
-   exp (expiration)
-   scope (permissions like user:read)
-   kid (key identifier)

Tokens are:

-   RSA signed
-   Audience validated
-   Expiry validated
-   Scope checked via middleware

Refresh tokens:

-   Stored in DB
-   Used to issue new access tokens
-   Expire independently

JWKS endpoint:

GET /jwks.json

------------------------------------------------------------------------

# 🗄 Database Migrations

Liquibase image:

liquibase/liquibase:5.0

If running outside Docker, provide PostgreSQL JDBC driver under:

database/lib

------------------------------------------------------------------------

# 🧾 Configuration

Configuration file:

src/config.json

Environment variables override file values via Viper.

Auth settings include:

-   PrivateKeyPath
-   PublicKeyPath
-   Issuer
-   Audience
-   TokenExpirySeconds
-   RefreshTokenExpirySeconds
-   Kid

------------------------------------------------------------------------

# 🛡 Security Notes

-   Short-lived access tokens
-   Stored refresh tokens
-   Audience validation
-   Scope-based authorization
-   JWKS for external verification

------------------------------------------------------------------------

Happy hacking 🚀

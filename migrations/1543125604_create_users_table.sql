-- +migrate Up
CREATE TABLE users (
                     id SERIAL PRIMARY KEY,
                     uuid uuid NOT NULL UNIQUE,
                     email TEXT NOT NULL UNIQUE,
                     is_guest BOOLEAN NOT NULL,
                     encrypted_password TEXT NOT NULL,
                     updated_at TIMESTAMPTZ NOT NULL,
                     created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX users_uuid ON users (uuid);
CREATE TABLE sessions (
                    id SERIAL PRIMARY KEY,
                    uuid uuid NOT NULL UNIQUE,
                    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
                    token TEXT NOT NULL,
                    expiration TIMESTAMPTZ NOT NULL,
                    created_at TIMESTAMPTZ NOT NULL

);
CREATE INDEX sessions_uuid ON sessions (uuid);
CREATE INDEX sessions_token ON sessions (token);
CREATE TABLE scope_groupings (
                    id SERIAL PRIMARY KEY,
                    uuid uuid NOT NULL UNIQUE,
                    session_id INTEGER REFERENCES sessions(id) ON DELETE CASCADE NOT NULL,
                    scopes TEXT[] NOT NULL,
                    expiration TIMESTAMPTZ NOT NULL,
                    created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX scope_groupings_uuid ON scope_groupings (uuid);
CREATE TABLE password_reset_tokens (
                   id SERIAL PRIMARY KEY,
                   uuid uuid NOT NULL UNIQUE,
                   user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
                   token TEXT NOT NULL UNIQUE,
                   expiration TIMESTAMPTZ NOT NULL,
                   created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX password_reset_tokens_uuid ON password_reset_tokens (uuid);
CREATE INDEX password_reset_token_token ON password_reset_tokens (token);

-- +migrate Down
DROP TABLE password_reset_tokens;
DROP TABLE scope_groupings;
DROP TABLE sessions;
DROP TABLE users;
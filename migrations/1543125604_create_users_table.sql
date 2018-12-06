-- +migrate Up
CREATE TABLE users (
                     uuid uuid NOT NULL UNIQUE PRIMARY KEY,
                     email TEXT NOT NULL UNIQUE,
                     is_guest BOOLEAN NOT NULL,
                     encrypted_password TEXT NOT NULL,
                     updated_at TIMESTAMPTZ NOT NULL,
                     created_at TIMESTAMPTZ NOT NULL
);
CREATE TABLE sessions (
                    uuid uuid NOT NULL UNIQUE PRIMARY KEY,
                    user_uuid uuid REFERENCES users(uuid) ON DELETE CASCADE NOT NULL,
                    token TEXT NOT NULL,
                    created_at TIMESTAMPTZ NOT NULL

);
CREATE INDEX sessions_uuid ON sessions (uuid);
CREATE INDEX sessions_token ON sessions (token);
CREATE TABLE scope_groupings (
                    uuid uuid NOT NULL UNIQUE PRIMARY KEY,
                    session_uuid uuid REFERENCES sessions(uuid) ON DELETE CASCADE NOT NULL,
                    scopes TEXT[] NOT NULL,
                    expiration TIMESTAMPTZ NOT NULL,
                    created_at TIMESTAMPTZ NOT NULL
);
CREATE TABLE password_resets (
                   uuid uuid NOT NULL UNIQUE PRIMARY KEY,
                   user_uuid uuid REFERENCES users(uuid) ON DELETE CASCADE NOT NULL,
                   token TEXT NOT NULL UNIQUE,
                   expiration TIMESTAMPTZ NOT NULL,
                   created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX password_resets_token ON password_resets (token);

-- +migrate Down
DROP TABLE password_resets;
DROP TABLE scope_groupings;
DROP TABLE sessions;
DROP TABLE users;
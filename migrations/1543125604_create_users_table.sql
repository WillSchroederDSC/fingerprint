-- +migrate Up
CREATE TABLE users (
                     id SERIAL PRIMARY KEY,
                     uuid uuid NOT NULL UNIQUE,
                     email TEXT NOT NULL UNIQUE,
                     encrypted_password TEXT NOT NULL,
                     created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX users_uuid ON users (uuid);
CREATE TABLE sessions (
                    id SERIAL PRIMARY KEY,
                    uuid uuid NOT NULL UNIQUE,
                    user_id INTEGER REFERENCES users(id),
                    token TEXT NOT NULL,
                    expiration TIMESTAMPTZ NOT NULL,
                    created_at TIMESTAMPTZ NOT NULL

);
CREATE INDEX sessions_uuid ON sessions (uuid);
create INDEX sessions_token ON sessions (token);
CREATE TABLE scope_groupings (
                    id SERIAL PRIMARY KEY,
                    uuid uuid NOT NULL UNIQUE,
                    session_id INTEGER REFERENCES sessions(id),
                    scopes TEXT[] NOT NULL,
                    expiration TIMESTAMPTZ NOT NULL,
                    created_at TIMESTAMPTZ NOT NULL
);
CREATE index scope_groupings_uuid ON scope_groupings (uuid);

-- +migrate Down
DROP TABLE scope_groupings;
DROP TABLE sessions;
DROP TABLE users;
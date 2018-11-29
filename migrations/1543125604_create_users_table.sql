-- +migrate Up
CREATE TABLE users (
                     id SERIAL PRIMARY KEY,
                     uuid uuid NOT NULL UNIQUE,
                     email TEXT NOT NULL UNIQUE,
                     encrypted_password TEXT NOT NULL,
                     created_at TIMESTAMPTZ
);
CREATE index users_uuid ON users (uuid);
CREATE TABLE sessions (
                    id SERIAL PRIMARY KEY,
                    uuid uuid NOT NULL UNIQUE,
                    user_id INTEGER REFERENCES users(id),
                    expiration TIMESTAMPTZ,
                    created_at TIMESTAMPTZ

);
CREATE index sessions_uuid ON sessions (uuid);
CREATE TABLE scope_groupings (
                    id SERIAL PRIMARY KEY,
                    uuid uuid NOT NULL UNIQUE,
                    session_id INTEGER REFERENCES sessions(id),
                    scopes TEXT[],
                    expiration TIMESTAMPTZ,
                    created_at TIMESTAMPTZ
);
CREATE index scope_groupings_uuid ON scope_groupings (uuid);

-- +migrate Down
DROP TABLE scope_groupings;
DROP TABLE sessions;
DROP TABLE users;